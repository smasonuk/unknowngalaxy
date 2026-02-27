package spacecraft

import (
	"github.com/smasonuk/unknowngalaxy/pkg/universe"
	"gocpu/pkg/cpu"
)

const NavigationPeripheralType = "NavigationPeripheral"

// NavigationPeripheral exposes a simple burn interface to the virtual CPU.
// Writing thrust vectors (int16 mm) to 0x02/0x04/0x06 and issuing command 1
// to 0x00 translates them into a Move() call on the physical GalacticPosition.
type NavigationPeripheral struct {
	c    *cpu.CPU
	slot uint8
	pos  *universe.GalacticPosition
	dx   int16
	dy   int16
	dz   int16
}

func NewNavigationPeripheral(c *cpu.CPU, slot uint8, pos *universe.GalacticPosition) *NavigationPeripheral {
	return &NavigationPeripheral{c: c, slot: slot, pos: pos}
}

func (n *NavigationPeripheral) Type() string { return NavigationPeripheralType }

func (n *NavigationPeripheral) Read16(offset uint16) uint16 {
	if offset >= 0x08 && offset <= 0x0E {
		return cpu.EncodePeripheralName("NAVSYS", offset)
	}
	switch offset {
	case 0x00:
		return 0
	case 0x02:
		return uint16(n.dx)
	case 0x04:
		return uint16(n.dy)
	case 0x06:
		return uint16(n.dz)
	}
	return 0
}

func (n *NavigationPeripheral) Write16(offset uint16, val uint16) {
	switch offset {
	case 0x00:
		if val == 1 {
			n.pos.Move(float64(n.dx), float64(n.dy), float64(n.dz))
			n.c.TriggerPeripheralInterrupt(n.slot)
		}
	case 0x02:
		n.dx = int16(val)
	case 0x04:
		n.dy = int16(val)
	case 0x06:
		n.dz = int16(val)
	}
}

func (n *NavigationPeripheral) Step() {}
