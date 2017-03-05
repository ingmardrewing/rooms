package main

type LevelElement interface {
	exists_at(p Point) bool
	get_tile(p Point) tiletype
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
			g.clear_level()
			g.next_level()
		}
	}

	new_pc_pos := Point{x, y}
	wkbl := g.level.get_walkable_points()
	if new_pc_pos.is_in_slice(wkbl) {
		g.pc.pos = new_pc_pos
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
func (l *Level) get_random_room() *Room {
	i := get_rand(len(l.rooms))
	r := l.rooms[i]
	return r
}
func (l *Level) put_player(pc *PlayerCharacter) {
	p := l.get_random_room().get_random_inner_point()
	pc.pos = p
	l.pc = pc
	l.elements = append(l.elements, pc)
	l.reverse_elements()
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
 * Player Character
 */

type PlayerCharacter struct {
	pos Point
}

func (pc *PlayerCharacter) exists_at(p Point) bool {
	return pc.pos.equals(p)
}
func (pc *PlayerCharacter) get_tile(p Point) tiletype {
	return PlayerTile
}
