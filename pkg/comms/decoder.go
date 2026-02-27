package comms

import (
	"fmt"
	"image"
	"image/color"
)

// DecodeRGB332 converts a flat RGB332 byte slice into an *image.RGBA.
// Each byte encodes one pixel: bits [7:5]=R3, [4:2]=G3, [1:0]=B2.
func DecodeRGB332(data []byte, width, height int) (*image.RGBA, error) {
	if len(data) != width*height {
		return nil, fmt.Errorf("DecodeRGB332: data length %d != %d×%d (%d)", len(data), width, height, width*height)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i, b := range data {
		r := uint8(uint16(b>>5) * 255 / 7)
		g := uint8(uint16((b>>2)&0x07) * 255 / 7)
		blue := uint8(uint16(b&0x03) * 255 / 3)
		x := i % width
		y := i / width
		img.Set(x, y, color.RGBA{R: r, G: g, B: blue, A: 255})
	}
	return img, nil
}
