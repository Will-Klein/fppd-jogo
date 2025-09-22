package main

import "time"

// canais da porta
var (
    botaoPress  = make(chan struct{}, 8) // qualquer 'B' pressiona aqui
    portaAcoes  = make(chan bool, 8)     // true = abrir; false = fechar (timeout)
)

// escolha a posição da porta (uma parede existente no mapa)
const portaX, portaY = 29, 3 // ajuste para uma célula '▤' do teu mapa

// goroutine da porta: abre ao apertar; fecha sozinha após 6s sem novos apertos
func IniciarPorta() {
    go func() {
        aberta := false
        for {
            if !aberta {
                // fechada: espera um botão
                <-botaoPress
                aberta = true
                portaAcoes <- true // abrir
            } else {
                // aberta: timeout para fechar, reinicia se apertarem de novo
                select {
                case <-botaoPress:
                    // apenas reinicia o prazo mantendo aberta
                case <-time.After(6 * time.Second):
                    aberta = false
                    portaAcoes <- false // fechar por timeout
                }
            }
        }
    }()
}