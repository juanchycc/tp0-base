# Ejercicio 1

### 1 - Modificar la definición del DockerCompose para agregar un nuevo cliente al proyecto

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
# Ejercicio 2

### Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. 
### La configuración a través del archivo correspondiente (config.ini y config.yaml, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: docker volumes).

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

Estos volumenes permiten tener un punto de contacto entre el host y los contenedores, de manera que si hago un cambio a nivel host, se ve plasmado en los contenedores que tengan ese volumen.

Ademas de los volumenes fue necesario extraer del Dockerfile del cliente la linea:

```
-COPY ./client/config.yaml /config.yaml
```

Ya que sino se estaría pisando la version del archivo de configuración.

Para probar su funcionamiento, levante los contenedores con el make docker-compose-up, luego realice modificaciones en los archivos de configuración e 
hice un start de los contenedores que se habían buildeado.
---
# Ejercicio 3

### Crear un script que permita verificar el correcto funcionamiento del servidor utilizando el comando netcat para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado. Netcat no debe ser instalado en la máquina host y no se puede exponer puertos del servidor para realizar la comunicación (hint: docker network).

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
