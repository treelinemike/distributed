package mazeviz

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Walltype int

const (
	W_none Walltype = iota
	W_true
	W_observed
)

type Line struct {
	X0, Y0, X1, Y1         float32 // line endpoints
	Xmin, Xmax, Ymin, Ymax float32
	Width                  float32
	Type                   Walltype
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

	op := &vector.StrokeOptions{}
	op.LineCap = vector.LineCapRound
	op.LineJoin = vector.LineJoinRound
	op.Width = l.Width
	vs, is := path.AppendVerticesAndIndicesForStroke([]ebiten.Vertex{}, []uint16{}, op)

	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		switch l.Type {
		case W_true:
			vs[i].ColorR = 0.85
			vs[i].ColorG = 0
			vs[i].ColorB = 0
		default:
			vs[i].ColorR = 0.5
			vs[i].ColorG = 0.5
			vs[i].ColorB = 0.5
		}
		vs[i].ColorA = 1
	}
	screen.DrawTriangles(vs, is, whiteSubImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: false,
	})

}
