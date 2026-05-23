/*
 * GOL, didactic project to learn go.
 * grid'll be a fixed size toroidal plane
 */

package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	void  = " "
	block = "██"
)

type Grid struct {
	w     int
	h     int
	cells []string
}

func newGrid(w, h int) *Grid {
	return &Grid{
		w:     w,
		h:     h,
		cells: make([]string, w*h),
	}
}

func (g *Grid) idx(x, y int) int {
	return g.w*y + x
}

func (g *Grid) cell(x, y int) string {
	return g.cells[g.idx(x, y)]
}

func (g *Grid) setCell(x, y int, v string) {
	g.cells[g.idx(x, y)] = v
}

func (g *Grid) plot() {
	for y := 0; y < g.h; y++ {
		ln := g.cells[(y * g.w):(y*g.w + g.w)]

		for _, v := range ln {
			if v != "" {
				fmt.Print(v)
			} else {
				fmt.Print(void)
			}
		}

		fmt.Println()
	}
}

func getTermSize() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Fatal error while reading terminal size")
		os.Exit(1)
	}

	return w / 2, h
}

func main() {
	w, h := getTermSize()

	g := newGrid(w, h)
	g.plot()
}
