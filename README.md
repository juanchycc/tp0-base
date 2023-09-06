# Ejercicio 5

### Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

---

Se define un protocolo, con mensajes de tipo texto divididos con delimitadores.

Los mensajes se conforman por un Header y un Payload, el Header nos dará la información necesaria para saber procesar el contenido del Payload:

**Estructura del Header:**     

`TIPO_MENSAJE;LARGO;AGENCIA\n`

Todo el header se delimita con un '\n' y los componentes del header con ';'. 

TIPO_DE_MENSAJE: Indica que se esta mandando en el Payload.<br>
LARGO: Indica la cantidad de Caracteres enviados en el Payload.<br>
AGENCIA: Indica el ID del cliente (agencia) que envia la apuesta.<br>

La estructura del Payload podría variar según su tipo, para este ejercicio se utiliza:

**Mensaje SINGLE_BET:** Representa una apuesta individual, el cliente la envia a el servidor.

`NOMBRE;APELLIDO;DOCUMENTO;CUMPLEAÑOS;NUMERO\n`

**Mensaje SUCCESS_BET:** Mensaje enviado de Servidor a Cliente, indica que la apuesta fue procesada. No tiene Payload.

**Ejemplo simple del protocolo:**

![diagrama_protocolo_simple](diagrama_protocolo_simple.png "Protocolo Simple")


Para evitar casos de short read se cruza el largo que debería haber llegado según el header y el largo real recibido, en caso de no ser iguales, se esperara un nuevo mensaje que será concatenado con este.

En el caso de sort write se chequea que la cantidad enviada sea la esperada sino se envia el Header (con el largo actualizado) más lo que quedo pendiente por enviarse.

