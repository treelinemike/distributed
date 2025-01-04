package mazeviz

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Mazeparams struct {
	M   int
	N   int
	CSZ int
	CSP int
	CSK int
	WW  int
	WH  int
}

type Interval struct {
	Start float32
	End   float32
	Index int
}

type Maze struct {
	p     Mazeparams
	cells []Cell
	walls []Wall
}

func (m *Maze) Draw(screen *ebiten.Image) {

	// draw cells first
	for _, c := range m.cells {
		c.Draw(screen)
	}

	// draw non-existant walls next
	for _, w := range m.walls {
		if w.Type == W_none {
			w.Draw(screen)
		}
	}

	// draw latent walls
	for _, w := range m.walls {
		if w.Type == W_latent {
			w.Draw(screen)
		}
	}

	// finally draw true walls that have been observed
	for _, w := range m.walls {
		if w.Type == W_observed {
			w.Draw(screen)
		}
	}

	// TODO: show phantom walls?
}

func NewMaze(p Mazeparams) (*Maze, error) {
	m := new(Maze)

	celltextinit(math.Round(float64(p.CSZ) / 3))
	m.p = p

	for col := 0; col <= p.N; col++ {

		// add column of vertical walls
		y0 := float32(p.CSP)/2 + float32(p.M*p.CSK)
		for y := y0; y > float32(p.CSP); y -= float32(p.CSK) {
			w := new(Wall)
			w.X0 = float32(p.CSP)/2 + float32(col*p.CSK)
			w.Y0 = float32(y)
			w.X1 = w.X0
			w.Y1 = w.Y0 - float32(p.CSK)
			w.Xmin = w.X0 - float32(p.CSP)/2
			w.Xmax = w.X0 + float32(p.CSP)/2
			w.Ymin = w.Y1 + float32(p.CSP)/2
			w.Ymax = w.Y0 - float32(p.CSP)/2
			w.Width = float32(p.CSP)
			w.Type = W_none
			m.walls = append(m.walls, *w)
		}

		// add column of horizontal lines and cells
		if col < p.N {
			for row := 0; row <= p.M; row++ {
				// horizontal walls
				w := new(Wall)
				w.X0 = float32(p.CSP)/2 + float32(col*p.CSK)
				w.Y0 = float32(p.CSP)/2 + float32(p.M-row)*float32(p.CSK)
				w.X1 = w.X0 + float32(p.CSK)
				w.Y1 = w.Y0
				w.Xmin = w.X0 + float32(p.CSP)/2
				w.Xmax = w.X1 - float32(p.CSP)/2
				w.Ymin = w.Y0 - float32(p.CSP)/2
				w.Ymax = w.Y0 + float32(p.CSP)/2
				w.Width = float32(p.CSP)
				w.Type = W_none
				m.walls = append(m.walls, *w)

				// cell on top of the wall
				if row < p.M {
					c := new(Cell)
					c.Size = p.CSZ
					c.X = p.CSP + p.CSK*col
					c.Y = p.CSP + (p.M-row-1)*(p.CSK)
					c.CX = float32(c.X) + float32(c.Size)/2
					c.CY = float32(c.Y) + float32(c.Size)/2
					c.Xmin = float32(c.X)
					c.Xmax = c.Xmin + float32(c.Size)
					c.Ymin = float32(c.Y)
					c.Ymax = c.Ymin + float32(c.Size)
					c.Color = color.RGBA{0x50, 0x50, 0x50, 0xff}
					c.Type = C_none
					m.cells = append(m.cells, *c)
				}
			}
		}
	}

	// return maze
	return m, nil
}
