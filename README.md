# Laboratorio 2 - Sistemas Distribuidos

# Integrantes:

-Vicente Balbontín 201804585-3

-Joaquín De Ferrari 201804543-8

# Ejecución de archivos:

-En todos los casos, se debe entrar a la carpeta "Lab2-SD" que estará en la carpeta principal "dist"

-Estando en la máquina dist129, utilizar el comando "make namenode"

-Estando en otra terminal en la máquina dist129, utilizar el comando "make combine"

-Estando en la máquina dist130, utilizar el comando "make datanode1"

-Estando en otra terminal en la máquina dist130, utilizar el comando "make rebeldes"

-Estando en la máquina dist131, utilizar el comando "make datanode2"

-Estando en la máquina dist132, utilizar el comando "make datanode3"

(La razón de que se use datanode<1,2,3> es que el numero permite setear el puerto de escucha)

(Es importante que se creen los nodos indicados en las máquinas indicadas para que no haya errores)

-Cuando todas las máquinas estén funcionando, en el combine se pueden ingresar strings con el formato "TIPO : ID : DATA". Notar que en "TIPO" se deben utilizar mayúsculas y la tilde en "LOGÍSTICA", además, el "ID" deben ser sólo números. Una vez escrito, presionar ENTER, tras lo cual se debería recibir una respuesta.

-En el rebels se pueden realizar búsquedas por tipo, ingresando simplemente este campo, con mayúsculas y tilde al igual que antes. Se deberían recibir los mensajes guardados como respuesta, o un mensaje diciendo que no se han encontrado, según sea el caso.

# Explicación general del código:

-Tanto el NameNode como los DataNodes levantan servidores donde están a la espera de mensajes GRPC, en puertos desde el 50050 al 50053.

-Rebels y Combine tienen terminales para recibir input del usuario, y una vez recibido, se comunican mediante GRPC al NameNode (puerto 50050), actuando como clientes.

-Cuando NameNode es llamado desde Combine, selecciona un DataNode al azar y le envía el mensaje para que lo almacene en su archivo DATA.txt. Además, guarda en su propio DATA.txt el tipo, idMensaje y idNodo.

-Cuando el NameNode es llamado desde Rebels, revisa su propio DATA.txt para ver qué nodos contienen el tipo de mensaje que este solicita. Luego obtiene todos los mensajes que sean de este tipo en cada nodo, para luego hacer que coincida el orden en que fueron ingresados originalmente a la hora de entregarlos a Rebels.