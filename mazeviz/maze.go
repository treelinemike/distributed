package mazeviz

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Interval struct {
	Start float32
	End   float32
	Index int
}

type Maze struct {
	cells []Cell
	lines []Line
}

func (m *Maze) Draw(screen *ebiten.Image) {

	// draw cells first
	for _, c := range m.cells {
		c.Draw(screen)
	}

	// draw non-existant lines next
	for _, l := range m.lines {
		if l.Type == W_none {
			l.Draw(screen)
		}
	}

	// draw latent lines
	for _, l := range m.lines {
		if l.Type == W_latent {
			l.Draw(screen)
		}
	}

	// finally draw true lines
	for _, l := range m.lines {
		if l.Type == W_observed {
			l.Draw(screen)
		}
	}
}

func NewMaze(p Params) (*Maze, error) {
	m := new(Maze)

	celltextinit(math.Round(float64(p.CSZ) / 3))

	for col := 0; col <= p.N; col++ {

		// add column of vertical lines
		y0 := float32(p.CSP)/2 + float32(p.M*p.CSK)
		for y := y0; y > float32(p.CSP); y -= float32(p.CSK) {
			l := new(Line)
			l.X0 = float32(p.CSP)/2 + float32(col*p.CSK)
			l.Y0 = float32(y)
			l.X1 = l.X0
			l.Y1 = l.Y0 - float32(p.CSK)
			l.Xmin = l.X0 - float32(p.CSP)/2
			l.Xmax = l.X0 + float32(p.CSP)/2
			l.Ymin = l.Y1 + float32(p.CSP)/2
			l.Ymax = l.Y0 - float32(p.CSP)/2
			l.Width = float32(p.CSP)
			l.Type = W_none
			m.lines = append(m.lines, *l)
		}

		// add column of horizontal lines and cells
		if col < p.N {
			for row := 0; row <= p.M; row++ {
				// horizontal line
				l := new(Line)
				l.X0 = float32(p.CSP)/2 + float32(col*p.CSK)
				l.Y0 = float32(p.CSP)/2 + float32(p.M-row)*float32(p.CSK)
				l.X1 = l.X0 + float32(p.CSK)
				l.Y1 = l.Y0
				l.Xmin = l.X0 + float32(p.CSP)/2
				l.Xmax = l.X1 - float32(p.CSP)/2
				l.Ymin = l.Y0 - float32(p.CSP)/2
				l.Ymax = l.Y0 + float32(p.CSP)/2
				l.Width = float32(p.CSP)
				l.Type = W_none
				m.lines = append(m.lines, *l)

				// cell on top of the line
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

	// add lines from JSON
	if len(p.Walltypes) == len(m.lines) {
		for i, w := range p.Walltypes {
			m.lines[i].Type = w
		}
	}

	// add cell types from JSON
	if len(p.Cellvals) == len(m.cells) {
		for i, cv := range p.Cellvals {
			m.cells[i].Text = fmt.Sprintf("%d", int(cv))
		}
	}

	// add cell values from JSON
	if len(p.Celltypes) == len(m.cells) {
		for i, ct := range p.Celltypes {
			m.cells[i].Type = ct
		}
	}

	// return maze
	return m, nil
}
