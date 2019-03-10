package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	capital     int     = 99
	meta        int     = 100
	num_estados int     = 101 // los dos "dummy" states son el 0 y el 100
	ph1         float32 = 0.4
	ph2         float32 = 0.25
	ε           float64 = 0.01
)

type Gambler_MDP struct {
	S [num_estados]estado
	A [capital]int
	R [num_estados]int
	γ float32
}
type estado struct {
	v float32
	π int
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var G Gambler_MDP
	G.init()
	G.iteracion_valores()
	G.imprimir_politica()
}

func (G *Gambler_MDP) iteracion_valores() {
	for Δ := 1.0; Δ > ε; {
		Δ = 0.0
		for s := meta - 1; s > 0; s-- {
			v := G.S[s].v
			nuevo_valor, a := G.arg_max(s)
			G.S[s].v = nuevo_valor
			G.S[s].π = a
			Δ = math.Max(float64(Δ), float64(math.Abs(float64(v-nuevo_valor))))
		}
	}
}

func (G *Gambler_MDP) arg_max(s int) (float32, int) {
	num_acciones := int(math.Min(float64(s), float64(meta-s)))
	max := ph1*G.S[s].v + (1.0-ph1)*G.S[s].v
	accion := 1
	for a := 1; a <= num_acciones; a++ {
		//nuevo estado s': ganas
		nuevo_valor := ph1*(float32(G.R[s+a])+G.S[s+a].v) + (1.0-ph1)*G.S[s-a].v
		if nuevo_valor > max {
			max = nuevo_valor
			accion = a
		}
	}
	return max, accion
}

func (G *Gambler_MDP) init() {
	G.R[meta] = 1
	G.γ = 1
}

func (G *Gambler_MDP) obtener_politicas() {
	for i := 0; i < meta; i++ {
		_, G.S[i].π = G.arg_max(i)
	}
}

func (G *Gambler_MDP) imprimir_politica() {
	for i := 0; i < meta; i++ {
		fmt.Println(G.S[i].π)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "El problema del apostador (aparte de ser apostador)"
	p.X.Label.Text = "La feria que tiene"
	p.Y.Label.Text = "La feria que deberia apostar"
	politica := make([]int, meta)
	for i := 1; i < meta; i++ {
		politica[i] = G.S[i].π
	}
	pts := Points(politica)
	err = plotutil.AddLinePoints(p,
		"0.4", pts)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}

}

func Points(p []int) plotter.XYs {
	pts := make(plotter.XYs, meta)
	for i := 1; i < meta; i++ {
		pts[i].Y = float64(p[i])
		pts[i].X = float64(i)
	}
	return pts
}
