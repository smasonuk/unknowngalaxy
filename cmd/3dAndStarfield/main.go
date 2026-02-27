package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"os"

	"github.com/smasonuk/si3d/pkg/si3d"
	"github.com/smasonuk/unknowngalaxy/pkg/universe"
)

const HEIGHT = 128
const WIDTH = 128

type Probe struct {
	StarfieldPosition si3d.Vector3 // Probe's location in the galaxy (Light-Years)
	Camera            *si3d.Camera // The Master Camera
}

func NewProbe(starfieldPosition si3d.Vector3, cam *si3d.Camera) *Probe {
	return &Probe{
		StarfieldPosition: starfieldPosition,
		Camera:            cam,
	}
}

func main() {
	// Bumped to 120 frames for a smoother, slightly slower 360-degree flight
	frames := 280
	outGif := &gif.GIF{}

	// 1. Where is the probe in the GALAXY?
	starfieldPosition := si3d.NewVector3(10000, 25000, 35000)

	// 2. Pre-generate the Mountains so we don't rebuild the heightmap every frame
	fmt.Println("Generating mountains...")
	mountains := si3d.NewSubdividedPlaneHeightMapPerlin(
		10000, 10000,
		color.RGBA{R: 153, G: 196, B: 210, A: 255},
		35, 800, 800, 42,
	)
	mountains.SetDrawLinesOnly(true)

	fmt.Println("Rendering animation frames...")

	for i := 0; i < frames; i++ {
		// 3. Calculate progression (0.0 to 1.0) and convert to radians
		progress := float64(i) / float64(frames)
		theta := progress * (math.Pi * 2.0) // Full 360 circle

		circleRadius := 400.0
		cameraHeight := -200.0

		// 4. Calculate camera POSITION on the circle
		localCamX := math.Cos(theta) * circleRadius
		localCamZ := math.Sin(theta) * circleRadius

		// 5. Calculate the TANGENT vector (the direction straight ahead on the curve)
		// The derivative of (cos, sin) is (-sin, cos)
		tangentX := -math.Sin(theta)
		tangentZ := math.Cos(theta)

		// 6. Project the LookAt target out along the tangent vector
		lookDistance := 100.0
		lookAtX := localCamX + (tangentX * lookDistance)
		lookAtZ := localCamZ + (tangentZ * lookDistance)
		lookAtY := cameraHeight - 20.0 // Look slightly down at the passing mountains

		// 7. Create the Master Camera for this specific frame
		masterCam := si3d.NewCamera(localCamX, cameraHeight, localCamZ, 0, 0, 0)
		masterCam.LookAt(si3d.NewVector3(lookAtX, lookAtY, lookAtZ), si3d.NewVector3(0, 1, 0))

		probe := NewProbe(starfieldPosition, masterCam)

		// 8. Render the Starfield Background
		field := universe.NewStarfield(probe.Camera, probe.StarfieldPosition)
		frameImg := field.GetStarField(HEIGHT, WIDTH) // Assuming GetStarField doesn't need w/h parameters if hardcoded

		// 9. Render the Mountains on top
		world := si3d.NewWorld3d()
		world.AddCamera(probe.Camera, localCamX, cameraHeight, localCamZ)
		world.AddObjectDrawFirst(&si3d.Entity{Model: mountains, X: 0, Y: 0, Z: 0})
		world.RenderToImage(frameImg)

		// 10. Convert the RGBA frame to a 256-color Paletted image for the GIF
		bounds := frameImg.Bounds()
		palettedImage := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, frameImg, bounds.Min, draw.Over)

		// 11. Append to the GIF sequence
		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 5) // 5 = 50ms delay per frame (20fps)

		fmt.Printf("Rendered frame %d/%d\n", i+1, frames)
	}

	// 12. Save the final animated GIF
	fmt.Println("Saving orbit.gif...")
	f, err := os.Create("orbit.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gif.EncodeAll(f, outGif)
	fmt.Println("Done! Open orbit.gif to see the result.")
}
