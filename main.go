package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/term"
)

type tiletype int

/**
 * Game
 */

type Game struct {
	level *Level
	pc    *PlayerCharacter
}

func (g *Game) handle_io() bool {
	b := string(getch())
	if b == "c" {
		return true
	} else {
		g.handle_user_input(b)
	}
	return false
}
func (g *Game) generate_level() {
	g.level = &Level{60, 32, nil, nil}
	g.level.init()
	g.pc = &PlayerCharacter{Point{0, 0}}
	g.level.put_player(g.pc)
}
func (g *Game) handle_user_input(c string) {
	x := g.pc.pos.x
	y := g.pc.pos.y
	switch c {
	case "j":
		y += 1
	case "h":
		x -= 1
	case "l":
		x += 1
	case "k":
		y -= 1
	}

	new_pc_pos := Point{x, y}
	wkbl := g.level.get_walkable_points()
	if new_pc_pos.is_in_slice(wkbl) {
		g.pc.pos = new_pc_pos
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

func (l *Level) init() {
	l.generate_rooms()
	l.generate_corridors()
}
func (l *Level) get_walkable_points() []Point {
	pts := []Point{}
	for _, r := range l.rooms {
		pts = append(pts, r.get_inner_points()...)
	}
	return pts
}
func (l *Level) new_room(a, b Point) Room {
	dx := b.x - a.x
	dy := b.y - a.y
	w := get_rand_range(3, dx)
	h := get_rand_range(3, dy)
	x := a.x + get_rand(dx-w)
	y := a.y + get_rand(dy-h)
	return Room{x, y, w, h}
}
func (l *Level) generate_rooms() {
	l.rooms = []Room{}
	row_height := l.height / 2
	col_width := l.width / 3
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			a := Point{j * col_width, i * row_height}
			b := Point{(j + 1) * col_width, (i + 1) * row_height}
			l.rooms = append(l.rooms, l.new_room(a, b))
		}
	}
}
func (l *Level) generate_corridors() {
	// TODO implement ...
}
func (l *Level) get_tile(p Point) tiletype {
	for _, r := range l.rooms {
		if r.exists_at(p) {
			if l.pc.pos.x == p.x && l.pc.pos.y == p.y {
				return 5
			} else {
				return r.get_tile(p)
			}
		}
	}
	return 0
}
func (l *Level) get_tiles() []tiletype {
	tiles := []tiletype{}
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			tiles = append(tiles, l.get_tile(Point{x, y}))
		}
	}
	return tiles
}
func (l *Level) get_random_room() Room {
	i := get_rand(len(l.rooms))
	return l.rooms[i]
}
func (l *Level) put_player(pc *PlayerCharacter) {
	r := l.get_random_room()
	p := r.get_random_inner_point()
	pc.pos = p
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
func (r *Room) get_points() []Point {
	a := Point{r.x, r.y}
	b := Point{r.x + r.w, r.y + r.h}
	return get_rect_points(a, b)
}
func (r *Room) is_my_point(p Point) bool {
	pts := r.get_points()
	return p.is_in_slice(pts)
}
func (r *Room) is_my_inner_point(p Point) bool {
	pts := r.get_inner_points()
	return p.is_in_slice(pts)
}
func (r *Room) exists_at(p Point) bool {
	return r.is_my_point(p)
}
func (r *Room) is_border(p Point) bool {
	return r.is_my_point(p) && !r.is_my_inner_point(p)
}
func (r *Room) get_tile(p Point) tiletype {
	if r.is_border(p) {
		return 2
	}
	return 1
}

/**
 *
 */
type Point struct {
	x, y int
}

func (p *Point) is_in_slice(s []Point) bool {
	for _, sp := range s {
		if sp.x == p.x && sp.y == p.y {
			return true
		}
	}
	return false
}

/**
 * Player Character
 */

type PlayerCharacter struct {
	pos Point
}

/**
 * renderer
 */
type Renderer struct{}

func new_renderer() *Renderer {
	return &Renderer{}
}
func (r *Renderer) clear() {
	fmt.Println("\033[H\033[2J")
}
func (r *Renderer) get_texture(tt tiletype) string {
	switch tt {
	case 1:
		return "."
	case 2:
		return "#"
	case 5:
		return "@"
	}
	return " "
}
func (r *Renderer) render(g *Game) {
	w := g.level.width
	h := g.level.height
	t := g.level.get_tiles()
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			indx := i*w + j
			texture := r.get_texture(t[indx])
			fmt.Print(texture)
		}
		fmt.Println()
	}
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

func get_rand_range(min int, max int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(max-min) + min
}
func get_rand(i int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(i)
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

func main() {
	renderer := new_renderer()
	game := new_game()
	for {
		renderer.clear()
		renderer.render(&game)
		if game.handle_io() {
			break
		}
	}
}
