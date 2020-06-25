package ofc

import (
	"automaton/automaton"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JamesHovious/w32"
	"golang.org/x/image/bmp"
)

type GameContext struct {
	Hwnd             w32.HWND // memu handle
	ImageServerHost  string
	SolverServerHost string
	ScreenShotDir    string // store screenshots here
}

type GameContextError struct {
	Msg string
}

func (e *GameContextError) Error() string {
	return fmt.Sprintf("GameContextError: %s", e.Msg)
}
func (gCtxt *GameContext) CaptureGameState(imageName string) error {
	// get screenshot
	img := automaton.ScreenCapture(gCtxt.Hwnd)

	// save image to file
	f, err := os.Create(gCtxt.ScreenShotDir + "\\" + imageName + ".bmp")
	if err != nil {
		return err
	}

	// save image to file
	err = bmp.Encode(f, &img) // *Bitmap implements Image, not Bitmap!
	if err != nil {
		return err
	}

	return nil
}

func (gCtxt *GameContext) ParseGameStateFromImage(imageName string) (*GameState, error) {
	u, err := url.Parse(gCtxt.ImageServerHost)
	if err != nil {
		return nil, err
	}

	// image server needs a filename parameter, and config parameter
	q := u.Query()
	q.Set("filename", gCtxt.ScreenShotDir+"\\"+imageName+".bmp")
	q.Set("config", "MEMU_536x983")
	u.RawQuery = q.Encode()
	fmt.Printf("constructed url: %v\n", u)

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	gs, err := parseGameStateFromJson(string(body))
	if err != nil {
		return nil, err
	}
	return gs, nil
}

func (gCtxt *GameContext) SolveGameState(gs *GameState) ([]Action, error) {
	u, err := url.Parse(gCtxt.SolverServerHost)
	if err != nil {
		return nil, err
	}

	// build query for the OfcSolver
	q := u.Query()
	q.Set("type", "progressive")
	fmt.Println(formatHand(&gs.MyHand))
	q.Set("my_hand", formatHand(&gs.MyHand))

	for _, otherHand := range gs.OtherHands {
		if !otherHand.Empty() {
			q.Add("other_hand", formatHand(&otherHand))
		}
	}
	q.Set("n_solves", "700")
	q.Set("n_decision_solves", "25")
	q.Set("pull", formatCardArray(gs.Pull))

	if len(gs.DeadCards) > 0 {
		q.Set("dead_cards", formatCardArray(gs.DeadCards))
	}

	u.RawQuery = q.Encode()
	fmt.Printf("constructed url: %v\n", u)

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// parse json response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	actionStr := resp["best"].(string)
	return parseActions(actionStr, gs)
}

func formatHand(h *Hand) string {
	s := ""
	s += formatCardArray(h.Top)
	s += "/"
	s += formatCardArray(h.Mid)
	s += "/"
	s += formatCardArray(h.Bot)
	return s
}

func formatCardArray(cards []Card) string {
	s := ""
	for _, c := range cards {
		if s != "" {
			s += " "
		}
		s += c.Val
	}
	return s
}

func parseActions(actionString string, gs *GameState) ([]Action, error) {
	// "3d bot Ad mid 3c bot Ah mid Qs top"
	tokens := strings.Split(actionString, " ")
	if len(tokens)%2 != 0 {
		return nil, &OfcError{"Actions string doesn't have an even number of actions!"}
	}

	actions := make([]Action, 0)
	for i := 0; i < len(tokens)/2; i++ {
		actions = append(actions, parseAction(tokens[2*i], tokens[2*i+1], gs))
	}

	return actions, nil
}

func parseAction(card string, position string, gs *GameState) Action {
	var c Card
	for _, pullCard := range gs.Pull {
		if card == pullCard.Val {
			c = pullCard
			break
		}
	}

	var p Position
	switch position {
	case "top":
		p = Top
	case "mid":
		p = Mid
	case "bot":
		p = Bot
	}

	return Action{Card: c, Pos: p}
}

func (gCtxt *GameContext) ExecuteActions(actions []Action, gs *GameState) error {
	// find an empty card slot, and drag the card there
	topCounter := 0
	midCounter := 0
	botCounter := 0
	for _, action := range actions {
		log.Printf("\nExecuting action: %+v\n", action)

		// find empty card slot in specified position
		var slotCoords Coord

		// FixMe: bounds checks?
		switch action.Pos {
		case Top:
			slotCoords = gs.EmptyCards.Top[topCounter]
			topCounter++
		case Mid:
			slotCoords = gs.EmptyCards.Mid[midCounter]
			midCounter++
		case Bot:
			slotCoords = gs.EmptyCards.Bot[botCounter]
			botCounter++
		}

		// log.Printf("Found slot %+v\n", slotCoords)

		cardCoord, err := gCtxt.findCardCoord(action, gs, len(actions) >= 13)
		if err != nil {
			return err
		}

		automaton.MoveMouse(gCtxt.Hwnd, cardCoord.X, cardCoord.Y)
		time.Sleep(50 * time.Millisecond)
		automaton.ClickDown()
		time.Sleep(50 * time.Millisecond)
		automaton.MoveMouse(gCtxt.Hwnd, slotCoords.X, slotCoords.Y)
		time.Sleep(50 * time.Millisecond)
		automaton.ClickUp()
		time.Sleep(350 * time.Millisecond)
	}

	return nil
}

func (gCtxt *GameContext) findCardCoord(action Action, gs *GameState, isFantasy bool) (*Coord, error) {
	if !isFantasy {
		return &action.Card.Coord, nil
	}

	// special case: Fantasy. If fantasy, every time we take a card, the other card
	// positions change. So after every executed action, we have to re-identify
	// the card positions.
	timeString := strconv.Itoa(int(time.Now().Unix()))
	err := gCtxt.CaptureGameState(timeString)
	if err != nil {
		return nil, err
	}
	newGs, err := gCtxt.ParseGameStateFromImage(timeString)
	if err != nil {
		return nil, err
	}

	// find the new card in newGs, return that coord
	for _, c := range newGs.Pull {
		if action.Card.Val == c.Val {
			return &c.Coord, nil
		}
	}

	return nil, &OfcError{"Couldn't find new card in image tmp"}
}
