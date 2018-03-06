package main

import (
	"fmt"
	"html"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type color struct {
	r int
	g int
	b int
}

func blendColors(a color, b color) color {
	newColor := color{0, 0, 0}
	newColor.r = (a.r + b.r) / 2
	newColor.g = (a.g + b.g) / 2
	newColor.b = (a.b + b.b) / 2

	return newColor
}

// The difference between two colors
func colorDifference(a color, b color) int {
	rC := math.Abs(float64(a.r - b.r))
	gC := math.Abs(float64(a.g - b.g))
	bC := math.Abs(float64(a.b - b.b))

	return int(rC + gC + bC)
}

type creature struct {
	age   int
	color color
}

type world struct {
	color color
}

type population []creature

func (p population) avgAge() int {
	cumulativeAge := 0
	for _, creature := range p {
		cumulativeAge += creature.age
	}
	return cumulativeAge / len(p)
}

func (p population) avgColor() color {
	cumulativeColor := color{0, 0, 0}
	for _, creature := range p {
		cumulativeColor.r = (cumulativeColor.r + creature.color.r) / 2
		cumulativeColor.g = (cumulativeColor.g + creature.color.g) / 2
		cumulativeColor.b = (cumulativeColor.b + creature.color.b) / 2
	}
	return cumulativeColor
}

//func (p people) Len() int           { return len(p) }
//func (p people) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
//func (p people) Less(i, j int) bool { return p[i].age < p[j].age }

func seedPopulation(num int) population {
	pop := population{}
	for i := 0; i < num; i++ {
		creature := creature{age: rand.Intn(100), color: color{rand.Intn(256), rand.Intn(256), rand.Intn(256)}}
		pop = append(pop, creature)

	}
	return pop
}

// Every tick, creatures between 0 and 15 years have a 20% chance of dying
// creatures between 15 and 50 years have a 10% chance of dying
// creatures between 50 and 80 have a 40% chance of dying
// creatures between 80 and 100 have a 80% chance of dying
func runSimulation(world world, pop population, wg *sync.WaitGroup, stop <-chan os.Signal) {
	ticker := time.NewTicker(time.Millisecond * 100)
	run := true
	for {
		select {
		case <-ticker.C:
			for i, cr := range pop {
				creatureDies := false
				age := cr.age
				color := cr.color
				age++

				// 25% chance that creature is killed by predators if it stands out
				if colorDifference(world.color, color) > 100 {
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
					pop[i] = creature{age: 0, color: blendColors(mother.color, father.color)}
				} else {
					// Save changes
					pop[i] = creature{age: age, color: color}
				}

			}
			fmt.Printf("Average age  : %d\n", pop.avgAge())
			fmt.Printf("Average color: %v\n", pop.avgColor())
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

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	initPop := seedPopulation(1e4)
	world := world{color: color{100, 50, 80}}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	go http.ListenAndServe(":8080", nil)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	runSimulation(world, initPop, wg, stop)

	wg.Wait()
}
