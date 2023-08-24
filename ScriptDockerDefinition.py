import sys
import os

filename = "docker-compose-dev"
argv = sys.argv

if len( argv ) < 2:
  print( "Error, debe incluir la cantidad de clientes." )
  exit()

try:
  cantidadClientes = int(argv[1]) 
except:
  print("Se esperaba un número que representa la cantidad de clientes.")
  exit() 
 

initialText = """version: '3.9'
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net

"""

clientText = """  client#:
    container_name: client#
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=#
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server

"""

finalText ="""networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

for i in range( 1, cantidadClientes + 1 ):
  initialText = initialText + clientText.replace( '#',str(i) )

with open( filename + '.txt', 'w' ) as f:
    f.write(initialText + finalText)
    
os.rename(filename + '.txt', filename + ".yaml")

print(f"Definición de DockerCompose finalizada con {cantidadClientes} Clientes" )