from common.loteria import *
import multiprocessing

import socket
import logging
import signal

from common.child import child_proccess

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._terminated = False
        self._client_socket: socket.socket
        self._childs = []
        self._loteria = Loteria()
        signal.signal(signal.SIGTERM, lambda s, _f: self.sigterm_handler( s ) ) 
        
    def sigterm_handler( self, signal ):
        logging.info(f'action: signal_detected | result: success | signal: {signal}')
        self.terminate()

    def terminate( self ):

        self._terminate = True
        if self._client_socket != None:
            self._client_socket.close()
        #Enviar SIGTERM a los procesos hijos:
        for p in self._childs:
            p.terminate() 

        self._server_socket.close()


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        barrier = multiprocessing.Barrier(CANTIDAD_AGENCIAS)
        
        while not self._terminated:
            client_sock = self.__accept_new_connection()
            if client_sock == None: break
            p = multiprocessing.Process(target= child_proccess, args=(self._loteria, client_sock, barrier))
            p.start()
            self._childs.append(p)   
            
        for p in self._childs:
            p.join()
            
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
            logging.info(
                f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except:
            return None