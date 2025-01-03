package mazeviz

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Celltype int

const (
	C_none Celltype = iota
	C_start
	C_goal
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
}
