package main

import (
	"bufio"
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	pb "github.com/Kendovvul/Ejemplo/Proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessageServiceServer
}

//=======FUNCIONES ASOCIADAS A PROTOBUFFER (llamados remotamente)=======//

/*
	función		:		Apagar
	input		:		<inputs por defecto de las funciones de protobuffer>
	output		:		(Mensaje que se dará de respuesta, Error en caso de no funcionar)
	explicación :		Llama a función ComunicacionDataNodes para enviar mensaje de apagado
						luego apaga el NameNode. Es llamado desde el cliente rebels.
*/

func (s *server) Apagar(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	ComunicacionDataNodes("", 1, "Apagar")
	ComunicacionDataNodes("", 2, "Apagar")
	ComunicacionDataNodes("", 3, "Apagar")
	os.Exit(1)
	return &pb.Message{Body: "DataNodes Apagados"}, nil
}

/*
	función		:		CombineMessage
	input		:		<inputs por defecto de las funciones de protobuffer>
	output		:		(Mensaje que se dará de respuesta, Error en caso de no funcionar)
	explicación :		Llamado desde el cliente combine. Chequea que el formato del mensaje sea
						correcto, selecciona un DataNode aleatoriopara almacenar el mensaje, y
						genera un registro en el DATA.txt del NameNode.
*/

func (s *server) CombineMessage(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	var respuesta = ""
	log.Println("Mensaje recibido de Combine, chequeando...")
	if msg.Body[len(msg.Body)-1] == '\n' {
		log.Println(msg.Body[:len(msg.Body)-1])
	} else {
		log.Println(msg.Body)
	}

	// Chequear que el mensaje de combine cumpla con el formato requerido:
	chequeo := CheckCombMsg(msg.Body)
	if chequeo == 1 {
		respuesta = "Mensaje aceptado"
		// Decidir aleatoriamente (del 1 al 3) en qué DataNode se almacenará el mensaje:
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		dataNodeID := r1.Intn(3) + 1
		log.Println("Mensaje aceptado, registrando y comunicando a DataNode " + strconv.Itoa(dataNodeID))
		// Enviar mensaje a DataNode:
		ComunicacionDataNodes(msg.Body, dataNodeID, "CreateRecord")
		// Guardar respaldo en DATA.txt:
		substrings := strings.Split(msg.Body, " : ")
		tipo := substrings[0]
		msgID := substrings[1]
		strToWrite := tipo + " : " + msgID + " : " + strconv.Itoa(dataNodeID) + "\n"
		writeFile(strToWrite)
	} else if chequeo == 0 {
		respuesta = "Su mensaje no cumple el formato requerido"
	} else {
		respuesta = "El ID de su mensaje ya existe"
	}
	log.Println("===Servidor escuchando en puerto :50050===")
	return &pb.Message{Body: respuesta}, nil
}

/*
	función		:		RebelsMessage
	input		:		<inputs por defecto de las funciones de protobuffer>
	output		:		(Mensaje que se dará de respuesta, Error en caso de no funcionar)
	explicación :		Se encarga de procesar mensajes del cliente "rebels". Si el mensaje
						tiene el formato correcto, se busca en los DataNodes los registros
						que coincidan con en el campo "TIPO" especificado en el mensaje.
*/

