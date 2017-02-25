package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"

	"github.com/pkg/term"
)

/**
 * Game
 */

type Game struct {
	level *Level
	pc    *PlayerCharacter
}

func (g *Game) generate_level() {
	g.level = &Level{60, 40, nil, nil}
	g.level.generate_rooms(3)
	g.pc = &PlayerCharacter{0, 0}
	g.level.put_player(g.pc)
}

func (g *Game) update() {
	fmt.Println("updating")
	g.level.render()
}

func (g *Game) handle_user_input(c string) {
	switch c {
	case "j":
		fmt.Println(g.pc.y)
		g.pc.y += 1
		fmt.Println(g.pc.y)
		fmt.Println("-")
	case "h":
		g.pc.x -= 1
	case "l":
		g.pc.x += 1
	case "k":
		g.pc.y -= 1
	}
}

func new_game() Game {
	g := Game{nil, nil}
	g.generate_level()
	return g
}

/**
 * Level
 */

type Level struct {
	width, height int
	rooms         []Room
	pc            *PlayerCharacter
}

func (l *Level) new_room() Room {
	w := get_rand(10) + 2
	h := get_rand(8) + 2
	x := get_rand(l.width - w)
	y := get_rand(l.height - h)
	return Room{x, y, w, h}
}
func (l *Level) generate_rooms(n int) {
	l.rooms = []Room{}
	for i := 0; i < n; i++ {
		l.rooms = append(l.rooms, l.new_room())
	}
}
func (l *Level) print_char(x int, y int) {
	dot := " "
	for _, r := range l.rooms {
		if r.exists_at(x, y) {
			if l.pc.x == x && l.pc.y == y {
				dot = "@"
			} else {
				dot = r.get_dot(x, y)
			}
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
func (l *Level) get_random_room() Room {
	i := get_rand(len(l.rooms))
	return l.rooms[i]
}
func (l *Level) put_player(pc *PlayerCharacter) {
	r := l.get_random_room()
	p := r.get_random_inner_point()
	pc.x = p.x
	pc.y = p.y
	l.pc = pc
}

/**
 * Room
 */

type Room struct {
	x, y int
	w, h int
}

func (r *Room) get_random_inner_point() Point {
	pts := r.get_inner_points()
	i := get_rand(len(pts))
	return pts[i]
}
func (r *Room) get_inner_points() []Point {
	a := Point{r.x + 1, r.y + 1}
	b := Point{r.x + r.w - 1, r.y + r.h - 1}
	return get_rect_points(a, b)
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

/**
 *
 */
type Point struct {
	x, y int
}

/**
 * Player Character
 */

type PlayerCharacter struct {
	x, y int
}

/**
 * main
 */
func get_rect_points(a Point, b Point) []Point {
	p := []Point{}
	for i := a.x; i <= b.x; i++ {
		for j := a.y; j <= b.y; j++ {
			p = append(p, Point{i, j})
		}
	}
	return p
}

func get_rand(i int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(i)
}

// since installing Goncurses on mac os x is a drag ...
func getch() []byte {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)
	numRead, err := t.Read(bytes)
	t.Restore()
	t.Close()
	if err != nil {
		return nil
	}
	return bytes[0:numRead]
}

func handle_io(g Game) bool {
	b := string(getch())
	if b == "c" {
		return true
	} else {
		g.handle_user_input(b)
	}
	return false
}

func main() {
	game := new_game()
	for {
		exec.Command("Clear").Run()
		game.update()
		if handle_io(game) {
			break
		}
	}
}
