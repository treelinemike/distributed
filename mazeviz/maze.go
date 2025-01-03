package mazeviz

import (
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

func pixdist(x0, y0, x1, y1 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x1-x0), 2.0) + math.Pow(float64(y1-y0), 2)))
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

	// finally draw true lines
	for _, l := range m.lines {
		if l.Type == W_true {
			l.Draw(screen)
		}
	}
}

func NewMaze(p Params) (*Maze, error) {
	m := new(Maze)

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
			if col == 0 || col == p.N {
				l.Type = W_true
			} else {
				l.Type = W_none
			}
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
				if row == 0 || row == p.M {
					l.Type = W_true
				} else {
					l.Type = W_none
				}
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
					m.cells = append(m.cells, *c)
				}

			}
		}
	}

	// return maze
	return m, nil
}
