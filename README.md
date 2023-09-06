# Ejercicio 4

### Modificar servidor y cliente para que ambos sistemas terminen de forma graceful al recibir la signal SIGTERM. Terminar la aplicación de forma graceful implica que todos los file descriptors (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag -t utilizado en el comando docker compose down).

Para este ejercicio se debio incorporar al Cliente y Servidor un handler de la señal SIGTERM. 

Esta señal es enviada por docker al hacer un stop a un container, el flag -t setea un tiempo el cual docker espera para handlear esta señal de manera graceful, sino directamente manda un SIGKILL.

Del lado del servidor (python) se creo un handler que se activa con la señal, aqui se setea un flag que detiene el loop del servidor y además se cierra el socket del mismo.

Para el cliente (go) entre cada iteración se chequea si se recibió la señal, en ese caso se cierra la conexion y termina el loop.

Para probar su funcionamiento se levanta todo con el make docker-compose-up y los logs con make docker-compose-logs y en otra terminal se hace un make docker-compose-stop
esperando una salida correcta del programa.