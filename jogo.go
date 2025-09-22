// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo   rune
	cor       Cor
	corFundo  Cor
	tangivel  bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa            [][]Elemento // grade 2D representando o mapa
	PosX, PosY      int          // posição atual do personagem
	FantasmaX, FantasmaY  int          // posição atual do fantasma
	UltimoVisitadoFantasma Elemento // elemento que estava na posição do fantasma antes de mover
	UltimoVisitado  Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg       string       // mensagem para a barra de status
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
	Pontinho   = Elemento{'•', CorAmarelo, CorPadrao, false}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado: Vazio}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
				jogo.FantasmaX, jogo.FantasmaY = x, y // registra a posição inicial do fantasma
				jogo.UltimoVisitadoFantasma = Vazio // inicializa o último visitado do fantasma como elemento Vazio
			case Vegetacao.simbolo:
				e = Vegetacao
			case Pontinho.simbolo:   
				e = Pontinho	
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy
	elemento := jogo.Mapa[y][x]

	// Só restaura o último visitado se não for pontinho
	if !EhPontinho(jogo.UltimoVisitado) {
		jogo.Mapa[y][x] = jogo.UltimoVisitado
	} else {
		jogo.Mapa[y][x] = Vazio
	}

	jogo.UltimoVisitado = elemento
	jogo.Mapa[ny][nx] = elemento
}

func jogoMoverFantasma(j *Jogo, dx, dy int) {
    fromX, fromY := j.FantasmaX, j.FantasmaY
    nx, ny := fromX+dx, fromY+dy
    if !jogoPodeMoverPara(j, nx, ny) {
        return
    }

	// correção
	if j.UltimoVisitadoFantasma.simbolo == 0 {
        j.UltimoVisitadoFantasma = Vazio
    }

    // repõe o que estava debaixo do fantasma na célula que ele está saindo
    j.Mapa[fromY][fromX] = j.UltimoVisitadoFantasma

    // guarda o que existe no destino (para repor quando sair de lá depois)
    j.UltimoVisitadoFantasma = j.Mapa[ny][nx]

    // coloca o fantasma na célula de destino
    j.Mapa[ny][nx] = Inimigo

    // atualiza coordenadas
    j.FantasmaX, j.FantasmaY = nx, ny
}

