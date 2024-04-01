package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Tamaño del tablero del juego
const boardSize = 2

// Estructura del juego
type Game struct {
	board     [][]bool // Tablero de juego
	boatAlive bool     // Indica si el barco está vivo o no
}

// Inicializa el juego colocando el barco en una posición aleatoria en el tablero.
func (game *Game) Initialize() {
	// Inicializa el tablero como una matriz booleana
	game.board = make([][]bool, boardSize)
	for i := range game.board {
		game.board[i] = make([]bool, boardSize)
	}

	// Genera una semilla aleatoria para los números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Coloca el barco en una posición aleatoria del tablero
	x := rand.Intn(boardSize)
	y := rand.Intn(boardSize)
	game.board[x][y] = true

	// Indica que el barco está vivo al inicio del juego
	game.boatAlive = true
}

// Dispara al tablero en la posición dada (x, y) y devuelve true si golpea el barco.
func (game *Game) Fire(x, y int) bool {
	// Verifica si las coordenadas están dentro del tablero
	if x < 0 || x >= boardSize || y < 0 || y >= boardSize {
		return false
	}
	// Verifica si hay un barco en las coordenadas dadas
	if game.board[x][y] {
		// Si hay un barco, se marca como destruido y devuelve true
		game.boatAlive = false
		return true
	}
	// Si no hay un barco, devuelve false
	return false
}

func main() {
	// Inicia el servidor y espera conexiones de clientes en el puerto 8080.
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close() // Asegura que se cierre el servidor al finalizar el programa.

	fmt.Println("Servidor iniciado. Esperando clientes...")

	for {
		// Acepta conexiones de clientes.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println("Cliente conectado:", conn.RemoteAddr())

		// Inicializa un nuevo juego y maneja la conexión con el cliente.
		game := Game{}
		game.Initialize()
		handleConnection(conn, &game)

		// Si el barco es destruido, muestra un mensaje y termina el programa.
		if !game.boatAlive {
			fmt.Println("¡El barco fue impactado! Cerrando el programa...")
			os.Exit(0)
		}
	}
}

// Maneja la conexión con el cliente.
func handleConnection(conn net.Conn, game *Game) {
	defer conn.Close() // Asegura que se cierre la conexión al finalizar la función.

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// Lee las coordenadas enviadas por el cliente y las procesa.
		coordinates := strings.Split(scanner.Text(), ",")
		if len(coordinates) != 2 {
			// Si el formato de las coordenadas es incorrecto, se informa al cliente.
			fmt.Fprintln(conn, "Entrada inválida. Por favor, proporcione coordenadas en el formato 'x,y'.")
			continue
		}

		x, err := strconv.Atoi(coordinates[0])
		if err != nil {
			// Si la coordenada x no es un número válido, se informa al cliente.
			fmt.Fprintln(conn, "Entrada inválida. Por favor, proporcione coordenadas enteras válidas.")
			continue
		}
		y, err := strconv.Atoi(coordinates[1])
		if err != nil {
			// Si la coordenada y no es un número válido, se informa al cliente.
			fmt.Fprintln(conn, "Entrada inválida. Por favor, proporcione coordenadas enteras válidas.")
			continue
		}

		// Realiza un disparo en las coordenadas dadas y envía la respuesta al cliente.
		if game.Fire(x, y) {
			fmt.Fprintln(conn, "¡Impacto!")
			return
		} else {
			fmt.Fprintln(conn, "¡Fallo!")
		}
	}
}
