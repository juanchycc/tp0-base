
import logging
import signal

from common.loteria import send_winners

def child_proccess( loteria, client_sock, barrier ):
  
  child = Child(client_sock)
  
  signal.signal(signal.SIGTERM, lambda s, _f: child.sigterm_handler( s ) )
  child.handle_client_connection(loteria,barrier)

class Child:
    def __init__(self, client_sock):
      self._client_sock = client_sock
      self._terminate = False 

    def sigterm_handler( self, s ):
      self._terminate = True
      self._client_sock.close()
      
    def handle_client_connection(self, loteria, barrier):

      try:
          repeat = True
          while repeat:
              if self._terminate:
                repeat = False
              else:
                repeat = loteria.add_bets( self._client_sock )

  
      except OSError as e:
        logging.error("action: receive_message | result: fail | error: {e}")
        self._client_sock.close()
        return
      finally:  
          self.esperar_agencias(barrier)
      return

    def esperar_agencias(self, barrier):
      
      #Esperar que todas las agencias finalicen con las apuestas
      barrier.wait()
      
      if self._terminate: return
            
      logging.info(f'action: sorteo | result: success') 
      send_winners(self._client_sock)

      self._client_sock.close()