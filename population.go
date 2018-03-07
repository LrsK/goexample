package main

import "math/rand"

type Creature struct {
	Age   int
	Color Color
}

type World struct {
	Color Color
}

func (w *World) SetColor(r int, g int, b int) {
	w.Color.SetColor(r, g, b)
}

type Population []Creature

func (p Population) AvgAge() int {
	cumulativeAge := 0
	for _, creature := range p {
		cumulativeAge += creature.Age
	}
	return cumulativeAge / len(p)
}

func (p Population) AvgColor() Color {
	cumulativeColor := Color{0, 0, 0}
	for _, creature := range p {
		cumulativeColor.R = (cumulativeColor.R + creature.Color.R) / 2
		cumulativeColor.G = (cumulativeColor.G + creature.Color.G) / 2
		cumulativeColor.B = (cumulativeColor.B + creature.Color.B) / 2
	}
	return cumulativeColor
}

func newPopulation(num int) Population {
	pop := Population{}
	for i := 0; i < num; i++ {
		creature := Creature{Age: rand.Intn(100), Color: Color{rand.Intn(256), rand.Intn(256), rand.Intn(256)}}
		pop = append(pop, creature)

	}
	return pop
}
