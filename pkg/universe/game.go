package universe

import (
	"math"

	"github.com/smasonuk/si3d/pkg/si3d"
)

// Constants for our 3-Tier coordinate system
const (
	MmPerAU = 149597870700000.0 // Exactly 1 AU in millimeters
	AUPerLY = 63241             // A standardized integer amount of AU in 1 LY
)

type GalacticPosition struct {
	SectorX, SectorY, SectorZ int64   // Macro:  Light-Years
	SystemX, SystemY, SystemZ int64   // Middle: Astronomical Units (AU)
	LocalX, LocalY, LocalZ    float64 // Micro:  Millimeters
}

func NewGalacticPosition(secX, secY, secZ, sysX, sysY, sysZ int64, locX, locY, locZ float64) *GalacticPosition {
	pos := &GalacticPosition{
		SectorX: secX, SectorY: secY, SectorZ: secZ,
		SystemX: sysX, SystemY: sysY, SystemZ: sysZ,
		LocalX: locX, LocalY: locY, LocalZ: locZ,
	}
	pos.Normalize()
	return pos
}

// Move translates the probe by a specific amount of millimeters.
func (pos *GalacticPosition) Move(deltaX, deltaY, deltaZ float64) {
	pos.LocalX += deltaX
	pos.LocalY += deltaY
	pos.LocalZ += deltaZ
	pos.Normalize()
}

// Normalize cascades coordinate overflows upward.
// Millimeters spill into AU, and AU spills into Light-Years.
func (pos *GalacticPosition) Normalize() {
	// 1. Roll Local (mm) up into System (AU)
	pos.LocalX, pos.SystemX = normalizeFloatTier(pos.LocalX, pos.SystemX, MmPerAU)
	pos.LocalY, pos.SystemY = normalizeFloatTier(pos.LocalY, pos.SystemY, MmPerAU)
	pos.LocalZ, pos.SystemZ = normalizeFloatTier(pos.LocalZ, pos.SystemZ, MmPerAU)

	// 2. Roll System (AU) up into Sector (LY)
	pos.SystemX, pos.SectorX = normalizeIntTier(pos.SystemX, pos.SectorX, AUPerLY)
	pos.SystemY, pos.SectorY = normalizeIntTier(pos.SystemY, pos.SectorY, AUPerLY)
	pos.SystemZ, pos.SectorZ = normalizeIntTier(pos.SystemZ, pos.SectorZ, AUPerLY)
}

// ToStarfieldPosition extracts the Macro and Middle tiers to generate the background.
// We combine LY and a fraction of AU so the starfield shifts smoothly as you traverse a system.
func (pos *GalacticPosition) ToStarfieldPosition() si3d.Vector3 {
	return si3d.NewVector3(
		float64(pos.SectorX)+(float64(pos.SystemX)/float64(AUPerLY)),
		float64(pos.SectorY)+(float64(pos.SystemY)/float64(AUPerLY)),
		float64(pos.SectorZ)+(float64(pos.SystemZ)/float64(AUPerLY)),
	)
}

// normalizeFloatTier centers a float64 inside a boundary and increments the parent int64
func normalizeFloatTier(micro float64, macro int64, limit float64) (float64, int64) {
	halfLimit := limit / 2.0
	if math.Abs(micro) > halfLimit {
		crossed := math.Round(micro / limit)
		macro += int64(crossed)
		micro -= crossed * limit
	}
	return micro, macro
}

// normalizeIntTier centers an int64 inside a boundary and increments the parent int64
func normalizeIntTier(micro int64, macro int64, limit int64) (int64, int64) {
	halfLimit := limit / 2
	if micro > halfLimit {
		crossed := (micro + halfLimit) / limit
		macro += crossed
		micro -= crossed * limit
	} else if micro < -halfLimit {
		crossed := (micro - halfLimit) / limit
		macro += crossed
		micro -= crossed * limit
	}
	return micro, macro
}
