package main

import (
	"math/rand"
	"time"

	"github.com/pkg/term" // for getting userinput without the user hitting enter
)

/**
 * main
 */

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
