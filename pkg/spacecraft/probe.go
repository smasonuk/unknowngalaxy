package spacecraft

import (
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"gocpu/pkg/compiler"
	"gocpu/pkg/cpu"
	"gocpu/pkg/peripherals"

	"github.com/smasonuk/unknowngalaxy/pkg/comms"
	"github.com/smasonuk/unknowngalaxy/pkg/universe"
)

// +go:embed assets/basic_probe_os.c
//
//go:embed assets/probe_os.c
var probeOSSource string

// SpaceProbe bundles a physical probe, its virtual CPU, and the message receiver
// peripheral, wiring them together through the game's message bus.
type SpaceProbe struct {
	Physical    *universe.Probe
	VM          *cpu.CPU
	MsgReceiver *peripherals.MessageReceiver
}

func ConvertToRGBA(img image.Image) *image.RGBA {
	// 1. Check if it is already an RGBA image via type assertion.
	// If it is, we can just return it and save memory/processing time.
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	// 2. If it's not, get the original image's bounds
	bounds := img.Bounds()

	// 3. Create a new RGBA image with those bounds
	rgba := image.NewRGBA(bounds)

	// 4. Draw the original image onto the new RGBA image.
	// We use draw.Src to overwrite the pixels entirely.
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	return rgba
}

func WriteImageToFile(img image.Image, filename string) {
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

func WriteImageToFile2(img image.RGBA, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("EARTH: Failed to create %s: %v\n", filename, err)
		return
	}
	defer f.Close()
	if err := png.Encode(f, &img); err != nil {
		fmt.Printf("EARTH: Failed to encode image: %v\n", err)
	}
}

// NewSpaceProbe creates a SpaceProbe, mounts its peripherals, and subscribes to the bus.
func NewSpaceProbe(id string, startPos *universe.GalacticPosition, scene *universe.LocalScene, bus *comms.MessageBus) *SpaceProbe {
	physical := universe.NewProbe(id, startPos)
	vm := cpu.NewCPU(id)

	// Slot 0: Message Sender — outbound messages go onto the bus.
	dispatchFunc := func(target string, body []byte) {
		bus.Send(id, target, body)
	}
	vm.MountPeripheral(0, peripherals.NewMessageSender(vm, 0, dispatchFunc))

	// Slot 1: Camera — captures the probe's local scene view.
	captureFunc := func() *image.RGBA {
		img := scene.TakePicture(physical, 128, 128)
		WriteImageToFile(img, "1111.png")

		ConvertToRGBA(img)
		rgba := ConvertToRGBA(img)

		WriteImageToFile2(*rgba, "2222.png")

		return rgba
	}
	vm.MountPeripheral(1, peripherals.NewCameraPeripheral(vm, 1, captureFunc))

	// Slot 2: Message Receiver — inbound messages land in the VFS queue.
	msgReceiver := peripherals.NewMessageReceiver(vm, 2)
	vm.MountPeripheral(2, msgReceiver)

	sp := &SpaceProbe{
		Physical:    physical,
		VM:          vm,
		MsgReceiver: msgReceiver,
	}

	// Subscribe to the bus so incoming messages are pushed into the receiver.
	bus.Subscribe(id, func(m comms.Message) {
		_ = msgReceiver.PushMessage(m.SenderID, m.Payload)
	})

	// Compile and load the probe OS into VM memory.
	_, mc, err := compiler.Compile(probeOSSource, "")
	if err != nil {
		fmt.Printf("[SpaceProbe %s] OS compile error: %v\n", id, err)
	} else if len(mc) > len(vm.Memory) {
		fmt.Printf("[SpaceProbe %s] OS binary too large (%d bytes)\n", id, len(mc))
	} else {
		copy(vm.Memory[:], mc)
	}

	return sp
}

// Tick advances the VM by the given number of cycles.
func (sp *SpaceProbe) Tick(cycles int) {
	for i := 0; i < cycles; i++ {
		sp.VM.Step()
	}
}
