# Ejercicio 7

### Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo. 

---

Dado el protocolo definido en el Ejercicio 6, las agencias ya notificaban con un FINISH_BET la finalización, por lo que solo se agrega el siguiente mensaje:

**WINNERS_BET:** Mensaje que contiene la cantidad de ganadores en cada agencia, el servidor lo envia a todas las agencias luego de que todas finalicen con el envio de apuestas.

`WINNERS_BET;LARGO\n` `AGENCIA;CANTIDAD\n...AGENCIA;CANTIDAD\n`


Una vez enviado el FINISH_BET el cliente esperara a recibir el WINNERS_BET, al recibirlo buscara su ID e imprimira la cantidad recibida. El Servidor contará la cantidad de agencias finalizadas y al llegar al valor esperado (configurable) enviará el mensaje.