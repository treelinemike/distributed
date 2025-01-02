package mazeviz

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Maze struct {
	cells []Cell
	lines []Line
}

func pixdist(x0, y0, x1, y1 float64) float64 {
	return math.Sqrt(math.Pow((x1-x0), 2.0) + math.Pow((y1-y0), 2))
}

func (m *Maze) Draw(screen *ebiten.Image) {
	for _, c := range m.cells {
		c.Draw(screen)
	}
	for _, l := range m.lines {
		l.Draw(screen)

		/*newsq := ebiten.NewImage(200, 200)
		drawopts := new(ebiten.DrawImageOptions)
		drawopts.GeoM.Translate(float64(100), float64(100))
		newsq.Fill(color.RGBA{0x00, 0xc0, 0x00, 0x50})
		screen.DrawImage(newsq, drawopts)
		*/
	}
}

func NewMaze(p Params) (*Maze, error) {
	m := new(Maze)

	// add cells to maze
	for cols := 0; cols < p.N; cols++ {
		for rows := p.M - 1; rows >= 0; rows-- {
			c := new(Cell)
			c.Size = p.CSZ
			c.X = p.CSP + p.CSK*cols
			c.Y = p.CSP + p.CSK*rows
			c.CX = float64(p.CSP) + float64(p.CSK)*(float64(cols)+0.5)
			c.CY = float64(p.CSP) + float64(p.CSK)*(float64(rows)+0.5)
			c.Color = color.RGBA{0x50, 0x50, 0x50, 0xff}
			m.cells = append(m.cells, *c)
		}
	}

	// add lines to maze
	l := new(Line)
	l.X0 = float32(p.CSP) / 2
	l.Y0 = float32(p.CSP) / 2
	l.X1 = float32(p.CSP) / 2
	l.Y1 = float32(p.CSP)/2 + float32(p.CSK)
	l.Width = float32(p.CSP)
	m.lines = append(m.lines, *l)

	//for wallidx := 0; wallidx < (2*p.M+1)*p.N+p.M; wallidx++ {

	//	}

	// return maze
	return m, nil
}
