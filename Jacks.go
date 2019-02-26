package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	max_carros       int     = 20
	max_transf       int     = 5
	recompensa_renta int     = 10
	costo_mover      int     = -2
	θ                float64 = 1e-6
	λ_solicitud_A    int     = 3
	λ_solicitud_B    int     = 4
	λ_regreso_A      int     = 3
	λ_regreso_B      int     = 2
)

type MDP_Jack struct {
	S [21][21][2]float64 // Estados (dos arreglos: uno que tiene V, y otro la politica)
	γ float64            // Descuento
	P [21][21][2]float64 //Probabilidades
	R [21]float64        //Recompensas
	A [11]int            //Acciones
}

func factorial(n float64) float64 {
	if n > 0 {
		return n * factorial(n-1.0)
	}
	return 1
}

func poisson(λ float64, n int64) float64 {
	return (math.Pow(λ, float64(n)) / (factorial(float64(n)) *
		math.Exp(float64(-λ))))
}

func main() {
	rand.Seed(time.Now().Unix())
	var j MDP_Jack
	j.init()

	j.iteracion_politicas()
	j.imprimir_politica()
}
func (M *MDP_Jack) imprimir_politica() {
	for i := 0; i <= max_carros; i++ {
		for j := 0; j <= max_carros; j++ {
			fmt.Println("%3v", M.S[i][j][1])
		}
		fmt.Println("")
	}
}
func (M *MDP_Jack) iteracion_politicas() {
	hubo_cambios := true
	for hubo_cambios {
		M.evaluacion_pol()
		hubo_cambios = M.mejora_pol()
	}
}

func (M *MDP_Jack) evaluacion_pol() {
	Δ := 1.0
	for Δ > θ {
		Δ = 0
		for n := 0; n <= max_carros; n++ {
			for m := 0; m <= max_carros; m++ {
				v := M.S[n][m][0]
				a := int(M.S[n][m][1])
				M.S[n][m][0] = M.obtenerValor(n, m, a)
				Δ = math.Max(Δ, math.Abs(float64(M.S[n][m][0])-float64(v)))
			}
		}
	}
}

func (M *MDP_Jack) obtenerValor(n int, m int, a int) float64 {
	val := float64(costo_mover) * math.Abs(float64(a))

	// act_1 es la variable que determina cuantos carros habra actualizado despues
	// de pasar "a" carros a "m" (a puede ser negativa)
	// act_2 es lo analogo
	act_n := n - a
	act_m := m + a

	for _n := 0; _n <= max_carros; _n++ {
		for _m := 0; _m <= max_carros; _m++ {
			val += M.P[act_n][_n][0]*M.P[act_m][_m][1]*M.R[act_n] +
				(M.R[act_m] + M.γ*M.S[act_n][_m][0])
		}
	}
	return val
}

func (M *MDP_Jack) mejora_pol() bool {
	politica_estable := true
	for i := 0; i <= max_carros; i++ {
		for j := 0; j <= max_carros; j++ {
			var vieja_accion = M.S[i][j][1]
			π := M.arg_max(i, j)
			M.S[i][j][1] = π
			if π != vieja_accion {
				politica_estable = false
			}
		}
	}
	return politica_estable
}

func (M *MDP_Jack) arg_max(n int, m int) float64 {
	v := M.S[n][m][0]
	mejor_accion := M.S[n][m][1]
	for _, a := range M.A {
		aux := M.obtenerValor(n, m, a)
		if aux > v {
			v = aux
			mejor_accion = float64(a)
		}
	}
	return mejor_accion
}

func (M *MDP_Jack) init() {
	for i := 0; i <= max_carros; i++ {
		for j := 0; j <= max_carros; j++ {
			// El 0 es para el valor del estado
			M.S[i][j][0] = 0
			// El 1 es para la accion que corresponde al estado segun la politica
			M.S[i][j][1] = 0
		}
	}
	for i := 0; i <= 10; i++ {
		M.A[i] = 0
	}
}
