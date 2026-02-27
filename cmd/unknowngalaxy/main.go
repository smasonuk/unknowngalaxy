package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smasonuk/si3d/pkg/si3d"
	"github.com/smasonuk/unknowngalaxy/pkg/comms"
	"github.com/smasonuk/unknowngalaxy/pkg/spacecraft"
	"github.com/smasonuk/unknowngalaxy/pkg/universe"
)

const probeID = "Voyager-1"

func saveImageToFile(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("EARTH: Failed to create %s: %v\n", filename, err)
		return
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		fmt.Printf("EARTH: Failed to encode image: %v\n", err)
	}
}

func handleEarthMessage(msg comms.Message) {
	if len(msg.Payload) == 16384 {
		img, err := comms.DecodeRGB332(msg.Payload, 128, 128)
		if err != nil {
			fmt.Printf("EARTH: Failed to decode image from %s: %v\n", msg.SenderID, err)
			return
		}
		filename := fmt.Sprintf("%s_capture.png", msg.SenderID)
		saveImageToFile(img, filename)
		fmt.Printf("EARTH: Received image from %s, saved to disk.\n", msg.SenderID)
	}
}

func main() {
	// Scene
	mountains := si3d.NewSubdividedPlaneHeightMapPerlin(
		10000, 10000,
		color.RGBA{R: 153, G: 196, B: 210, A: 255},
		35, 800, 800, 42,
	)
	// mountains.SetDrawLinesOnly(false)

	scene := universe.NewLocalScene(10000, 25000, 35000, 0, 0, 0)
	scene.AddEntity(&si3d.Entity{Model: mountains, X: 0, Y: 0, Z: 0})

	// Message bus
	bus := comms.NewMessageBus()
	bus.Subscribe("Earth", handleEarthMessage)

	// Probe
	startPos := universe.NewGalacticPosition(10000, 25000, 35000, 0, 0, 0, 0, -200.0, -400.0)
	probe := spacecraft.NewSpaceProbe(probeID, startPos, scene, bus)

	lookAtTarget := si3d.NewVector3(0, -200.0, 0)
	probe.Physical.PointCamera(lookAtTarget)

	// Graceful shutdown on SIGINT / SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(time.Millisecond * 16) // ~60 Hz
	defer ticker.Stop()

	fmt.Println("Simulation running. Press Ctrl+C to stop.")

	for {
		select {
		case <-stop:
			fmt.Println("Shutting down.")
			return
		case <-ticker.C:
			bus.Tick()
			probe.Tick(1000)
		}
	}
}
