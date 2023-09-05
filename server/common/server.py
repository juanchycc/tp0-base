from common.loteria import *
import multiprocessing

import socket
import logging
import signal


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._terminated = False
        self._loteria = Loteria()
        self._finished_agencies = 0
        self._client_sockets = []
        signal.signal(signal.SIGTERM, lambda s, _f: self.sigterm_handler( s ) ) 
        
    def sigterm_handler( self, signal ):
        logging.info(f'action: signal_detected | result: success | signal: {signal}')
        self.terminate()

    def terminate( self ):
        self._terminate = True
        self._server_socket.close()
        if len(self._client_sockets) > 0:
            for c in self._client_sockets:
                c.close()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        process = []
        semaphore = multiprocessing.Semaphore()
        agencias = multiprocessing.Value('i', 0)
        
        while not self._terminated:
            client_sock = self.__accept_new_connection()
            if client_sock == None: break
            p = multiprocessing.Process(target=self.__handle_client_connection, args=(client_sock,agencias, semaphore))
            logging.info('MANDO PROCESO')
            p.start()
            process.append(p)
            logging.info('SIGO')
            

            '''with agencias.get_lock():
                if agencias.value == CANTIDAD_AGENCIAS: # type: ignore
                    logging.info('TERMINARON TODAS, ENVIO WIN')
                    for p in process:
                        p.join()
                    logging.info(f'action: sorteo | result: success')
                
                    send_winners(self._client_sockets)
                    self._client_sockets = []'''
                
            
        self.terminate()

    def __handle_client_connection(self, client_sock, agencias, semaphore):
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
            semaphore.acquire()
            agencias.value += 1    
            logging.info(f'AGENCIAS: {agencias.value}')
            semaphore.release()
            #TODO:ver que hacer con los sockets
            #self._client_sockets.append(client_sock)
            esperar_agencias(agencias,semaphore, client_sock)
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

def esperar_agencias(agencias, semaphore, client_sock):
    termina = False
    while not termina:
        semaphore.acquire()
        if agencias.value == CANTIDAD_AGENCIAS: # type: ignore
            logging.info(f'action: sorteo | result: success') 
            send_winners(client_sock)
            termina = True
        semaphore.release()
