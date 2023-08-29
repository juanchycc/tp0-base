import logging
from common.utils import Bet, store_bets

BET_FIELDS = 6

class Loteria:
  def __init__(self):
    self.bet = None
  
  def add_bets( self, msg ):
    
    fields =msg.split(';')
    if( len(fields) != BET_FIELDS ): return
    agency, name, last_name, document, birthday, number = fields
    
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