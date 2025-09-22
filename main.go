// main.go - Loop principal do jogo
package main

import "os"

// Processa eventos do teclado e é a ÚNICA que altera o estado e redesenha.
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
            oldX, oldY := jogo.FantasmaX, jogo.FantasmaY
            jogoMoverFantasma(jogo, mov.MX, mov.MY)
            if jogo.FantasmaX != oldX || jogo.FantasmaY != oldY {
                jogo.StatusMsg = "☠ moveu"
            } else {
                jogo.StatusMsg = "☠ tentou mover, mas bateu"
            }
            interfaceDesenharJogo(jogo)

        // reaparecimento de pontinhos controlado pela processadora
        case p := <-pontinhosProntos:
            jogo.Mapa[p.y][p.x] = Pontinho
            jogo.StatusMsg = "• reapareceu"
            interfaceDesenharJogo(jogo)

        // ações da porta (true = abrir, false = fechar)
        case abrir := <-portaAcoes:
            if abrir {
                jogo.Mapa[portaY][portaX] = Vazio
                jogo.StatusMsg = "🚪 abriu"
            } else {
                jogo.Mapa[portaY][portaX] = Parede
                jogo.StatusMsg = "🚪 fechou (timeout)"
            }
            interfaceDesenharJogo(jogo)
        }
    }
}

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

	IniciarPorta()
	IniciarPontinhos()

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	eventos := make(chan EventoTeclado, 16)
	movimentoFantasma := make(chan MovimentoFantasma, 8)
	encerrarFantasma := make(chan struct{}) // canal para encerrar a goroutine do fantasma

	// Inicia a goroutine processadora (única dona do estado)
	go processadoraDeEventos(eventos, movimentoFantasma, &jogo)
	go iniciarFantasma(movimentoFantasma, encerrarFantasma)

	// Loop principal de entrada
	for {
		ev := interfaceLerEventoTeclado()

		// Envia o evento para a processadora aplicar e redesenhar
		eventos <- ev

		// Se for para sair, fechamos o canal e encerramos a main
		if ev.Tipo == "sair" {
			close(eventos)          // processadora sai do range e termina
			close(encerrarFantasma) // sinaliza para a goroutine do fantasma encerrar
			return
		}
	}
}