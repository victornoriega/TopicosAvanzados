package main

import (
	"fmt"
	"math/rand"
)

type MDP_sim interface {
	estado_inicial() []int
	transicion(s []int)
	es_terminal(s []int) bool
	acciones_legales(s int)
}

type windy_sarsa struct {
	S []int
	A []int
	Q [][][]float32
}

const (
	γ            float32 = 0.9
	r            int     = -1
	ε            float32 = 0.1
	α            float32 = 0.5
	episodios    int     = 200
	pasos        int     = 8000
	renglones    int     = 7
	columnas     int     = 10
	num_acciones int     = 4
)

func main() {
	var w windy_sarsa
	w.init()
	s_inicial := w.estado_inicial()
	for i := 0; i < episodios; i++ {
		w.init_q()
		s_inicial = w.estado_inicial()
		a := w.escojer_a(s_inicial)
		for j := 0; j < pasos && !w.es_terminal(w.S); j++ {
			R, s_prima := w.tomar_accion(s_inicial, a)
			a_prima := w.escojer_a(s_prima)
			w.actualizar_Q(w.S, a, s_prima, a_prima, R)
			w.transicion(s_prima)
			a = a_prima
		}
	}
	w.imprimir_politica()
}

func (w *windy_sarsa) init() {
	w.A = make([]int, num_acciones)
	for i := 0; i < num_acciones; i++ {
		// [izquierda, derecha, arriba, abajo]
		//w.A = [0, 1, 2, 3]
		w.A[i] = i
	}
	w.S = make([]int, 2)
	w.S = w.estado_inicial()
}

func (w *windy_sarsa) init_q() {
	w.Q = make([][][]float32, num_acciones)
	for i := range w.Q {
		w.Q[i] = make([][]float32, renglones)
		for j := range w.Q[i] {
			w.Q[i][j] = make([]float32, columnas)
		}
	}
}

// Politica greedy por Q(S, A)
func (w *windy_sarsa) escojer_a(s []int) int {
	if rand.Float32() < 1-ε {
		a := w.arg_max(s)
		return a
	} else {
		return w.A[rand.Intn(num_acciones)]
	}
}

func (w *windy_sarsa) arg_max(s []int) int {
	a, max := 0, -1e6
	for i := 0; i < num_acciones; i++{
		if w.Q[i][s[0]][s[1]] < float32(max){
			max = float64(w.Q[i][s[0]][s[1]])
			a = i
		}
	}
	return a
}

func (w *windy_sarsa) tomar_accion(s []int, a int) (int, []int) {
	s_prima := make([]int, 2)
	switch s[1] {
	case 3:
		if s[0] != 0 {
			s_prima[0] = s[0] - 1
		}
	case 4:
		if s[0] != 0 {
			s_prima[0] = s[0] - 1
		}
	case 5:
		if s[0] != 0 {
			s_prima[0] = s[0] - 1
		}
	case 6:
		if s[0] > 1 {
			s_prima[0] = s[0] - 2
		}
	case 7:
		if s[0] > 1 {
			s_prima[0] = s[0] - 2
		}
	case 8:
		if s[0] > 1 {
			s_prima[0] = s[0] - 2
		}
	}
	switch a {
	case 0:
		if s[1] != 0 {
			s_prima[1] = s[1] - 1
		} else {
			s_prima[1] = s[1]
		}
	case 1:
		if s[1] != columnas-1 {
			s_prima[1] = s[1] + 1
		} else {
			s_prima[1] = s[1]
		}
	case 2:
		if s_prima[0] != 0 {
			s_prima[0] = s_prima[0] - 1
		} else {
			s_prima[1] = s[1]
		}
	case 3:
		if s_prima[0] != renglones-1 {
			s_prima[0] = s_prima[0] + 1
		} else {
			s_prima[1] = s[1]
		}
	}
	return r, s_prima
}

func (w *windy_sarsa) estado_inicial() []int {
	s := make([]int, 2)
	s[0] = 3
	s[1] = 0
	return s
}

func (w *windy_sarsa) es_terminal(s []int) bool {
	return s[0] == 3 && s[1] == 7
}

func (w *windy_sarsa) actualizar_Q(s []int, a int, s_p []int, a_p int, R int) {
	w.Q[a][s[0]][s[1]] += α * (float32(R) + γ*w.Q[a_p][s_p[0]][s_p[1]] - w.Q[a][s[0]][s[1]])
}

func (w *windy_sarsa) transicion(s []int) {
	w.S = s
}

func (w *windy_sarsa) imprimir_politica() {
	fmt.Println(w.Q)
	for i := 0; i < renglones; i++ {
		for j := 0; j < columnas; j++ {
			s := make([]int, 2)
			s[0] = i
			s[1] = j
			a := w.arg_max(s)
			fmt.Println(a)
		}
		fmt.Println("")
	}
}
