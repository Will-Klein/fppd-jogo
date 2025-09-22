package main

import "time"

// canais do ciclo de vida do pontinho
var reaparecerPontinho = make(chan struct{ x, y int }, 32)
var pontinhosProntos   = make(chan struct{ x, y int }, 32)

// agora essa função só agenda por timer e devolve coords prontas pra processadora
func IniciarPontinhos() {
    go func() {
        for {
            p := <-reaparecerPontinho
            go func(px, py int) {
                time.Sleep(6 * time.Second)
                pontinhosProntos <- struct{ x, y int }{px, py}
            }(p.x, p.y)
        }
    }()
}

// chamada DENTRO da processadora (via personagemMover): remove já
func ColetarPontinho(jogo *Jogo, x, y int) {
    jogo.Mapa[y][x] = Vazio
    reaparecerPontinho <- struct{ x, y int }{x, y}
}

func EhPontinho(e Elemento) bool { return e.simbolo == Pontinho.simbolo }