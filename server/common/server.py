from common.loteria import *
import multiprocessing

import socket
import logging
import signal

class ServerStatus:
    def __init__(self):
        self.semaphore = multiprocessing.Semaphore()
        self.rec_signal = multiprocessing.Value('i', False)
        self.agencias_procesadas = multiprocessing.Value('i', 0)

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._terminated = False
        self._loteria = Loteria()
        self._status = ServerStatus()
        signal.signal(signal.SIGTERM, lambda s, _f: self.sigterm_handler( s ) ) 
        
    def sigterm_handler( self, signal ):
        logging.info(f'action: signal_detected | result: success | signal: {signal}')
        self.terminate()

    def terminate( self ):
        self._status.semaphore.acquire()
        self._status.rec_signal.value = True # type: ignore
        self._status.semaphore.release()
        self._terminate = True
        self._server_socket.close()


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        process = []
        
        while not self._terminated:
            client_sock = self.__accept_new_connection()
            if client_sock == None: break
            p = multiprocessing.Process(target=self.__handle_client_connection, args=(client_sock,self._status))
            p.start()
            process.append(p)   
            
        for p in process:
            p.join()
        

    def __handle_client_connection(self, client_sock, status):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:

            repeat = True
            while repeat:
                repeat = self._loteria.add_bets( client_sock )

            addr = client_sock.getpeername()
 
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
            client_sock.close()
        finally:
            status.semaphore.acquire()
            status.agencias_procesadas.value += 1   
            status.semaphore.release()

            esperar_agencias( status, client_sock)
        return
            
    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        try:
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except:
            return None

def esperar_agencias(status, client_sock):
    termina = False
    while not termina:
        status.semaphore.acquire()
        
        if status.rec_signal.value:
            termina = True
        elif status.agencias_procesadas.value  == CANTIDAD_AGENCIAS: # type: ignore
            logging.info(f'action: sorteo | result: success') 
            send_winners(client_sock)
            termina = True
        status.semaphore.release()
    client_sock.close()
