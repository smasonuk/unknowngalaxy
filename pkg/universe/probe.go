package universe

import (
	"image"
	"image/png"
	"os"

	"github.com/smasonuk/si3d/pkg/si3d"
)

const HEIGHT = 128
const WIDTH = 128

type Probe struct {
	ID       string
	Position *GalacticPosition
	Camera   *si3d.Camera
}

func NewProbe(id string, startPos *GalacticPosition) *Probe {
	return &Probe{
		ID:       id,
		Position: startPos,
		Camera:   si3d.NewCamera(0, 0, 0, 0, 0, 0), // Base camera, position updated by Scene
	}
}

type LocalScene struct {
	SectorX, SectorY, SectorZ int64
	SystemX, SystemY, SystemZ int64
	Entities                  []*si3d.Entity
}

func NewLocalScene(secX, secY, secZ, sysX, sysY, sysZ int64) *LocalScene {
	return &LocalScene{
		SectorX: secX, SectorY: secY, SectorZ: secZ,
		SystemX: sysX, SystemY: sysY, SystemZ: sysZ,
		Entities: make([]*si3d.Entity, 0),
	}
}

func (s *LocalScene) AddEntity(e *si3d.Entity) {
	s.Entities = append(s.Entities, e)
}

func (s *LocalScene) TakePicture(probe *Probe, width, height int) image.Image {
	starfieldPos := probe.Position.ToStarfieldPosition()
	field := NewStarfield(probe.Camera, starfieldPos)
	frameImg := field.GetStarField(height, width)

	world := si3d.NewWorld3d()
	world.AddCamera(probe.Camera, probe.Position.LocalX, probe.Position.LocalY, probe.Position.LocalZ)

	for _, entity := range s.Entities {
		world.AddObjectDrawFirst(entity)
	}

	world.RenderToImage(frameImg)
	return frameImg
}

// utils

func saveImageToFile(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}
}
func (p *Probe) PointCamera(target si3d.Vector3) {
	// 1. Re-center the internal si3d camera to the probe's actual current location
	p.Camera = si3d.NewCamera(p.Position.LocalX, p.Position.LocalY, p.Position.LocalZ, 0, 0, 0)

	// 2. NOW calculate the view angle to the target
	p.Camera.LookAt(target, si3d.NewVector3(0, 1, 0))
}
