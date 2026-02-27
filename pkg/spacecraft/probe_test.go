package spacecraft

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/smasonuk/unknowngalaxy/pkg/comms"
	"github.com/smasonuk/unknowngalaxy/pkg/universe"
)

func TestSpaceProbe_ReceivesMessageViaBus(t *testing.T) {
	bus := comms.NewMessageBus()
	pos := universe.NewGalacticPosition(0, 0, 0, 0, 0, 0, 0, 0, 0)
	scene := universe.NewLocalScene(0, 0, 0, 0, 0, 0)

	probe := NewSpaceProbe("Probe1", pos, scene, bus)

	payload := []byte("ping")
	bus.Send("Earth", "Probe1", payload)

	// Message must not be in the queue before Tick
	_, err := probe.VM.Disk.Read(".msgq.sys")
	if err == nil {
		t.Fatal("expected .msgq.sys to be absent before bus.Tick()")
	}

	bus.Tick()

	// After Tick the subscriber fires and PushMessage writes to .msgq.sys
	data, err := probe.VM.Disk.Read(".msgq.sys")
	if err != nil {
		t.Fatalf(".msgq.sys not found after bus.Tick(): %v", err)
	}

	// Queue format: [len: uint16 LE][payload bytes]
	if len(data) < 2+len(payload) {
		t.Fatalf("queue too short: got %d bytes, want at least %d", len(data), 2+len(payload))
	}
	storedLen := binary.LittleEndian.Uint16(data[0:2])
	if int(storedLen) != len(payload) {
		t.Errorf("stored length = %d, want %d", storedLen, len(payload))
	}
	if !bytes.Equal(data[2:2+len(payload)], payload) {
		t.Errorf("stored payload = %q, want %q", data[2:], payload)
	}
}
