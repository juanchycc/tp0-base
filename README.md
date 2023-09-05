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