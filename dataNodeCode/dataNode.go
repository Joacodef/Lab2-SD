package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	pb "github.com/Kendovvul/Ejemplo/Proto"
	"google.golang.org/grpc"
)

var PORT = "" //Seleccionar puerto 50051, 50052 o 50053

type server struct {
	pb.UnimplementedMessageServiceServer
}

func writeFile(line string) {
	f, _ := os.OpenFile("DATA.txt", os.O_APPEND|os.O_WRONLY, 0644)

	_, _ = f.WriteString(line)

	f.Close()
}

func (s *server) Apagar(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	log.Println("Apagando dataNode...")
	os.Exit(1)
	return &pb.Message{Body: ""}, nil
}

func (s *server) CreateRecord(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	log.Println("Solicitud de NameNode recibida, mensaje registrado: ", msg.Body)
	writeFile(msg.Body)
	log.Println("Enviando respuesta...\nServidor escuchando en puerto " + PORT)
	return &pb.Message{Body: "Registro creado correctamente"}, nil
}

func (s *server) SearchRecord(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	log.Println("Solicitud de NameNode recibida, se requiere buscar: ", msg.Body)
	//BUSCAR REGISTROS DEL TIPO ESPECIFICADO EN msg.Body
	file, err := os.Open("DATA.txt")
	if err != nil {
		panic(err)
	}

	var respuesta = "\n"
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var i = 0

	for fileScanner.Scan() {
		if fileScanner.Text() == "" {
			continue
		}

		arreglo := strings.Split(fileScanner.Text(), " : ")
		if strings.TrimSpace(arreglo[0]) == strings.TrimSpace(msg.Body) {
			i++
			respuesta = respuesta + arreglo[1] + " : " + arreglo[2] + "\n"
		}
	}

	log.Println("Se han enviado " + strconv.Itoa(i) + " registros como respuesta")
	log.Println("===Servidor escuchando en puerto " + PORT + "===")
	return &pb.Message{Body: respuesta}, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var serviceCliente pb.MessageServiceClient
var solicitudes [4]int
var serv *grpc.Server

func main() {
	file, _ := os.Create("DATA.txt")
	file.Close()
	PORT = ":"+strings.TrimSpace(os.Args[1])

	listener, err := net.Listen("tcp", PORT) //conexion sincrona
	if err != nil {
		panic("La conexion no se pudo crear" + err.Error())
	}
	serv = grpc.NewServer()
	log.Println("Servidor escuchando en puerto " + PORT)
	pb.RegisterMessageServiceServer(serv, &server{})
	if err = serv.Serve(listener); err != nil {
		panic("El server no se pudo iniciar" + err.Error())
	}
}
