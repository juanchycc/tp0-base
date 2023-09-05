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