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
