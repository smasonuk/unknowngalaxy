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

	// Queue format: [SenderLen: uint8][SenderStr][BodyLen: uint16][Body]
	if len(data) < 1 {
		t.Fatalf("queue too short: got %d bytes, want at least 1", len(data))
	}
	
	senderLen := int(data[0])
	expectedSender := "Earth"
	if senderLen != len(expectedSender) {
		t.Errorf("stored sender length = %d, want %d", senderLen, len(expectedSender))
	}
	
	if len(data) < 1+senderLen+2+len(payload) {
		t.Fatalf("queue too short for full payload")
	}

	storedSender := string(data[1 : 1+senderLen])
	if storedSender != expectedSender {
		t.Errorf("stored sender = %q, want %q", storedSender, expectedSender)
	}

	storedLen := binary.LittleEndian.Uint16(data[1+senderLen : 1+senderLen+2])
	if int(storedLen) != len(payload) {
		t.Errorf("stored length = %d, want %d", storedLen, len(payload))
	}
	
	storedPayload := data[1+senderLen+2 : 1+senderLen+2+int(storedLen)]
	if !bytes.Equal(storedPayload, payload) {
		t.Errorf("stored payload = %q, want %q", storedPayload, payload)
	}
}
