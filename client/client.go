package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Establece una conexión con el servidor en el puerto 8080.
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close() // Asegura que la conexión se cierre al finalizar el programa.

	// Configura un escáner para leer la entrada del usuario desde la consola.
	scanner := bufio.NewScanner(os.Stdin)

	// Bucle principal para enviar coordenadas al servidor.
	for {
		// Solicita al usuario que ingrese las coordenadas.
		fmt.Print("Ingrese las coordenadas (x,y): ")
		scanner.Scan()                // Escanea la entrada del usuario.
		coordinates := scanner.Text() // Lee las coordenadas ingresadas por el usuario.

		// Envía las coordenadas al servidor a través de la conexión establecida.
		_, err := fmt.Fprintf(conn, "%s\n", coordinates)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Lee la respuesta del servidor.
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Muestra la respuesta del servidor.
		fmt.Println("Respuesta del servidor:", response)

		// Si la respuesta indica que se ha acertado un objetivo, el programa finaliza.
		if response == "¡Impacto!\n" {
			fmt.Println("¡Se ha impactado un barco! ¡Saliendo del programa...")
			return
		}
	}
}