func (s *server) RebelsMessage(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	var mensaje = strings.TrimSpace(strings.ToUpper(msg.Body))

	if mensaje == "MILITAR" || mensaje == "LOGÍSTICA" || mensaje == "FINANCIERA" {
		var tipo = mensaje
		if mensaje[len(mensaje)-1] == '\n' {
			tipo = mensaje[:len(mensaje)-1] //tipo sin salto de linea
		}
		log.Println("Mensaje recibido de Rebels, solicitando registros de tipo " + tipo)
		var nodo1, nodo2, nodo3 = false, false, false

		file, err := os.Open("DATA.txt")
		if err != nil {
			panic(err)
		}

		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() && !(nodo1 == true && nodo2 == true && nodo3 == true) {
			if fileScanner.Text() == "" {
				log.Println("Error, linea vacia en DATA.txt")
				continue
			}
			// "arreglo" es la linea de texto de DATA.txt actual, dividida según " : "
			arreglo := strings.Split(fileScanner.Text(), " : ")
			//log.Println("Se encontró un registro de tipo " + arreglo[0])
			//log.Println("El ID del nodo es: " + arreglo[2])
			if strings.TrimSpace(arreglo[0]) == strings.TrimSpace(tipo) {
				if arreglo[2] == "1" {
					nodo1 = true
				} else if arreglo[2] == "2" {
					nodo2 = true
				} else if arreglo[2] == "3" {
					nodo3 = true
				}
			}
		}

		var informacion = ""

		if nodo1 == true {
			msgNodo1 := ComunicacionDataNodes(tipo, 1, "SearchRecord")
			informacion = msgNodo1
			log.Println("Registros solicitados a nodo 1")
		}
		if nodo2 == true {
			msgNodo2 := ComunicacionDataNodes(tipo, 2, "SearchRecord")
			informacion = informacion + msgNodo2
			log.Println("Registros solicitados a nodo 2")
		}
		if nodo3 == true {
			msgNodo3 := ComunicacionDataNodes(tipo, 3, "SearchRecord")
			informacion = informacion + msgNodo3
			log.Println("Registros solicitados a nodo 3")
		}
		var strRespuesta = ""

		file2, err := os.Open("DATA.txt")
		if err != nil {
			panic(err)
		}
		fileScanner2 := bufio.NewScanner(file2)
		infoArreglo := strings.Split(informacion, "\n")
		fileScanner2.Split(bufio.ScanLines)
		for fileScanner2.Scan() {
			if fileScanner2.Text() == "" {
				continue
			}
			lineaDATA := strings.Split(fileScanner2.Text(), " : ")
			for i := 0; i < len(infoArreglo); i++ {
				if infoArreglo[i] == "" {
					continue
				}
				lineaNODO := strings.Split(infoArreglo[i], " : ")
				if strings.TrimSpace(lineaNODO[0]) == strings.TrimSpace(lineaDATA[1]) {
					strRespuesta = strRespuesta + infoArreglo[i] + "\n"
				}
			}
		}
		if strings.TrimSpace(strRespuesta) == ""{
			log.Println("No se han encontrado registros, aviso enviado a Rebels")
		}else{
			log.Println("Se han enviado registros a Rebels")
		}
		log.Println("===Servidor escuchando en puerto :50050===")
		return &pb.Message{Body: strRespuesta}, nil
	} else {
		log.Println("Formato de input incorrecto, aviso enviado a Rebels")
		log.Println("===Servidor escuchando en puerto :50050===")
		return &pb.Message{Body: "Formato de input incorrecto, inténtelo de nuevo."}, nil
	}
}

//=========FUNCIONES LOCALES=========//

/*
	función		:		ComunicacionDataNodes
	input		:		string msg - mensaje que se desea enviar a algún DataNode
						int dataNodeID - identificador del DataNode con el que conectará (1, 2 o 3) define el puerto TCP de conexión
						funToCall - función del DataNode que se desea llamar
	output		:		int respuesta - la respuesta recibida desde el DataNode
	explicación : 		Establece una conexión sincrónica como cliente con alguno de los DataNodes.
						Hace una llamada a alguna de las funciones remotas de los DataNodes, enviando
						el mensaje recibido como parámetro (msg).

*/

