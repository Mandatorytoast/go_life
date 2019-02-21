package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	FG   = termbox.ColorWhite
	BG   = termbox.ColorDefault
	RUNE = 'O'
)

type Colony struct {
	XPos        int
	YPos        int
	alive       bool
	surrounding []string
}

func (c *Colony) createKey() string {
	return strconv.Itoa(c.XPos) + "|" + strconv.Itoa(c.YPos)
}

func (c *Colony) generateSurroundingKeys() []string {
	var keys []string
	for i := -1; i <= 1; i++ {
		for n := -1; n <= 1; n++ {
			if i == 0 && n == 0 {
				continue
			} else {
				keys = append(keys, strconv.Itoa(c.XPos+i)+"|"+strconv.Itoa(c.YPos+n))
			}
		}
	}
	return keys
}

func (c *Colony) checkSurrounding(colonies map[string]Colony) int {
	var count int
	for _, i := range c.surrounding {
		if colonies[i].alive {
			count += 1
		}
	}
	return count
}

func main() {
	iterations, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	initerr := termbox.Init()
	defer termbox.Close()
	rand.Seed(time.Now().UnixNano())
	event_queue := make(chan termbox.Event)
	draw_tick := time.NewTicker(60 * time.Millisecond)
	if initerr != nil {
		fmt.Println("Something fucked up...")
	}

	width, height := termbox.Size()
	colonies := make(map[string]Colony)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r := rand.Intn(20)
			var a bool
			if r == 1 {
				a = true
			} else {
				a = false
			}
			tmp := Colony{x, y, a, nil}
			tmp.surrounding = tmp.generateSurroundingKeys()
			colonies[tmp.createKey()] = tmp
		}
	}
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

loop:
	for i := 0; i < iterations; i++ {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		case <-draw_tick.C:
			for _, value := range colonies {
				if value.alive {
					termbox.SetCell(value.XPos, value.YPos, RUNE, FG, BG)
				} else {
					termbox.SetCell(value.XPos, value.YPos, ' ', BG, BG)
				}
			}
			for key, value := range colonies {
				tmp := colonies[key]
				if value.alive {
					if value.checkSurrounding(colonies) <= 1 || value.checkSurrounding(colonies) >= 4 {
						tmp.alive = false
					}
				} else {
					if value.checkSurrounding(colonies) == 3 {
						tmp.alive = true
					}
				}
				colonies[key] = tmp
			}
			termbox.Flush()
		}
	}
}
