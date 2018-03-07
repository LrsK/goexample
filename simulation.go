package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Every tick, creatures between 0 and 15 years have a 20% chance of dying
// creatures between 15 and 50 years have a 10% chance of dying
// creatures between 50 and 80 have a 40% chance of dying
// creatures between 80 and 100 have a 80% chance of dying
func runSimulation(world *World, pop Population, wg *sync.WaitGroup, stop <-chan os.Signal) {
	ticker := time.NewTicker(time.Millisecond * 100)
	run := true
	for {
		select {
		case <-ticker.C:
			for i, cr := range pop {
				creatureDies := false
				age := cr.Age
				color := cr.Color
				age++

				// 25% chance that creature is killed by predators if it stands out
				if colorDifference(world.Color, color) > 200 {
					if rand.Intn(4) == 1 {
						creatureDies = true
					}
				}

				if age >= 0 && age < 15 {
					// 10% chance this creature dies
					if rand.Intn(10) == 1 {
						creatureDies = true
					}
				} else if age >= 16 && age < 50 {
					// 5% chance this creature dies
					if rand.Intn(20) == 1 {
						creatureDies = true
					}
				} else if age >= 50 && age < 80 {
					// 20% chance this creature dies
					if rand.Intn(5) == 1 {
						creatureDies = true
					}
				} else if age >= 80 {
					// 80% chance this creature dies
					if rand.Intn(100) >= 20 {
						creatureDies = true
					}
				}

				if creatureDies {
					// Replace with a new creature
					mother := pop[rand.Intn(len(pop))]
					father := pop[rand.Intn(len(pop))]
					pop[i] = Creature{Age: 0, Color: blendColors(mother.Color, father.Color)}
				} else {
					// Save changes
					pop[i] = Creature{Age: age, Color: color}
				}

			}
			fmt.Printf("Average age  : %d\n", pop.AvgAge())
			fmt.Printf("Average color: %v\n", pop.AvgColor())
		case <-stop:
			fmt.Println("I got interrupted so I'm stopping")
			ticker.Stop()
			run = false
		}
		if !run {
			break
		}
	}
	fmt.Println("We're done!")
	wg.Done()
}
