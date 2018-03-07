package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	initPop := newPopulation(1e4)
	world := &World{Color: Color{0, 0, 0}}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cr, _ := strconv.Atoi(r.FormValue("r"))
		cg, _ := strconv.Atoi(r.FormValue("g"))
		cb, _ := strconv.Atoi(r.FormValue("b"))

		world.SetColor(cr, cg, cb)

		fmt.Printf("%d\n", colorDifference(world.Color, initPop.AvgColor()))

		tmpl := template.Must(template.ParseFiles("main.html"))
		tmpl.Execute(w, world)
	})

	http.HandleFunc("/pop", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pop.html"))
		tmpl.Execute(w, initPop)
	})

	// Run web server
	go http.ListenAndServe(":8080", nil)

	// Run population sim
	go runSimulation(world, initPop, wg, stop)

	wg.Wait()
}
