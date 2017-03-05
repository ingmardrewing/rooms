package main

import "fmt"

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
