package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/Kendovvul/Ejemplo/Proto"
	"google.golang.org/grpc"
)

var serviceCliente pb.MessageServiceClient

func main() {
	// Conexión con el nameNode (el combine es cliente):
	port := ":50050"
	connS, err := grpc.Dial("localhost"+port, grpc.WithInsecure())
	if err != nil {
		panic("No se pudo conectar con el servidor" + err.Error())
	}
	defer connS.Close()
	serviceCliente = pb.NewMessageServiceClient(connS)

	for {
		// Interfaz para el usuario:
		fmt.Println("Ingrese información: (<TIPO> : <ID> : <DATA>) \nPara salir ecribir 'exit' ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n') // Line = informacion
		if err != nil {
			log.Fatal(err)
		}
		if strings.TrimSpace(line) == "exit" {
			break
		}
		// Llamada remota a NameNode para registrar mensaje:
		respuesta, _ := serviceCliente.CombineMessage(context.Background(), &pb.Message{
			Body: line,
		})
		log.Println(respuesta.Body)

	}

}
