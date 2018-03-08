package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Color struct {
	R int
	G int
	B int
}

func (c Color) ToHex() string {
	s := fmt.Sprintf("%02x%02x%02x", c.R, c.G, c.B)
	return s
}

func (c *Color) SetColor(r int, g int, b int) {
	c.R = r % 256
	c.G = g % 256
	c.B = b % 256
}

// mix two colors
func blendColors(a Color, b Color) Color {
	newColor := Color{0, 0, 0}
	newColor.R = mixColorElement(a.R, b.R)
	newColor.G = mixColorElement(a.G, b.G)
	newColor.B = mixColorElement(a.B, b.B)

	return newColor
}

// take e.g. two r-values and take random parts of them and combine
func mixColorElement(motherElement int, fatherElement int) int {
	pos := [8]int{0, 1, 2, 3, 4, 5, 6, 7}

	// Mix order
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(pos); n > 0; n-- {
		randIndex := r.Intn(n)
		pos[n-1], pos[randIndex] = pos[randIndex], pos[n-1]
	}

	motherPos := pos[0:4]
	fatherPos := pos[4:8]

	newElement := 0
	for _, bit := range motherPos {
		gene := int(math.Exp2(float64(bit)))
		if motherElement&gene != 0 {
			newElement += gene
			if r.Intn(100) >= 90 {
				newElement -= gene
			}
		} else {
			if r.Intn(100) >= 90 {
				newElement += gene
			}
		}
	}

	for _, bit := range fatherPos {
		gene := int(math.Exp2(float64(bit)))
		if fatherElement&gene != 0 {
			newElement += gene
			if r.Intn(100) >= 90 {
				newElement -= gene
			}
		} else {
			if r.Intn(100) >= 90 {
				newElement += gene
			}
		}
	}

	return newElement
}

// The difference between two colors
func colorDifference(a Color, b Color) int {
	rC := float64(a.R - b.R)
	gC := float64(a.G - b.G)
	bC := float64(a.B - b.B)

	return int(math.Sqrt(rC*rC + gC*gC + bC*bC))
}
