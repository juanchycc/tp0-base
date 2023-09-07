# **Trabajo Práctico 0**
# **Materia:** Sistemas Distribuidos - FIUBA
# **Alumno:** Juan Cruz Caserío

---

## DOCKER

---

## Ejercicio 1

### 1 - Modificar la definición del DockerCompose para agregar un nuevo cliente al proyecto

---

Para poder agregar un nuevo cliente, fue necesario editar el archivo 'docker-compose-dev.yaml' agregando debajo del Client1 un Client2.

El client 2 va a tener un nuevo nombre de container y otro CLI_ID.

### 1.1 - Definir un script (en el lenguaje deseado) que permita crear una definición de DockerCompose con una cantidad configurable de clientes.

Siguiendo con la lógica del ejercicio 1, para este punto programe un Script en Python el cual genera el archivo docker-compose-dev.yaml' dinámicamente según la cantid
de clientes recibidos por parámetros. 

Nuevamente, cada uno tendrá su nombre de contenedor correspondiente asi como el CLI_ID.

Para ejecutarlo, es necesario estar dentro de la carpeta del proyecto, ahi se encontrará el archivo 'ScriptDockerDefinition.py'. 

Ejecutar el siguiente comando definiendo un n (cantidad de clientes):

```
python3 ScriptDockerDefinition.py n
```
---
## Ejercicio 2

### Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (config.ini y config.yaml, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: docker volumes).

---

Para la resolución de este punto primero fue necesario agregar volumenes para el servidor y para el/los cliente/s.

Para eso se sumaron a la definicion de los contenedores las siguientes lineas:

```
#Servidor
volumes:
  - ./server/config.ini:/config.ini

#Cliente
volumes:
  - ./client/config.yaml:/config.yaml
```

Estos volumenes permiten tener un punto de contacto entre el host y los contenedores, de manera que si se realiza un cambio a nivel host, se ve plasmado en los contenedores que tengan ese volumen.

Ademas de los volumenes fue necesario extraer del Dockerfile del cliente la linea:

```
-COPY ./client/config.yaml /config.yaml
```

Ya que sino se estaría pisando la version del archivo de configuración.

Para probar su correcto funcionamiento, levante los contenedores con el make docker-compose-up, luego realice modificaciones en los archivos de configuración e 
hice un start de los contenedores que se habían buildeado viendo si los cambios fueron registrados.

---
## Ejercicio 3

### Crear un script que permita verificar el correcto funcionamiento del servidor utilizando el comando netcat para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado. Netcat no debe ser instalado en la máquina host y no se puede exponer puertos del servidor para realizar la comunicación (hint: docker network).

---

