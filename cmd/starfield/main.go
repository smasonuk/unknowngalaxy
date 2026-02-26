package main

import (
	"image/png"
	"os"

	"github.com/smasonuk/si3d/pkg/si3d"
	"github.com/smasonuk/unknowngalaxy/pkg/universe"
)

func main() {
	probePos := si3d.NewVector3(10000, 25000, 35000)
	cam := si3d.NewCamera(probePos.X, probePos.Y, probePos.Z, 0, 0, 0)
	cam.LookAt(si3d.NewVector3(0, 0, 0), si3d.NewVector3(0, 1, 0))

	field := universe.NewStarfield(cam, probePos)
	snapshot := field.GetStarField(512, 512)

	filename := ".temp.png"
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	png.Encode(f, snapshot)
	f.Close()

}
