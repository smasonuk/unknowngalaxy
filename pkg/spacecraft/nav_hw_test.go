package spacecraft

import (
	"testing"

	"github.com/smasonuk/unknowngalaxy/pkg/universe"
	"gocpu/pkg/cpu"
)

func TestNavigationPeripheral_Move(t *testing.T) {
	pos := universe.NewGalacticPosition(0, 0, 0, 0, 0, 0, 0, 0, 0)
	c := cpu.NewCPU()
	nav := NewNavigationPeripheral(c, 3, pos)
	c.MountPeripheral(3, nav)

	// Write thrust vector: X=100mm, Y=-50mm, Z=200mm
	dy := int16(-50)
	nav.Write16(0x02, uint16(int16(100)))
	nav.Write16(0x04, uint16(dy))
	nav.Write16(0x06, uint16(int16(200)))

	// Execute burn
	nav.Write16(0x00, 1)

	if pos.LocalX != 100.0 {
		t.Errorf("LocalX: want 100.0, got %v", pos.LocalX)
	}
	if pos.LocalY != -50.0 {
		t.Errorf("LocalY: want -50.0, got %v", pos.LocalY)
	}
	if pos.LocalZ != 200.0 {
		t.Errorf("LocalZ: want 200.0, got %v", pos.LocalZ)
	}
}

func TestNavigationPeripheral_InterruptFired(t *testing.T) {
	pos := universe.NewGalacticPosition(0, 0, 0, 0, 0, 0, 0, 0, 0)
	c := cpu.NewCPU()
	nav := NewNavigationPeripheral(c, 3, pos)
	c.MountPeripheral(3, nav)

	nav.Write16(0x00, 1)

	if c.PeripheralIntMask&(1<<3) == 0 {
		t.Error("expected peripheral interrupt bit for slot 3 to be set after burn")
	}
}

func TestNavigationPeripheral_ReadBack(t *testing.T) {
	pos := universe.NewGalacticPosition(0, 0, 0, 0, 0, 0, 0, 0, 0)
	c := cpu.NewCPU()
	nav := NewNavigationPeripheral(c, 0, pos)

	vx := int16(-10)
	nav.Write16(0x02, uint16(vx))
	nav.Write16(0x04, uint16(int16(42)))
	nav.Write16(0x06, uint16(int16(7)))

	if int16(nav.Read16(0x02)) != -10 {
		t.Errorf("0x02 readback: want -10, got %d", int16(nav.Read16(0x02)))
	}
	if int16(nav.Read16(0x04)) != 42 {
		t.Errorf("0x04 readback: want 42, got %d", int16(nav.Read16(0x04)))
	}
	if int16(nav.Read16(0x06)) != 7 {
		t.Errorf("0x06 readback: want 7, got %d", int16(nav.Read16(0x06)))
	}
}
