package main

import (
	"math/rand"
	"time"
)

type MovimentoFantasma struct {
	MX, MY int
}

// Goroutine do fantasma (move aleatoriamente a cada 450ms), não altera o estado do jogom envia "sugestões" para a processadora de eventos
func iniciarFantasma(saidas chan<- MovimentoFantasma, encerrar <-chan struct{}) {
	tick := time.NewTicker(450 * time.Millisecond)
	defer tick.Stop()
	//RNG
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-encerrar: // canal para encerrar a goroutine graciosamente
			return
		case <-tick.C: // a cada 450ms o canal recebe um "tick"
			mx, my := 0, 0
			switch r.Intn(4) {
			case 0:
				mx = +1 //direita
			case 1:
				mx = -1 //esquerda
			case 2:
				my = +1 //baixo
			case 3:
				my = -1 //cima
			}

			saidas <- MovimentoFantasma{MX: mx, MY: my}
		}
	}
}
