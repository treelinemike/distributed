package mazeviz

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Line struct {
	X0, Y0, X1, Y1 float32
	Width          float32
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func (l *Line) Draw(screen *ebiten.Image) {

	whiteImage.Fill(color.White)

	var path vector.Path
	path.MoveTo(l.X0, l.Y0)
	path.LineTo(l.X1, l.Y1)
	//fmt.Printf("(%0.2f, %0.2f) to (%0.2f, %0.2f)\n", l.X0, l.Y0, l.X1, l.Y1)

	op := &vector.StrokeOptions{}
	op.LineCap = vector.LineCapRound
	op.LineJoin = vector.LineJoinRound
	op.Width = l.Width
	vs, is := path.AppendVerticesAndIndicesForStroke([]ebiten.Vertex{}, []uint16{}, op)

	//vs, is := path.AppendVerticesAndIndicesForStroke(m.vertices[:0], m.indices[:0], op)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 0xff
		vs[i].ColorG = 0x00
		vs[i].ColorB = 0x00
		vs[i].ColorA = 1
	}
	screen.DrawTriangles(vs, is, whiteSubImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: false,
	})

}