func ComunicacionDataNodes(msg string, dataNodeID int, funToCall string) (resp string) {
	// NameNode actua como Cliente para comunicarse con DataNodes
	// Los puertos son 50051 para el datanode 1, 50052 para el datanode 2 y lo mismo para el 3
	port := ":5005" + strconv.Itoa(dataNodeID)
	// Las maquinas son dist130 para datanode 1, dist131 para el 2 y dist132 para el 3
	connS, err := grpc.Dial("dist13"+strconv.Itoa(dataNodeID-1)+".inf.santiago.usm.cl"+port, grpc.WithInsecure())
	if err != nil {
		panic("No se pudo conectar con el servidor, puerto" + port + err.Error())
	}
	defer connS.Close()
	serviceCliente = pb.NewMessageServiceClient(connS)
	resp = ""
	if funToCall == "CreateRecord" {
		respuesta, err := serviceCliente.CreateRecord(context.Background(), &pb.Message{
			Body: msg,
		})
		if err != nil {
			panic("Error al llamar funcion CreateRecord" + err.Error())
		}
		resp = respuesta.Body
		log.Println("Respuesta recibida de dataNode, mensaje recibido:", "\""+respuesta.Body+"\"")
	} else if funToCall == "SearchRecord" {
		respuesta, err := serviceCliente.SearchRecord(context.Background(), &pb.Message{
			Body: msg,
		})
		if err != nil {
			panic("Error al llamar funcion SearchRecord" + err.Error())
		}
		resp = respuesta.Body
	} else if funToCall == "Apagar" {
		serviceCliente.Apagar(context.Background(), &pb.Message{
			Body: "",
		})
		if err != nil {
			panic("Error al llamar funcion SearchRecord" + err.Error())
		}
		resp = ""
	}
	return resp
}

// función writeFile: escribe en el archivo DATA.txt (lo crea si no existe) la string recibida como parámetro.

func writeFile(line string) {
	f, _ := os.OpenFile("DATA.txt", os.O_APPEND|os.O_WRONLY, 0644)
	_, _ = f.WriteString(line)
	f.Close()
}

/*
	función		:		CheckCombMsg
	input		:		string msg - mensaje recibido desde el combine, que se va a chequear
	output		:		int respuesta - -1 : Error, el ID del mensaje ya existe
										0  : Error, el formato del mensaje no cumple lo esperado
										1  : String aceptada
	explicación : 		Verifica que los mensajes recibidos desde el cliente combine cumplan el formato
						<TIPO> : <ID> : <DATA>, y que el tipo sea MILITAR, LOGÍSTICA O FINANCIERA

*/

func CheckCombMsg(msg string) (respuesta int) {
	match, _ := regexp.MatchString("([A-Z]|[ ]|Í)*[ ][:][ ][0-9]*[ ][:][ ]([A-Z]|[ ])*", msg)
	file, err := os.Open("DATA.txt")

	if err != nil {
		panic(err)
	}
	if match {
		substrings := strings.Split(msg, " : ")
		tipo := substrings[0]
		ID := substrings[1]
		if tipo == "MILITAR" || tipo == "LOGÍSTICA" || tipo == "FINANCIERA" {
			fileScanner := bufio.NewScanner(file)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				if fileScanner.Text() == "" {
					log.Println("Error, linea vacia en DATA.txt")
					continue
				}
				arreglo := strings.Split(fileScanner.Text(), " : ")
				if arreglo[1] == ID {
					return -1
				}
			}
		} else {
			return 0
		}
	} else {
		return 0
	}
	return 1
}

var serviceCliente pb.MessageServiceClient
var solicitudes [4]int
var serv *grpc.Server

func main() {
	// Crear servidor RPC en el puerto 50050
	file, _ := os.Create("DATA.txt")
	file.Close()
	listener, err := net.Listen("tcp", ":50050")
	if err != nil {
		panic("La conexion no se pudo crear" + err.Error())
	}
	serv = grpc.NewServer()
	log.Println("===Servidor escuchando en puerto :50050===")
	pb.RegisterMessageServiceServer(serv, &server{})
	if err = serv.Serve(listener); err != nil {
		panic("El server no se pudo iniciar" + err.Error())
	}
}
