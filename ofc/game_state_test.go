package ofc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidGameState(t *testing.T) {
	str := `{
		"MyHand" : { 
			"Top": [
				{"Val" : "Ah", "Coord" : { "X" : 20, "Y" : 30}}, 
				{"Val" : "Kh", "Coord" : { "X" : 30, "Y" : 30}}
			], 
			"Mid": [
				{"Val" : "3s", "Coord" : { "X" : 100, "Y" : 100}}
			],
			"Bot" : [
				{"Val" : "5d", "Coord" : { "X" : 40, "Y" : 40}}
			]
		},
		"OtherHands" : [ 
			{ 
				"Top": [], 
				"Mid": [], 
				"Bot": []
			}
		],
		"Pull" : [ 
			{"Val" : "4s", "Coord" : { "X" : 50, "Y" : 60 }}, 
			{"Val" : "8d", "Coord" : { "X" : 40, "Y" : 90 }}, 
			{"Val" : "9d", "Coord" : { "X" : 80, "Y" : 100}}
		],
		"DeadCards" : [ 
			{"Val" : "9h", "Coord" : { "X" : 90, "Y": 30 }}
		]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.Nil(t, err)
	assert.Equal(t, gs.MyHand.Top[0], Card{"Ah", Coord{20, 30}})
}

/*
func TestInvalidGameStateInvalidMyHand(t *testing.T) {
	str := `{
		"MyHand" : { "Top": ["Ah", "Kh", "Kd", "4d"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [],
		"Pull" : { "4s" : { "X" : -50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "9h" ]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.NotNil(t, err)
	assert.Nil(t, gs)
}
func TestInvalidGameStateNoOtherHands(t *testing.T) {
	str := `{
		"MyHand" : { "Top": ["Ah", "Kh"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [],
		"Pull" : { "4s" : { "X" : -50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "9h" ]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.NotNil(t, err)
	assert.Nil(t, gs)
}

func TestInvalidGameStateInvalidOtherHands(t *testing.T) {
	str := `{
		"MyHand" : { "Top": ["Ah", "Kh"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [{ "Top": ["4h"], "Mid": ["2d", "3d", "4d", "5d", "6d", "7d"], "Bot" : ["5s"] }],
		"Pull" : { "4s" : { "X" : -50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "9h" ]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.NotNil(t, err)
	assert.Nil(t, gs)
}

func TestInvalidGameStateBadPullLength(t *testing.T) {
	str := `{
		"MyHand" : { "Top": ["Ah", "Kh"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [{ "Top": ["4d"], "Mid": ["2d", "3d"], "Bot" : ["5s"] }],
		"Pull" : { "4s" : { "X" : -50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90} },
		"DeadCards" : [ "9h" ]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.NotNil(t, err)
	assert.Nil(t, gs)
}
func TestInvalidGameStateBadPullCoords(t *testing.T) {
	str := `{
		"MyHand" : { "Top": ["Ah", "Kh"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [{ "Top": ["4d"], "Mid": ["2d", "3d"], "Bot" : ["5s"] }],
		"Pull" : { "4s" : { "X" : -50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "9h" ]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.NotNil(t, err)
	assert.Nil(t, gs)
}

func TestInvalidGameStateInvalidCard(t *testing.T) {
	str := `{
		"MyHand" : { "Top": ["Ah", "Kha"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [{ "Top": ["4d"], "Mid": ["2d", "3d"], "Bot" : ["5s"] }],
		"Pull" : { "4s" : { "X" : 50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "9h" ]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.NotNil(t, err)
	assert.Nil(t, gs)
}

func TestInvalidGameStateDuplicateCard(t *testing.T) {
	str := `{
		"MyHand" : { "Top": ["Ah", "Kh"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [{ "Top": ["4d"], "Mid": ["2d", "3d"], "Bot" : ["5s"] }],
		"Pull" : { "4s" : { "X" : 50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "4d" ]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.NotNil(t, err)
	assert.Nil(t, gs)
}

func TestStateChanged(t *testing.T) {
	str1 := `{
		"MyHand" : { "Top": ["Ah", "Kh"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [{ "Top": ["4d"], "Mid": ["2d", "3d"], "Bot" : ["5s"] }],
		"Pull" : { "4s" : { "X" : 50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "5h" ]
	}`
	str2 := `{
		"MyHand" : { "Top": ["Ah", "Kh"], "Mid": ["3s", "3c"], "Bot" : ["5d", "5c"] },
		"OtherHands" : [{ "Top": ["4d"], "Mid": ["2d", "3d"], "Bot" : ["5s"] }],
		"Pull" : { "4s" : { "X" : 50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "5h" ]
	}`

	gs1, _ := parseGameStateFromJson(str1)
	gs2, _ := parseGameStateFromJson(str2)
	assert.True(t, StateChanged(gs1, gs2))
	assert.False(t, StateChanged(gs1, gs1))
}

func TestDecisionRequired(t *testing.T) {
	str := `{
		"MyHand" : { "Top": ["Ah", "Kh"], "Mid": ["3s"], "Bot" : ["5d"] },
		"OtherHands" : [{ "Top": ["4d"], "Mid": ["2d", "3d"], "Bot" : ["5s"] }],
		"Pull" : { "4s" : { "X" : 50, "Y" : 60 }, "8d" : { "X" : 40, "Y" : 90}, "9d" : { "X" : 80, "Y" : 100 } },
		"DeadCards" : [ "5h" ]
	}`

	gs, err := parseGameStateFromJson(str)
	assert.Nil(t, err)
	assert.NotNil(t, gs)
	assert.True(t, gs.DecisionRequired())

	gs.Pull = make(map[Card]Coord)
	assert.False(t, gs.DecisionRequired())
}
*/
