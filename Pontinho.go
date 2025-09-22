package main

import (
	"time"
)

// Canal global para reaparecer pontinhos
var reaparecerPontinho = make(chan struct{ x, y int })

// Inicializa a goroutine que escuta os pontinhos e faz eles reaparecerem
func IniciarPontinhos(jogo *Jogo) {
	go func() {
		for {
			select {
			case ponto := <-reaparecerPontinho:
				// Goroutine que aguarda 6s para reaparecer o pontinho
				go func(px, py int) {
					select {
					case <-time.After(6 * time.Second):
						jogo.Mapa[py][px] = Pontinho
						interfaceDesenharJogo(jogo)
					}
				}(ponto.x, ponto.y)

			case <-time.After(1 * time.Minute):
				// Timeout de segurança para evitar que a goroutine fique presa indefinidamente
			}
		}
	}()
}

// Coleta o pontinho: remove do mapa e envia para o canal reaparecerPontinho
func ColetarPontinho(jogo *Jogo, x, y int) {
	jogo.Mapa[y][x] = Vazio
	reaparecerPontinho <- struct{ x, y int }{x, y}
}

// Verifica se o elemento é pontinho
func EhPontinho(elem Elemento) bool {
	return elem.simbolo == Pontinho.simbolo
}
