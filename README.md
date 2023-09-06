# Ejercicio 8

### Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.

---

Para la implementación de este ejercicio hice uso de la libreria de multiprocessing de python. Se creará un proceso para atender a cada Agencia (cliente), si bien los procesos son mas pesados, dado el alcance del ejercicio (5 agencias) no se consideraron problemas de espacio, además se favorece el paralelismo.

Como mecanismos de sincronización se utilizaron:

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