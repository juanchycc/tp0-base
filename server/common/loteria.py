import logging
from common.utils import Bet, has_won, load_bets, store_bets

CANTIDAD_AGENCIAS = 5

MAX_BUFFER = 8192
BET_FIELDS = 5
FINISH_FIELDS = 3
SINGLE_BET_TYPE = "SINGLE_BET"
MULTIPLE_BET_TYPE = "MULTIPLE_BET"
FINISH_BET_TYPE = "FINISH_BET"
SUCCESS_BET_TYPE = "SUCCESS_BET"
WINNERS_BET_TYPE = "WINNERS_BET"

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
          
    if type == MULTIPLE_BET_TYPE:
      store_multiple_bet( lines, id )
      self.successMsg( socket, id)
      return True
    elif type == SINGLE_BET_TYPE:
      store_single_bet(msg, id)
      self.successMsg( socket, id)
    elif type == FINISH_BET_TYPE:
      logging.info(f'action: finish_multiple_bet | result: success')
    else:
      logging.info(f'action: add_bets | result: fail | error: Wrong type of Message: {type}')
    
    return False

  def successMsg( self, socket, id ):
            
    msg = SUCCESS_BET_TYPE + ";" + "0" + ";" + id + "\n"         
    socket.send(msg.encode('utf-8'))
  
  
def store_single_bet( msg, id ):
  
  bet = get_msg_to_bet( msg, id )
  if bet == None: return
            
  bets = []
  bets.append( bet )
  store_bets( bets )
  logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')

  return 

def store_multiple_bet( lines, id ):

  bets = []
 
  for l in lines:
    if l == "": break
    bet = get_msg_to_bet( l ,id )
    if bet == None: return
    bets.append( bet )
  
  store_bets(bets)
  
def get_msg_to_bet( msg, id ):
    fields = msg.split(';')

    if len(fields) != BET_FIELDS : 
      logging.info(f'action: add_bets | result: fail | error: Number of fields incomplete')
      return None

    name, last_name, document, birthday, number = fields
    return Bet( id, name, last_name, document, birthday, number)
  
def get_winners( ):
  
  agencias = {}
  for i in range(1, CANTIDAD_AGENCIAS + 1):
    agencias[i] = 0
  
  bets = list(load_bets())
  
  for b in bets:
    if has_won(b):
      agencias[b.agency] += 1
  return agencias

def send_winners( sockets ):
  
  winners = get_winners()
  msg = ""

  for agencia, cantidad in winners.items():
      msg += f"{agencia};{cantidad}\n"
      
  header = WINNERS_BET_TYPE + ';' + str(len(msg)) + '\n'
  
  packet = header + msg
  
  for s in sockets:
    s.send(packet.encode('utf-8'))
    s.close()
  
  
    