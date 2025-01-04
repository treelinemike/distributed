package mazeviz

import (
	"engg415/mazeio"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (p *Mazeparams) Setparams(m int, n int) {
	p.M = m
	p.N = n
	numcells := max(m, n)
	p.CSZ = int(math.Round(708.0 / (1.1*float64(numcells) + 0.1)))
	p.CSP = int(math.Floor(((708.0 - float64(numcells*p.CSZ)) / float64(numcells+1))))
	p.CSK = p.CSZ + p.CSP
	p.WW = p.CSP + n*p.CSK
	p.WH = p.CSP + m*p.CSK
}

type Game struct {
	image *ebiten.Image
	maze  *Maze
}

func (g *Game) Loadmaze(jsonfilename string) (ww int, wh int, err error) {

	// default values
	ww = 0
	wh = 0
	err = nil

	// load a maze configuration from json
	readmaze, err := mazeio.Readjsonmaze(jsonfilename)
	if err != nil {
		return
	}

	// set parameters
	p := new(Mazeparams)
	p.Setparams(int(readmaze.M), int(readmaze.N))
	ww = p.WW
	wh = p.WH

	// generate maze object
	mz, err := NewMaze(*p)
	if err != nil {
		return
	}

	// add data from json maze element
	for _, e := range readmaze.Elements {
		switch e.Type {
		case 0: // wall types
			for i, v := range e.Data {
				var wt Walltype
				switch v {
				case 0:
					wt = W_none
				case 1:
					wt = W_latent
				case 2:
					wt = W_observed
				case 3:
					wt = W_phantom
				}
				mz.walls[i].Type = wt
			}
		case 100: // cell types
			for i, v := range e.Data {
				var ct Celltype
				switch v {
				case 0:
					ct = C_none
				case 1:
					ct = C_goal
				case 2:
					ct = C_start
				}
				mz.cells[i].Type = ct
			}
		case 101: // cell values
			for i, v := range e.Data {
				mz.cells[i].Text = fmt.Sprintf("%d", int(v))
			}
		}
	}
	g.maze = mz
	return
}

func (g *Game) Savemaze(jsonfilename string) error {
	// generate a maze configuration
	writemaze := new(mazeio.Mazedata)
	writemaze.Title = "test output"
	writemaze.Author = "kokko"
	writemaze.Description = ""
	writemaze.M = int32(g.maze.p.M)
	writemaze.N = int32(g.maze.p.N)
	newelement := new(mazeio.Mazeelement)

	// wall types
	newelement.Type = 0
	newelement.Description = ""
	newelement.Data = []float32{}
	for _, w := range g.maze.walls {
		var wt float32
		switch w.Type {
		case W_none:
			wt = 0
		case W_latent:
			wt = 1
		case W_observed:
			wt = 2
		case W_phantom:
			wt = 3
		}
		newelement.Data = append(newelement.Data, wt)
	}
	writemaze.Elements = append(writemaze.Elements, *newelement)

	// cell types
	newelement.Type = 100
	newelement.Description = ""
	newelement.Data = []float32{}
	for _, c := range g.maze.cells {
		var ct float32
		switch c.Type {
		case C_none:
			ct = 0
		case C_goal:
			ct = 1
		case C_start:
			ct = 2
		}
		newelement.Data = append(newelement.Data, ct)
	}
	writemaze.Elements = append(writemaze.Elements, *newelement)

	// TODO: WRITE OUT CELL VALUES

	// write the maze configuration to json
	err := mazeio.Writejsonmaze(jsonfilename, *writemaze)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Newmaze(m, n int) (ww int, wh int, err error) {

	// default values
	ww = 0
	wh = 0
	err = nil

	p := new(Mazeparams)
	p.Setparams(m, n)
	mz, err := NewMaze(*p)
	if err != nil {
		return
	}
	ww = p.WW
	wh = p.WH
	g.maze = mz
	return
}

func Newgame() (*Game, error) {
	g := new(Game)
	g.image = new(ebiten.Image)
	return g, nil
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		// search walls first
		for i, l := range g.maze.walls {
			if (float32(x) >= l.Xmin) && (float32(x) <= l.Xmax) && (float32(y) >= l.Ymin) && (float32(y) <= l.Ymax) {
				switch l.Type {
				case W_none:
					g.maze.walls[i].Type = W_latent
				case W_latent:
					g.maze.walls[i].Type = W_observed
				case W_observed:
					g.maze.walls[i].Type = W_none
				}
				break
			}
		}

		// now search cells
		for i, c := range g.maze.cells {
			if (float32(x) > c.Xmin) && (float32(x) < c.Xmax) && (float32(y) > c.Ymin) && (float32(y) < c.Ymax) {
				if c.Type == C_none {
					g.maze.cells[i].Type = C_goal
				} else {
					g.maze.cells[i].Type = C_none
				}
				break
			}
		}

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.maze.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
}
