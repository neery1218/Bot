package automaton

import (
	"image/color"
	"unsafe"

	"github.com/JamesHovious/w32"
)

func ScreenCapture(hwnd w32.HWND) Bitmap {

	// get relative coordinates of window (used to compute size)
	rect := w32.GetClientRect(hwnd)
	// fmt.Printf("left right top down %v %v %v %v\n", rect.Left, rect.Right, rect.Top, rect.Bottom)

	// dc = device context. idk what these are for
	hdcScreen := w32.GetDC(hwnd)
	hdc := w32.CreateCompatibleDC(hdcScreen)

	// initialize memory to store image of window.
	hbmp := w32.CreateCompatibleBitmap(hdcScreen, uint(rect.Right-rect.Left), uint(rect.Bottom-rect.Top))
	w32.SelectObject(hdc, w32.HGDIOBJ(hbmp))
	PrintWindow(hwnd, hdc) // write bitmap of hwnd to the memory pointed to by hdc
	// hbmp isn't readable by us yet, not sure why it's implemented this way.

	// use hbmp to make a system call to figure out dims of bitmap. (height, width)
	var dib w32.DIBSECTION
	ret := w32.GetObject(w32.HGDIOBJ(hbmp), unsafe.Sizeof(dib), unsafe.Pointer(&dib))
	if ret == 0 { // ret holds #bytes written to dib
		panic("Getting object failed!")
	}

	// fmt.Printf("%+v\n", dib)
	if dib.DsBm.BmBitsPixel != 32 {
		panic("assertion failed. bits per pixel should be 32")
	}

	// initialize BITMAPINFO struct with dimension values. this is needed to find the rgb values
	var bmpInfo w32.BITMAPINFO
	var bmpHeader w32.BITMAPINFOHEADER

	bmpInfo.BmiHeader.BiSize = uint32(unsafe.Sizeof(bmpHeader))
	bmpInfo.BmiHeader.BiWidth = dib.DsBm.BmWidth
	bmpInfo.BmiHeader.BiHeight = dib.DsBm.BmHeight
	bmpInfo.BmiHeader.BiPlanes = dib.DsBm.BmPlanes
	bmpInfo.BmiHeader.BiBitCount = dib.DsBm.BmBitsPixel
	bmpInfo.BmiHeader.BiCompression = w32.BI_RGB
	// this is in bytes
	bmpInfo.BmiHeader.BiSizeImage = uint32(((bmpInfo.BmiHeader.BiWidth*32 + 31) & ^31) / 8 * bmpInfo.BmiHeader.BiHeight)
	bmpInfo.BmiHeader.BiClrImportant = 0

	// Copying the bitmap only works if i use this specific function.
	// i think it's bc the windows system calls can't just write to _any_ memory region. idfk
	buf := w32.GlobalAlloc(w32.GMEM_MOVEABLE, bmpInfo.BmiHeader.BiSizeImage)
	defer w32.GlobalFree(buf) // TODO: monitor mem usage to make sure this actually works

	// lock the memory so no other process can use it.
	memptr := w32.GlobalLock(buf)
	defer w32.GlobalUnlock(buf)

	// copy rgb array from wherever it is to memptr
	numLines := GetDIBits(
		hdc,
		hbmp,
		0,
		bmpInfo.BmiHeader.BiHeight,
		uintptr(memptr),
		uintptr(unsafe.Pointer(&bmpInfo)),
		w32.DIB_RGB_COLORS)
	// fmt.Printf("Num lines written: %v\n", numLines)
	if numLines == 0 {
		panic("Numlines shouldn't be zero!")
	}

	// initialize our image structure
	var img Bitmap
	img.ColorModelVal = color.RGBAModel
	img.ColorArray = make([][]color.Color, bmpInfo.BmiHeader.BiHeight)
	for i := range img.ColorArray {
		img.ColorArray[i] = make([]color.Color, bmpInfo.BmiHeader.BiWidth)
	}

	// copy rgba values from memptr to win32bitmap struct
	head := uintptr(memptr)
	for i := uint32(0); i < bmpInfo.BmiHeader.BiSizeImage; i += 4 {
		// r g b a
		b := *(*uint8)(unsafe.Pointer(head))
		g := *(*uint8)(unsafe.Pointer(head + 1))
		r := *(*uint8)(unsafe.Pointer(head + 2))
		a := *(*uint8)(unsafe.Pointer(head + 3))

		row := (i / 4) / uint32(bmpInfo.BmiHeader.BiWidth)
		col := (i / 4) - row*uint32(bmpInfo.BmiHeader.BiWidth)

		img.ColorArray[len(img.ColorArray)-int(row)-1][col] = color.RGBA{R: r, G: g, B: b, A: a}
		head += 4
	}

	return img
}
