// main.go - Loop principal do jogo
package main

import "os"

// Processa eventos do teclado e é a ÚNICA que altera o estado e redesenha.
func processadoraDeEventos(eventos <-chan EventoTeclado, movimentoFantasma <-chan MovimentoFantasma, jogo *Jogo) {
	for {
		select {
		case ev, aberto := <-eventos:
			if !aberto {
				return
			}
			if continuar := personagemExecutarAcao(ev, jogo); !continuar {
				return // ESC
			}
			interfaceDesenharJogo(jogo)

		case mov := <-movimentoFantasma:
			oldX, oldY := jogo.FantasmaX, jogo.FantasmaY 			//teste
			jogoMoverFantasma(jogo, mov.MX, mov.MY)
			if jogo.FantasmaX != oldX || jogo.FantasmaY != oldY {   //teste
				jogo.StatusMsg = "☠ moveu"
			} else { 												//teste
				jogo.StatusMsg = "☠ tentou mover, mas bateu"
			}
			interfaceDesenharJogo(jogo)
		}
	}
}

func main() {
	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// ⚡ Inicializa pontinhos com goroutine
	IniciarPontinhos(&jogo)

	interfaceDesenharJogo(&jogo)
	eventos := make(chan EventoTeclado, 16)
	movimentoFantasma := make(chan MovimentoFantasma, 8)
	encerrarFantasma := make(chan struct{})

	go processadoraDeEventos(eventos, movimentoFantasma, &jogo)
	go iniciarFantasma(movimentoFantasma, encerrarFantasma)

	for {
		ev := interfaceLerEventoTeclado()
		eventos <- ev
		if ev.Tipo == "sair" {
			close(eventos)
			close(encerrarFantasma)
			return
		}
	}
}