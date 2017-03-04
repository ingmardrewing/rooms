package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/bradfitz/slice" // for sorting
	"github.com/pkg/term"       // for getting userinput without the user hitting enter
)

/**
 * tiles
 */
type tiletype int

const (
	VoidTile          = iota
	FloorTile         = iota
	WallTile          = iota
	DoorTile          = iota
	PlayerTile        = iota
	StaircaseDownTile = iota
	StaircaseUpTile   = iota
)

var down = false

/**
 * Game
 */

type Game struct {
	level *Level
	pc    *PlayerCharacter
}

func (g *Game) status() string {
	return "h to go left, j to go down, k to go up, l to go right.\nq to quit"
}
func (g *Game) handle_io() bool {
	b := string(getch())
	if b == "q" {
		return true
	}
	if b == "c" {
		return true
	} else {
		g.handle_user_input(b)
	}
	return false
}
func (g *Game) generate_level() {
	g.level = &Level{60, 32, nil, nil, nil, nil, nil}
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
	case "d":
		// TODO implement transition to next level
		if g.level.is_staircase(g.pc.pos) {
			down = true
		}
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
	corridors     []Corridor
	doors         []Door
	staircases    []Staircase
	pc            *PlayerCharacter
}

func (l *Level) init() {
	l.generate_rooms()
	l.generate_corridors()
	l.generate_doors()
	l.generate_staircases()
}
func (l *Level) get_walkable_points() []Point {
	pts := []Point{}
	for _, r := range l.rooms {
		pts = append(pts, r.get_inner_points()...)
	}
	for _, c := range l.corridors {
		pts = append(pts, c.get_points()...)
	}
	return pts
}
func (l *Level) generate_rooms() {
	l.rooms = []Room{}
	row_height := l.height / 2
	col_width := l.width / 3
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			a := Point{j*col_width + j, i*row_height + i}
			b := Point{(j + 1) * col_width, (i + 1) * row_height}
			l.rooms = append(l.rooms, new_room(a, b))
		}
	}
}
func (l *Level) generate_corridors() {
	from := []int{0, 1, 3, 4, 0, 1, 2}
	to := []int{1, 2, 4, 5, 3, 4, 5}
	crs := []Corridor{}
	for i, _ := range from {
		crs = append(crs, new_corridor(l, from[i], to[i]))
	}
	l.corridors = crs
}
func (l *Level) generate_doors() {
	wall_pts := l.get_all_wall_points()
	corr_pts := l.get_all_corridor_points()
	door_pts := l.find_point_set_intersections(wall_pts, corr_pts)
	l.doors = l.generate_doors_at(door_pts)
}
func (l *Level) generate_doors_at(pts []Point) []Door {
	doors := []Door{}
	for _, p := range pts {
		doors = append(doors, new_door(p))
	}
	return doors
}
func (l *Level) generate_staircases() {
	s := new_staircase(l.get_random_room().get_random_inner_point())
	l.staircases = []Staircase{s}
}
func (l *Level) is_staircase(p Point) bool {
	for _, s := range l.staircases {
		if s.pos.equals(p) {
			return true
		}
	}
	return false
}
func (l *Level) get_all_wall_points() []Point {
	pts := []Point{}
	for _, r := range l.rooms {
		pts = append(pts, r.get_wall_points()...)
	}
	return pts
}
func (l *Level) get_all_corridor_points() []Point {
	pts := []Point{}
	for _, c := range l.corridors {
		pts = append(pts, c.get_points()...)
	}
	return pts
}
func (l *Level) find_point_set_intersections(s1 []Point, s2 []Point) []Point {
	pts := []Point{}
	for _, p1 := range s1 {
		if p1.is_in_slice(s2) {
			pts = append(pts, p1)
		}
	}
	return pts
}

