import logging
from common.utils import Bet, store_bets

BET_FIELDS = 7
SINGLE_BET_TYPE = "SINGLE_BET"

class Loteria:
  def __init__(self):
    self.bet = None
  
  def add_bets( self, msg ):
    
    fields =msg.split(';')
    if( len(fields) != BET_FIELDS ): 
      logging.info(f'action: add_bets | result: fail | error: Number of fields incomplete')
      return
    type, agency, name, last_name, document, birthday, number = fields
    
    if( type != SINGLE_BET_TYPE ): 
      logging.info(f'action: add_bets | result: fail | error: Wrong type of Message: {type} expected {SINGLE_BET_TYPE}')
      return
    
    bet = Bet( agency, name, last_name, document, birthday, number)
            
    self.bet = bet
  
  def store_bets( self ) -> bool:
    if( self.bet == None ): 
      logging.error("action: store_bet | result: fail | error: Invalid Bet")
      return False
    
    bets = []
    bets.append( self.bet )
    store_bets( bets )
    logging.info(f'action: apuesta_almacenada | result: success | dni: {self.bet.document} | numero: {self.bet.number}')
    return True
  
  def successMsg( self ):
    return "success" + ";" + self.bet.document + ";" + str(self.bet.number) + "\n"