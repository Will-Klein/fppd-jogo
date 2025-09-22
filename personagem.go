package main

import "fmt"

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
    dx, dy := 0, 0
    switch tecla {
    case 'w':
        dy = -1 // Move para cima
    case 'a':
        dx = -1 // Move para a esquerda
    case 's':
        dy = 1  // Move para baixo
    case 'd':
        dx = 1  // Move para a direita
    }

    nx, ny := jogo.PosX+dx, jogo.PosY+dy

	if jogoPodeMoverPara(jogo, nx, ny) {
		// aciona botão se a célula de destino for 'B'
		if jogo.Mapa[ny][nx].simbolo == Botao.simbolo {
			botaoPress <- struct{}{}
		}

		// coleta pontinho antes de mover (como já está)
		if EhPontinho(jogo.Mapa[ny][nx]) {
			ColetarPontinho(jogo, nx, ny)
		}

		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
func personagemInteragir(jogo *Jogo) {
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
		personagemInteragir(jogo)
	case "mover":
		personagemMover(ev.Tecla, jogo)
	}
	return true
}
