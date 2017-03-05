package main

import (
	"math"

	"github.com/bradfitz/slice"
)

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
 * helper functions
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
