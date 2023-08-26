echo "Iniciando Test"
for i in $(seq 1 3) 
do
  if echo "Testing server" | nc -w 1 server 12345 | grep -q "Testing server"; then
    echo "Server Ok"
  else
    echo "Server Err"
  fi
  sleep 1.5
done
echo "Test Finalizado"
