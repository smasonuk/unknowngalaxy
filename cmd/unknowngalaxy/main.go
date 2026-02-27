package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/smasonuk/si3d/pkg/si3d"
	"github.com/smasonuk/unknowngalaxy/pkg/universe"
)

const HEIGHT = 128
const WIDTH = 128

func saveImageToFile(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
func main() {
	fmt.Println("Generating mountains...")
	mountains := si3d.NewSubdividedPlaneHeightMapPerlin(
		10000, 10000,
		color.RGBA{R: 153, G: 196, B: 210, A: 255},
		35, 800, 800, 42,
	)
	mountains.SetDrawLinesOnly(true)

	// 1. Setup the Scene
	scene := universe.NewLocalScene(10000, 25000, 35000, 0, 0, 0)
	scene.AddEntity(&si3d.Entity{Model: mountains, X: 0, Y: 0, Z: 0})

	// 2. Position the Probe
	// We'll place it at LocalZ = -400, high up at LocalY = -200
	startPos := universe.NewGalacticPosition(10000, 25000, 35000, 0, 0, 0, 0, -200.0, -400.0)
	probe := universe.NewProbe("Voyager-1", startPos)

	// 3. Point the probe's camera straight ahead (towards Z = 0)
	// We'll look slightly down at Y = -20.0 to see the passing mountains
	lookAtTarget := si3d.NewVector3(0, -200.0, 0)
	probe.PointCamera(lookAtTarget)

	// Take the picture
	fmt.Println("Taking picture...")
	frameImg := scene.TakePicture(probe, WIDTH, HEIGHT)

	// Save the final image
	filename := "probe_capture.png"
	saveImageToFile(frameImg, filename)
	fmt.Println("Done! Saved to", filename)
}
