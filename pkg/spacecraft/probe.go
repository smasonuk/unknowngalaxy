package spacecraft

import (
	_ "embed"
	"fmt"
	"image"

	"gocpu/pkg/compiler"
	"gocpu/pkg/cpu"
	"gocpu/pkg/peripherals"

	"github.com/smasonuk/unknowngalaxy/pkg/comms"
	"github.com/smasonuk/unknowngalaxy/pkg/universe"
)

//go:embed assets/probe_os.c
var probeOSSource string

// SpaceProbe bundles a physical probe, its virtual CPU, and the message receiver
// peripheral, wiring them together through the game's message bus.
type SpaceProbe struct {
	Physical    *universe.Probe
	VM          *cpu.CPU
	MsgReceiver *peripherals.MessageReceiver
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
		rgba, _ := img.(*image.RGBA)
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
