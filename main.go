// main.go - Loop principal do jogo
package main

import (
	"os" 
	"time");

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Goroutine que escuta os pontinhos para reaparecer
	go func() {
		for {
			ponto := <-reaparecerPontinho
			go func(px, py int) {
				<-time.After(6 * time.Second) // espera 6000ms
				jogo.Mapa[py][px] = Pontinho
				interfaceDesenharJogo(&jogo)
			}(ponto.x, ponto.y)
		}
	}()

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}