package main

import "math"

const (
	Top    = iota
	Left   = iota
	Bottom = iota
	Right  = iota
)

type LevelElement interface {
	exists_at(p Point) bool
	get_tile(p Point) tiletype
	get_gamepoint(p Point) *GamePoint
}

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

func (g *Game) clear_level() {
	g.level = nil
}

func (g *Game) generate_level() {
	g.level = &Level{60, 32, nil, nil, nil, nil, nil, nil}
	g.level.init()
}

func (g *Game) init_player() {
	g.pc = new_playercharacter()
	g.level.put_player(g.pc)
}

func (g *Game) handle_user_input(c string) {
	x := g.pc.gp.pos.x
	y := g.pc.gp.pos.y
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
		if g.level.is_staircase(g.pc.gp.pos) {
			g.clear_level()
			g.next_level()
		}
	}

	new_pc_pos := Point{x, y}
	wkbl := g.level.get_walkable_points()
	if new_pc_pos.is_in_slice(wkbl) {
		g.pc.move_to(new_pc_pos)
	}
}

func (g *Game) next_level() {
	g.clear_level()
	g.generate_level()
	g.level.put_player(g.pc)
}

func new_game() Game {
	g := Game{nil, nil}
	g.generate_level()
	g.init_player()
	return g
}

/**
 * Level
 */

type Level struct {
	width, height int
	elements      []LevelElement
	rooms         []*Room
	corridors     []*Corridor
	doors         []*Door
	staircases    []*Staircase
	pc            *PlayerCharacter
}

func (l *Level) init() {
	l.elements = []LevelElement{}
	l.rooms = l.generate_rooms()
	l.corridors = l.generate_corridors()
	l.doors = l.generate_doors()
	l.staircases = l.generate_staircases()
}

func (l *Level) put_player(pc *PlayerCharacter) {
	r := l.get_random_room()
	r.mark_gamepoints_as_seen()
	p := r.get_random_inner_point()
	pc.move_to(p)
	l.pc = pc
	l.elements = append(l.elements, pc)
	l.reverse_elements()
}

func (l *Level) reverse_elements() {
	reversed := []LevelElement{}
	for i := len(l.elements) - 1; i >= 0; i-- {
		reversed = append(reversed, l.elements[i])
	}
	l.elements = reversed
}

func (l *Level) add_element(e LevelElement) {
	l.elements = append(l.elements, e)
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

func (l *Level) generate_rooms() []*Room {
	rooms := []*Room{}
	row_height := l.height / 2
	col_width := l.width / 3
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			a := Point{j*col_width + j, i*row_height + i}
			b := Point{(j + 1) * col_width, (i + 1) * row_height}
			r := new_room(a, b)
			l.add_element(&r)
			rooms = append(rooms, &r)
		}
	}
	return rooms
}

func (l *Level) generate_corridors() []*Corridor {
	from := []int{0, 1, 3, 4, 0, 1, 2}
	to := []int{1, 2, 4, 5, 3, 4, 5}
	crs := []*Corridor{}
	for i, _ := range from {
		c := new_corridor(l, from[i], to[i])
		l.add_element(c)
		crs = append(crs, c)
	}
	return crs
}

func (l *Level) generate_doors() []*Door {
	wall_pts := l.get_all_wall_points()
	corr_pts := l.get_all_corridor_points()
	door_pts := l.find_point_set_intersections(wall_pts, corr_pts)
	return l.generate_doors_at(door_pts)
}

func (l *Level) generate_doors_at(pts []Point) []*Door {
	doors := []*Door{}
	for _, p := range pts {
		d := new_door(p)
		l.add_element(d)
		doors = append(doors, d)
	}
	return doors
}

