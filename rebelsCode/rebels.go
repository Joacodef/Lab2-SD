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
	connS, err := grpc.Dial("dist129.inf.santiago.usm.cl"+port, grpc.WithInsecure())
	if err != nil {
		panic("No se pudo conectar con el servidor" + err.Error())
	}
	defer connS.Close()
	serviceCliente = pb.NewMessageServiceClient(connS)

	for {
		// Interfaz para el usuario:
		fmt.Println("===Ingrese Categoria: MILITAR,LOGÍSTICA,FINANCIERA=== \nPara apagar todo escriba 'exit' ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Llamadas remotas a NameNode:
		if strings.TrimSpace(line) == "exit" {
			// Apagar todo:
			log.Println("Apagando NameNode y DataNodes...")
			serviceCliente.Apagar(context.Background(), &pb.Message{
				Body: "",
			})
			log.Println("Apagado exitoso. ¡Nos vemos!")
			break
		} else {
			// Búsqueda de registros por tipo:
			respuesta, _ := serviceCliente.RebelsMessage(context.Background(), &pb.Message{
				Body: line,
			})
			if strings.TrimSpace(respuesta.Body) == "" {
				log.Println("<REGISTROS NO ENCONTRADOS>")
			} else {
				log.Println("Registros recibidos (ID : DATA):\n" + respuesta.Body)
			}
		}

	}
}
