package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

/**
 * Level
 */

type Level struct {
	width, height int
	rooms         []Room
}

func (l *Level) get_rand(i int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(i)
}
func (l *Level) new_room() Room {
	w := l.get_rand(10) + 2
	h := l.get_rand(8) + 2
	x := l.get_rand(l.width - w)
	y := l.get_rand(l.height - h)
	return Room{x, y, w, h}
}
func (l *Level) generate_rooms(n int) {
	for i := 0; i < n; i++ {
		l.rooms = append(l.rooms, l.new_room())
	}
}
func (l *Level) print_char(x int, y int) {
	dot := " "
	for _, r := range l.rooms {
		if r.exists_at(x, y) {
			dot = r.get_dot(x, y)
		}
	}
	fmt.Print(dot)
}
func (l *Level) render() {
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			l.print_char(x, y)
		}
		fmt.Println()
	}
}

func NewLevel() {
}

/**
 * Room
 */

type Room struct {
	x, y int
	w, h int
}

func (r *Room) exists_at(x int, y int) bool {
	xw := r.x + r.w
	yh := r.y + r.h
	return x >= r.x && x <= xw && y >= r.y && y <= yh
}
func (r *Room) is_border(x int, y int) bool {
	xw := r.x + r.w
	yh := r.y + r.h
	return x == r.x || x == xw || y == r.y || y == yh
}
func (r *Room) get_dot(x int, y int) string {
	if r.is_border(x, y) {
		return "#"
	}
	return "."
}

func generate_level() Level {
	l := Level{60, 40, []Room{}}
	l.generate_rooms(3)
	return l
}

func draw() {
	l := generate_level()
	l.render()
}

func scan() bool {
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	s := scan.Text()
	if strings.Contains(s, "c") {
		return true
	}
	return false
}

func main() {
	for {
		draw()
		if scan() {
			break
		}
	}
}
