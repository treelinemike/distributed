package mazeviz

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Cell struct {
	X     int // RELATIVE TO SCREEN (not maze, x/y sawpped in maze)
	Y     int // RELATIVE TO SCREEN (not maze, x/y sawpped in maze)
	CX    float64
	CY    float64
	Size  int
	Color color.RGBA
	Text  string
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
	cell.Fill(c.Color)
	screen.DrawImage(cell, drawopts)
}
