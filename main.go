package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/term"
)

type CellState int

const (
	CellDead CellState = iota
	CellAlive
)

type Grid struct {
	w     int
	h     int
	cells []CellState
}

func wrapCoord(x, size int) int {
	return ((x % size) + size) % size
}

func newGrid(w, h int) *Grid {
	return &Grid{
		w:     w,
		h:     h,
		cells: make([]CellState, w*h),
	}
}

func (g *Grid) idx(x, y int) int {
	return g.w*wrapCoord(y, g.h) + wrapCoord(x, g.w)
}

func (g *Grid) cell(x, y int) CellState {
	return g.cells[g.idx(x, y)]
}

func (g *Grid) setCell(x, y int, v CellState) {
	g.cells[g.idx(x, y)] = v
}

func (g *Grid) plot(buf *bytes.Buffer) {
	fmt.Fprint(buf, "\033[H")
	for y := 0; y < g.h; y++ {
		ln := g.cells[(y * g.w):(y*g.w + g.w)]

		for _, v := range ln {
			render := ""
			if v == CellDead {
				render = "  "
			} else {
				render = "██"
			}
			fmt.Fprint(buf, render)
		}
		fmt.Fprint(buf, "\n")
	}
}

func (g *Grid) cellAliveNeighbours(x, y int) int {
	res := 0
	for xo := -1; xo <= 1; xo++ {
		for yo := -1; yo <= 1; yo++ {
			if xo == 0 && yo == 0 {
				continue
			}
			if g.cell(x+xo, y+yo) == CellAlive {
				res++
			}
		}
	}

	return res
}

func (g *Grid) cellNextState(x, y int) CellState {
	cs := g.cell(x, y)
	an := g.cellAliveNeighbours(x, y)

	if cs == CellDead {
		if an == 3 {
			return CellAlive
		}
		return CellDead
	}

	switch {
	case an < 2:
		return CellDead
	case an == 2 || an == 3:
		return CellAlive
	default: // an > 3
		return CellDead
	}
}

func (g *Grid) randomize() {
	for i := range g.cells {
		if rand.Intn(10) == 0 {
			g.cells[i] = CellAlive
		}
	}
}

func nextGrid(oldG *Grid) *Grid {
	newG := newGrid(oldG.w, oldG.h)
	for x := 0; x < oldG.w; x++ {
		for y := 0; y < oldG.h; y++ {
			newG.setCell(x, y, oldG.cellNextState(x, y))
		}
	}

	return newG
}

func main() {
	mustAtoi := func(s string) int {
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Error: invalid arguments")
			os.Exit(1)
		}
		return n
	}

	usageAndExit := func() {
		fmt.Println("Usage: gol <size w h | term>")
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		usageAndExit()
	}

	var w, h int
	switch args[0] {
	case "size":
		if len(args) != 3 {
			usageAndExit()
		}
		w = mustAtoi(args[1])
		h = mustAtoi(args[2])
	case "term":
		var err error
		w, h, err = term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			fmt.Println("Error: couldn't read terminal size")
			os.Exit(1)
		}
		w /= 2
	default:
		usageAndExit()
	}

	if w < 3 || h < 3 {
		fmt.Println("Error: grid must be at least 3x3")
		os.Exit(1)
	}
	g := newGrid(w, h)
	g.randomize()

	fmt.Print("\033[2J")
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	buf := new(bytes.Buffer)
	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()

	for range ticker.C {
		buf.Reset()
		g.plot(buf)
		_, err := buf.WriteTo(os.Stdout)
		if err != nil {
			fmt.Print("Error: problem while rendering")
			os.Exit(1)
		}
		g = nextGrid(g)
	}
}