func (l *Level) get_tile(p Point) tiletype {
	if l.pc.pos.x == p.x && l.pc.pos.y == p.y {
		return PlayerTile
	}
	// TODO use a map to manage points / tiles
	for _, d := range l.doors {
		if d.exists_at(p) {
			return d.get_tile(p)
		}
	}
	for _, s := range l.staircases {
		if s.exists_at(p) {
			return s.get_tile(p)
		}
	}
	for _, c := range l.corridors {
		if c.exists_at(p) {
			return c.get_tile(p)
		}
	}
	for _, r := range l.rooms {
		if r.exists_at(p) {
			return r.get_tile(p)
		}
	}
	return VoidTile
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
func (l *Level) get_random_room() *Room {
	i := get_rand(len(l.rooms))
	r := l.rooms[i]
	return &r
}
func (l *Level) put_player(pc *PlayerCharacter) {
	p := l.get_random_room().get_random_inner_point()
	pc.pos = p
	l.pc = pc
}

/**
 * Room
 */

type Room struct {
	x, y         int
	w, h         int
	points       []Point
	inner_points []Point
}

func new_room(a, b Point) Room {
	dx := b.x - a.x
	dy := b.y - a.y
	w := get_rand_range(3, dx-1)
	h := get_rand_range(3, dy-1)
	x := a.x + get_rand(dx-w)
	y := a.y + get_rand(dy-h)
	r := Room{x, y, w, h, nil, nil}
	r.init()
	return r
}
func (r *Room) init() {
	r.points = r.get_points()
	r.inner_points = r.get_inner_points()
}
func (r *Room) get_random_inner_point() Point {
	pts := r.get_inner_points()
	i := get_rand(len(pts))
	return pts[i]
}
func (r *Room) get_central_point() Point {
	return Point{r.x + r.w/2, r.y + r.h/2}
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
func (r *Room) get_wall_points() []Point {
	pts := []Point{}
	for _, p := range r.points {
		if r.is_wall(p) {
			pts = append(pts, p)
		}
	}
	return pts
}
func (r *Room) is_my_point(p Point) bool {
	return p.is_in_slice(r.points)
}
func (r *Room) is_my_inner_point(p Point) bool {
	return p.is_in_slice(r.inner_points)
}
func (r *Room) exists_at(p Point) bool {
	return r.is_my_point(p)
}
func (r *Room) is_wall(p Point) bool {
	return r.is_my_point(p) && !r.is_my_inner_point(p)
}
func (r *Room) get_tile(p Point) tiletype {
	if r.is_wall(p) {
		return WallTile
	}
	return FloorTile
}

/**
 * Corridor
 */
type Corridor struct {
	room_a Room
	room_b Room
	points []Point
}

func new_corridor(l *Level, i int, j int) Corridor {
	c := Corridor{l.rooms[i], l.rooms[j], nil}
	c.init()
	return c
}
func (c *Corridor) init() {
	c.points = c.get_points()
}
func (c *Corridor) get_tile(p Point) tiletype {
	return FloorTile
}
func (c *Corridor) exists_at(p Point) bool {
	return p.is_in_slice(c.points)
}
func (c *Corridor) get_points() []Point {
	a := c.room_a.get_central_point()
	b := c.room_b.get_central_point()
	m := a.get_point_between(b)
	lines := []Line{
		Line{a, Point{m.x, a.y}},
		Line{Point{m.x, a.y}, m},
		Line{m, Point{b.x, m.y}},
		Line{Point{b.x, m.y}, b}}
	pts := []Point{}
	for _, l := range lines {
		pts = append(pts, l.get_points()...)
	}
	return pts
}

/**
 * Door
 */

type Door struct {
	pos    Point
	locked bool
	hidden bool
}

func new_door(pos Point) Door {
	d := Door{pos, false, false}
	return d
}
func (d *Door) exists_at(p Point) bool {
	return d.pos.equals(p)
}
func (d *Door) get_tile(p Point) tiletype {
	return DoorTile
}

/**
 * Staircase
 */

type Staircase struct {
	pos Point
	up  bool
}

func new_staircase(pos Point) Staircase {
	s := Staircase{pos, false}
	return s
}
func (s *Staircase) exists_at(p Point) bool {
	return s.pos.equals(p)
}
func (s *Staircase) get_tile(p Point) tiletype {
	if s.up {
		return StaircaseUpTile
	}
	return StaircaseDownTile
}

/**
 * Line
 */
type Line struct {
	a, b Point
}

func (l *Line) get_points() []Point {
	cp := l.a
	pts := []Point{l.a}
	for !cp.equals(l.b) {
		adjacent := cp.get_surrounding_points()
		slice.Sort(adjacent[:], func(i, j int) bool {
			return adjacent[i].get_distance_to(l.b) < adjacent[j].get_distance_to(l.b)
		})
		pts = append(pts, adjacent[0])
		cp = adjacent[0]
	}
	return pts
}

/**
 * point
 */
type Point struct {
	x, y int
}

func (p *Point) get_surrounding_points() []Point {
	return []Point{
		Point{p.x + 1, p.y},
		Point{p.x, p.y + 1},
		Point{p.x - 1, p.y},
		Point{p.x, p.y - 1}}
}
func (p *Point) equals(p1 Point) bool {
	return p.x == p1.x && p.y == p1.y
}
func (p *Point) get_distance_to(p1 Point) float64 {
	dx := float64(p1.x - p.x)
	dy := float64(p1.y - p.y)
	return math.Sqrt(dx*dx + dy*dy)
}
func (p *Point) get_point_between(p1 Point) Point {
	x := int((p1.x - p.x) / 2)
	y := int((p1.y - p.y) / 2)
	return Point{p.x + x, p.y + y}
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

type Renderer struct {
	tiles map[tiletype]string
}

func new_renderer() *Renderer {
	tiles := map[tiletype]string{
		VoidTile:          " ",
		FloorTile:         ".",
		WallTile:          "#",
		DoorTile:          "█",
		PlayerTile:        "@",
		StaircaseDownTile: "▼",
		StaircaseUpTile:   "▲"}
	return &Renderer{tiles}
}
func (r *Renderer) clear() {
	fmt.Println("\033[H\033[2J")
}
func (r *Renderer) get_texture(tt tiletype) string {
	tile_char, exists := r.tiles[tt]
	if exists {
		return tile_char
	}
	return r.tiles[VoidTile]
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
	fmt.Println()
	fmt.Println(g.status())
	// TODO remove, implement level transition
	if down {
		fmt.Println("down")
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
