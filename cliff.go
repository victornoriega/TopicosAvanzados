package main

import (
	"fmt"
	"math/rand"
)

const (
	altura       = 4
	largo        = 12
	num_acciones = 4
	num_estados  = altura * largo
	r_cliff      = -100
	r            = -1
	s_inicial    = 36
	s_final      = 47
	ε            = 0.1
	α            = 0.9
	episodios    = 2000
	pasos        = 800
	γ            = 0.9
)

type Accion int

type MDP_sim interface {
	estado_inicial() Accion
	transicion(s int, a Accion)
	es_terminal(s int) bool
	acciones_legales(s int)
}

type Qlearning struct {
	acciones [num_acciones]Accion
	Q        [num_estados][num_acciones]float32
}

func main() {
	var cliff Qlearning
	for i := 0; i < episodios; i++ {
		s := cliff.estado_inicial()
		a := cliff.escojer_a(s)
		for j := 0; j < pasos && !cliff.es_terminal(s); j++ {
			R, s_prima := cliff.tomar_accion(s, a)
			a_prima := cliff.escojer_a(s_prima)
			cliff.actualizar_Q(s, a, s_prima, a_prima, R)
			s = s_prima
		}
	}
	cliff.imprimir_politica()
}

func (c *Qlearning) estado_inicial() int {
	return s_inicial
}
func (c *Qlearning) escojer_a(s int) Accion {
	if rand.Float32() < 1-ε {
		a := c.arg_max(s)
		return a
	} else {
		return c.acciones[rand.Intn(num_acciones)]
	}
}

func (c *Qlearning) arg_max(s int) Accion {
	a, max := 0, -1e6
	for i := 0; i < num_acciones; i++ {
		if c.Q[s][i] > float32(max) {
			max = float64(c.Q[s][i])
			a = i
		}
	}
	return Accion(a)
}

func (c *Qlearning) es_terminal(s int) bool {
	return s == s_final
}

func (c *Qlearning) tomar_accion(s int, a Accion) (int, int) {
	switch a {
	case 0:
		// Subir
		if s >= 0 && s < 12 {
			return -1, s
		} else {
			return -1, s - 12
		}
	case 1:
		//bajar
		if s == s_inicial || s == s_final {
			return -1, s
		} else if s > s_inicial && s < s_final {
			return -100, s
		} else {
			return -1, s + 12
		}
	case 2:
		//izquierda
		for i := 0; i < 4; i++ {
			if s == i*12 {
				return -1, s
			}
		}
		if s > s_inicial && s <= s_final {
			return -100, s - 1
		} else {
			return -1, s - 1
		}
	case 3:
		//derecha
		for i := 11; i <= 35; i += 12 {
			if s == i {
				return -1, s
			}
		}
		if s >= s_inicial && s < s_final-1 {
			return -100, s + 1
		} else {
			return -1, s + 1
		}
	}
	return -1, s
}

func (c *Qlearning) actualizar_Q(s int, a Accion, s_p int, a_p Accion, R int) {
	c.Q[s][a] += α * (float32(R) + γ*float32(c.arg_max(s)) - float32(c.Q[s][a]))
}

func (c *Qlearning) imprimir_politica() {
	for i := 0; i < s_final; i++ {
		fmt.Println(i, "\t", c.arg_max(i))
	}
}
