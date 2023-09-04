import logging
from common.utils import Bet, store_bets

MAX_BUFFER = 1024
BET_FIELDS = 5
FINISH_FIELDS = 3
SINGLE_BET_TYPE = "SINGLE_BET"
SUCCESS_BET_TYPE = "SUCCESS_BET"


class Loteria:
    def __init__(self):
        self.bet = None


    def add_bets( self, socket ):
        
        read = True
        lines, type, id, msg = [], "", "", ""
        recBytes = 0
        
        while read:

            newMsg = socket.recv(MAX_BUFFER).decode('utf-8')
            if not newMsg: return False 
            #Obtener Header:
            newLines = newMsg.split('\n')
                
            headerLen = len(newLines[0]) + 1
            recBytes += ( len(newMsg) - headerLen )
                
            type, readBytes, id  = newLines[0].split(';')
            newLines.pop(0)
                
            if lines == []:
                lines = newLines
            else:
                lines.append(newLines)
            msg = msg + newMsg
            #Si no llego todo, sigo leyendo
            read = ( readBytes != str(recBytes) )
            
        if type == SINGLE_BET_TYPE:
            store_single_bet(lines, id)
            self.successMsg( socket, id)
        else:
            logging.info(f'action: add_bets | result: fail | error: Wrong type of Message: {type}')
        
        return False

    def successMsg( self, socket, id ):
                
        msg = SUCCESS_BET_TYPE + ";" + "0" + ";" + id + "\n"         
        socket.send(msg.encode('utf-8'))
    
    
def store_single_bet( msg, id ):

    bet = get_msg_to_bet( msg[0], id )
    if bet == None: return
                    
    bets = []
    bets.append( bet )
    store_bets( bets )
    logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')

    return 
    
def get_msg_to_bet( msg, id ):
    fields = msg.split(';')

    if len(fields) != BET_FIELDS : 
        logging.info(f'action: add_bets | result: fail | error: Number of fields incomplete')
        return None

    name, last_name, document, birthday, number = fields
    return Bet( id, name, last_name, document, birthday, number)
    