func (l *Level) generate_staircases() []*Staircase {
	s := new_staircase(l.get_random_room().get_random_inner_point())
	l.add_element(&s)
	return []*Staircase{&s}
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
	for _, e := range l.elements {
		if e.exists_at(p) {
			return e.get_tile(p)
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

func (l *Level) get_gamepoints() []*GamePoint {
	gps := []*GamePoint{}
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			gp := l.get_gamepoint(Point{x, y})
			gps = append(gps, gp)
		}
	}
	return gps
}

func (l *Level) get_gamepoint(p Point) *GamePoint {
	for _, e := range l.elements {
		if e.exists_at(p) {
			gp := e.get_gamepoint(p)
			return gp
		}
	}
	gp := new_gamepoint(
		p,
		VoidTile,
		true,
		false)
	return &gp
}

func (l *Level) get_random_room() *Room {
	i := get_rand(len(l.rooms))
	r := l.rooms[i]
	return r
}

/**
 * Room
 */

type Room struct {
	pos          Point
	w, h         int
	points       []Point
	inner_points []Point
	gamepoints   []*GamePoint
}

func new_room(a, b Point) Room {
	dx := b.x - a.x
	dy := b.y - a.y
	w := get_rand_range(3, dx-1)
	h := get_rand_range(3, dy-1)
	x := a.x + get_rand(dx-w)
	y := a.y + get_rand(dy-h)
	r := Room{Point{x, y}, w, h, nil, nil, nil}
	r.init()
	return r
}

func (r *Room) init() {
	r.points = r.get_points()
	r.inner_points = r.get_inner_points()
	r.init_gamepoints()
}

func (r *Room) get_random_inner_point() Point {
	pts := r.get_inner_points()
	i := get_rand(len(pts))
	return pts[i]
}

func (r *Room) get_central_point() Point {
	return Point{r.pos.x + r.w/2, r.pos.y + r.h/2}
}

func (r *Room) get_inner_points() []Point {
	a := Point{r.pos.x + 1, r.pos.y + 1}
	b := Point{r.pos.x + r.w - 1, r.pos.y + r.h - 1}
	return get_rect_points(a, b)
}

func (r *Room) get_points() []Point {
	b := Point{r.pos.x + r.w, r.pos.y + r.h}
	return get_rect_points(r.pos, b)
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

func (r *Room) get_central_wall_point(side int) Point {
	wall := r.get_wall_points()
	m := r.get_central_point()
	for _, wp := range wall {
		if side == Top && wp.y == r.pos.y && wp.x == m.x {
			return wp
		}
		if side == Bottom && wp.y == r.pos.y+r.h && wp.x == m.x {
			return wp
		}
		if side == Left && wp.x == r.pos.x && wp.y == m.y {
			return wp
		}
		if side == Right && wp.x == r.pos.x+r.w && wp.y == m.y {
			return wp
		}
	}
	return Point{0, 0}
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

func (r *Room) init_gamepoints() {
	pts := r.get_points()
	gps := []*GamePoint{}
	for _, p := range pts {
		gp := new_gamepoint(
			p,
			r.get_tile(p),
			true,
			false)
		gps = append(gps, &gp)
	}
	r.gamepoints = gps

}

func (r *Room) mark_gamepoints_as_seen() {
	for _, gp := range r.gamepoints {
		gp.seen = true
	}
}

func (r *Room) get_gamepoint(p Point) *GamePoint {
	for _, gp := range r.gamepoints {
		if gp.pos.equals(p) {
			return gp
		}
	}
	return nil
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
	room_a *Room
	room_b *Room
	points []Point
}

func new_corridor(l *Level, i int, j int) *Corridor {
	c := Corridor{l.rooms[i], l.rooms[j], nil}
	c.init()
	return &c
}

func (c *Corridor) init() {
	c.points = c.get_points()
}

func (c *Corridor) get_tile(p Point) tiletype {
	return FloorTile
}

func (c *Corridor) get_gamepoint(p Point) *GamePoint {
	gp := new_gamepoint(
		p,
		FloorTile,
		true,
		false)
	return &gp
}

func (c *Corridor) exists_at(p Point) bool {
	return p.is_in_slice(c.points)
}

func (c *Corridor) get_points() []Point {
	ws_a, ws_b := c.get_wall_sides()
	start, corner_1, middle, corner_2, end := c.get_defining_points(ws_a, ws_b, ws_a == Top || ws_a == Bottom)
	lines := c.get_lines(start, corner_1, middle, corner_2, end)
	return c.get_points_from_lines(lines)
}

func (c *Corridor) get_defining_points(ws_a int, ws_b int, vertical bool) (Point, Point, Point, Point, Point) {
	start := c.room_a.get_central_wall_point(ws_a)
	end := c.room_b.get_central_wall_point(ws_b)
	middle := start.get_point_between(end)
	if vertical {
		return start, Point{start.x, middle.y}, middle, Point{end.x, middle.y}, end
	}
	return start, Point{middle.x, start.y}, middle, Point{middle.x, end.y}, end
}

func (c *Corridor) get_points_from_lines(lines []Line) []Point {
	pts := []Point{}
	for _, l := range lines {
		pts = append(pts, l.get_points()...)
	}
	return pts
}

func (c *Corridor) get_lines(wp_a Point, corner_1 Point, m Point, corner_2 Point, wp_b Point) []Line {
	return []Line{
		Line{wp_a, corner_1},
		Line{corner_1, m},
		Line{m, corner_2},
		Line{corner_2, wp_b}}
}

func (c *Corridor) get_wall_sides() (int, int) {
	a := c.room_a.get_central_point()
	b := c.room_b.get_central_point()
	dx := float64(b.x - a.x)
	dy := float64(b.y - a.y)
	if math.Abs(dx) > math.Abs(dy) && dx < 0 {
		return Left, Right
	} else if math.Abs(dx) > math.Abs(dy) {
		return Right, Left
	} else if dy < 0 {
		return Top, Bottom
	}
	return Bottom, Top
}

/**
 * Door
 */

type Door struct {
	pos    Point
	locked bool
	hidden bool
}

func new_door(pos Point) *Door {
	d := Door{pos, false, false}
	return &d
}

func (d *Door) exists_at(p Point) bool {
	return d.pos.equals(p)
}

func (d *Door) get_tile(p Point) tiletype {
	return DoorTile
}

func (d *Door) get_gamepoint(p Point) *GamePoint {
	gp := new_gamepoint(
		p,
		DoorTile,
		true,
		false)
	return &gp
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

func (s *Staircase) get_gamepoint(p Point) *GamePoint {
	gp := new_gamepoint(
		p,
		s.get_tile(p),
		true,
		false)
	return &gp
}

/**
 * Player Character
 */

type PlayerCharacter struct {
	gp *GamePoint
}

func new_playercharacter() *PlayerCharacter {
	gp := new_gamepoint(
		Point{0, 0},
		PlayerTile,
		true,
		false)
	gp.seen = true
	pc := PlayerCharacter{&gp}
	return &pc
}

func (pc *PlayerCharacter) exists_at(p Point) bool {
	return pc.gp.pos.equals(p)
}

func (pc *PlayerCharacter) get_tile(p Point) tiletype {
	return PlayerTile
}

func (pc *PlayerCharacter) set_gamepoint(gp GamePoint) {
	pc.gp = &gp
}

func (pc *PlayerCharacter) move_to(p Point) {
	pc.gp.pos = p
}

func (pc *PlayerCharacter) get_gamepoint(p Point) *GamePoint {
	return pc.gp
}

/**
 * GamePoint
 */

type GamePoint struct {
	tile       tiletype
	pos        Point
	seen       bool
	persistent bool
	moving     bool
}

func (gp *GamePoint) get_tile() tiletype {
	if gp.seen {
		return gp.tile
	}
	return VoidTile
}

func new_gamepoint(pos Point, tile tiletype, persistent bool, moving bool) GamePoint {
	return GamePoint{tile, pos, false, persistent, moving}
}
