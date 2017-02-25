package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
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
	g.level.render()
}

func (g *Game) over() bool {
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	switch scan.Text() {
	case "c":
		return true
	case "j":
		g.pc.y += 1
	case "h":
		g.pc.x -= 1
	case "l":
		g.pc.x += 1
	case "k":
		g.pc.y -= 1
	}
	return false
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
func (l *Level) put_player(p *PlayerCharacter) {
	x := l.rooms[0].x + 1
	y := l.rooms[0].y + 1
	p.x, p.y = x, y
	l.pc = p
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

/**
 * Player Character
 */

type PlayerCharacter struct {
	x, y int
}

/**
 * main
 */

func main() {
	game := new_game()
	for {
		game.update()
		if game.over() {
			break
		}
	}
}