Para la resolución de este ejercicio cree una carpeta dedicada llamada 'tester' aqui se encuentran:
* Archivo Dockerfile el cual levanta la imagen de [subfuzion/netcat](https://hub.docker.com/r/subfuzion/netcat).
* Archivo de configuración para poder definir los parámetros usados por el script a la hora de testear: puerto, ip, repeticiones, esperas entre mensajes, timeout para espera conexion, mensaje a enviar
* Archivo bash generado por el Script, el mismo será ejecutado por el contenedor una vez iniciado.

Para mandar mensajes al servidor se da uso a la docker network 'testing_net' y a netcat como como herramienta para el envio de mensajes.

Con el siguiente codigo bash se envia un msg al servidor y dado que es un echo server se toma la respuesta que se espera sea la misma a la enviada:

```
if echo {msg} | nc -w {timeout} server {port} | grep -q {msg}; then
    echo "Server Ok"
  else
    echo "Server Err"
```

Programe un Script python que genera este código bash a partir de los parámetros de configuración y además buildea el Dockerfile de tester y levanta el contenedor.
Para usarlo (necesario estar dentro de la carpeta del proyecto) ejecutar:

```
python3 ScriptServerTester.py
```
---
## Ejercicio 4

### Modificar servidor y cliente para que ambos sistemas terminen de forma graceful al recibir la signal SIGTERM. Terminar la aplicación de forma graceful implica que todos los file descriptors (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag -t utilizado en el comando docker compose down).

---

Para este ejercicio se debio incorporar al Cliente y Servidor un handler de la señal SIGTERM. 

Esta señal es enviada por docker al hacer un stop a un container, el flag -t setea un tiempo el cual docker espera para handlear esta señal de manera graceful, sino directamente manda un SIGKILL.

Del lado del servidor (python) se creo un handler que se activa con la señal, aqui se setea un flag que detiene el loop del servidor y además se cierra el socket del mismo.

Para el cliente (go) entre cada iteración se chequea si se recibió la señal, en ese caso se cierra la conexion y termina el loop.

Para probar su funcionamiento se levanta todo con el make docker-compose-up y los logs con make docker-compose-logs y en otra terminal se hace un make docker-compose-down
esperando una salida correcta del programa.

---

## COMUNICACIONES

---


## Ejercicio 5

### Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

---

Se define un protocolo, con mensajes de tipo texto divididos con delimitadores.

Los mensajes se conforman por un **Header** y un **Payload**, el Header nos dará la información necesaria para saber procesar el contenido del Payload:

* **Estructura del Header:**     

    `TIPO_MENSAJE;LARGO;AGENCIA\n`

Todo el header se delimita con un '\n' y los componentes del header con ';'. 

- TIPO_DE_MENSAJE: Indica que se esta mandando en el Payload.<br>
- LARGO: Indica la cantidad de Caracteres enviados en el Payload.<br>
- AGENCIA: Indica el ID del cliente (agencia) que envia la apuesta.<br>

La estructura del Payload podría variar según su tipo, para este ejercicio se utiliza:

* **Mensaje SINGLE_BET:** Representa una apuesta individual, el cliente la envia a el servidor.

    `NOMBRE;APELLIDO;DOCUMENTO;CUMPLEAÑOS;NUMERO\n`

* **Mensaje SUCCESS_BET:** Mensaje enviado de Servidor a Cliente, indica que la apuesta fue procesada. No tiene Payload.

**Ejemplo simple del protocolo:**

![diagrama_protocolo_simple](diagrama_protocolo_simple.png "Protocolo Simple")


Para evitar casos de short read se cruza el largo que debería haber llegado según el header y el largo real recibido, en caso de no ser iguales, se esperara un nuevo mensaje que será concatenado con este.

En el caso de short write, para evitarlos se chequea que la cantidad enviada sea la esperada sino se envia el Header (con el largo actualizado) más lo que quedo pendiente por enviarse.

---

## Ejercicio6

### Modificar los clientes para que envíen varias apuestas a la vez 

---

Para este Ejercicio se extiende el protocolo planteado para el Ejercicio 5. Siendo el header identico.

Pero se agregan nuevos tipos de mensajes:

* **MULTIPLE_BET:** Misma idea que SINGLE_BET pero permitiendo tener varias apuestas a la vez, las mismas se separaran por un \n.

    `SINGLE_BET\n`...
`SINGLE_BET\n`

* **FINISH_BET:** Al empezar un envio de multiples bets se utiliza este mensaje para marcar el fin del mismo. No cuenta con Payload.

<p>Por otro lado, se agregan constantes en el codigo que permiten modificar la cantidad de apuestas por mensaje, asi como el buffer maximo.<br>
La idea con el short read y write se mantiene identica al ejercicio 5.<br>
Fue necesario que los archivos csv sean copiados con el dockerfile del cliente para tener acceso a los mismos,y se procesan por su correspondiente agencia linea a linea.</p>


**Ejemplo del protocolo:**

![diagrama_protocolo_multiple](diagrama_protocolo_multiple.png "Protocolo multiple")

---

## Ejercicio 7

### Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo. 

---

Dado el protocolo definido en el Ejercicio 6, las agencias ya notificaban con un **FINISH_BET** la finalización, por lo que solo se agrega el siguiente mensaje:

* **WINNERS_BET:** Mensaje que contiene la cantidad de ganadores en cada agencia, el servidor lo envia a todas las agencias luego de que todas finalicen con el envio de apuestas.

    `WINNERS_BET;LARGO\n` `AGENCIA;CANTIDAD\n...AGENCIA;CANTIDAD\n`


Una vez enviado el FINISH_BET el cliente esperara a recibir el WINNERS_BET, al recibirlo buscara su ID e imprimira la cantidad recibida. <br>
El Servidor contará la cantidad de agencias finalizadas y al llegar al valor esperado (configurable) enviará el mensaje.

---

## CONCURRENCIA


---

## Ejercicio 8

### Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.

---

Para la implementación de este ejercicio hice uso de la libreria de **multiprocessing** de python. <br>
Se creará un proceso para atender a cada Agencia (cliente), si bien los procesos son mas pesados, dado el alcance del ejercicio (5 agencias) no se consideraron problemas de espacio, además se favorece el paralelismo.

Como mecanismos de **sincronización** se utilizaron:

* Lock a la hora de escribir el archivo bets.csv, de esta manera los procesos no se solaparan al momento de guardar los bets.
* Semaforo para acceder al estado compartido del Servidor, el cual implica la cantidad de agencias completas y un flag que indica si se detecto un SIGTERM

Resumiendo el funcionamiento del programa:

- Al llegar una nueva conexión se crea un proceso.
- El proceso recibe las apuestas y responde, entre cada respuesta accede al semaforo para chequear si hubo un SIGTERM.
- Cuando se finaliza el envio de apuestas el proceso lockea el archivo, cuando puede acceder, guarda las apuestas en el.
- Se hace un acquire al semaforo, al acceder se aumenta el numero de agencias procesadas y libera el semaforo.
- Entra en un loop en el cual pide acceso al semaforo y verifica si tuvo un SIGTERM:
  - Si hubo un SIGTERM, cierra el socket con el cliente y finaliza.
  - En caso de que no, verifica si todas las agencias finalizaron, en ese caso envia el mensaje WINNERS_BETS y finaliza, sino libera el lock y sigue loopeando.
