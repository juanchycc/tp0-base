from configparser import ConfigParser
import os

filename = "./tester/test.sh"

def initialize_config():

    config = ConfigParser(os.environ)
    config.read("./tester/config.ini")

    config_params = {}
    try:
        config_params["port"] = int(os.getenv('TESTER_PORT', config["DEFAULT"]["TESTER_PORT"]))
        config_params["reps"] = int(os.getenv('TESTER_REPS', config["DEFAULT"]["TESTER_REPS"]))
        config_params["wait"] = os.getenv('TESTER_WAIT', config["DEFAULT"]["TESTER_WAIT"])
        config_params["timeout"] = int(os.getenv('TESTER_TIMEOUT', config["DEFAULT"]["TESTER_TIMEOUT"]))
        config_params["msg"] = os.getenv('TESTER_MSG', config["DEFAULT"]["TESTER_MSG"])

    except KeyError as e:
        raise KeyError("Key was not found. Error: {} .Aborting server".format(e))
    except ValueError as e:
        raise ValueError("Key could not be parsed. Error: {}. Aborting server".format(e))

    return config_params

def create_bash():
  config_params = initialize_config()
  
  port = config_params["port"]
  msg = config_params["msg"]
  reps = config_params["reps"]
  wait = config_params["wait"]
  timeout = config_params["timeout"]
  
  bash =f"""echo "Iniciando Test"
for i in $(seq 1 {reps}) 
do
  if echo {msg} | nc -w {timeout} server {port} | grep -q {msg}; then
    echo "Server Ok"
  else
    echo "Server Err"
  fi
  sleep {wait}
done
echo "Test Finalizado"
"""


  with open( filename, 'w' ) as f:
    f.write(bash)

def main():
    create_bash()

    os.system('docker build -f ./tester/Dockerfile -t "tester:latest" .')
    os.system('docker compose -f docker-compose-dev.yaml run tester')

if __name__ == "__main__":
    main()