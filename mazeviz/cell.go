package mazeviz

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Celltype int

const (
	C_none Celltype = iota
	C_goal
	C_start
)

type Cell struct {
	X                      int // RELATIVE TO SCREEN (not maze, x/y sawpped in maze)
	Y                      int // RELATIVE TO SCREEN (not maze, x/y sawpped in maze)
	Xmin, Xmax, Ymin, Ymax float32
	CX                     float32
	CY                     float32
	Size                   int
	Color                  color.RGBA
	Text                   string
	Type                   Celltype
}

var labelfontsource *text.GoTextFaceSource
var labelfontface *text.GoTextFace

func celltextinit(fontsize float64) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	labelfontsource = s
	labelfontface = &text.GoTextFace{
		Source: labelfontsource,
		Size:   fontsize,
	}
}

func (c *Cell) Updpate(x int, y int, color color.RGBA, text string) {
	c.X = x
	c.Y = y
	c.Color = color
	c.Text = text
}

func (c *Cell) Draw(screen *ebiten.Image) {
	cell := ebiten.NewImage(c.Size, c.Size)
	drawopts := new(ebiten.DrawImageOptions)
	drawopts.GeoM.Translate(float64(c.X), float64(c.Y))
	switch c.Type {
	case C_goal:
		cell.Fill(color.RGBA{0x00, 0xa0, 0x00, 0xff})
	default:
		cell.Fill(color.RGBA{0x55, 0x55, 0x55, 0xff})
	}
	screen.DrawImage(cell, drawopts)

	w, h := text.Measure(c.Text, labelfontface, labelfontface.Size)
	x := float64(c.X) + (float64(c.Size)-float64(w))/2
	y := float64(c.Y) + (float64(c.Size)-float64(h))/2

	//vector.DrawFilledRect(screen, float32(x_, y, float32(w), float32(h), color.RGBA{0xff, 0xff, 0xff, 0xff}, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.RGBA{0xb0, 0xb0, 0xb0, 0xff})
	op.LineSpacing = labelfontface.Size
	text.Draw(screen, c.Text, labelfontface, op)

}
