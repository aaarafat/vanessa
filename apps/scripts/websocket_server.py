#!/usr/bin/python

# autor: Ramon dos Reis Fontes
# livro: Emulando Redes sem Fio com Mininet-WiFi
# github: https://github.com/ramonfontes/mn-wifi-book-pt

import os
import time
import threading

from simple_websocket_server import WebSocket, WebSocketServer


# --------------Global Nodes------------- # 
cars = []
# --------------------------------------- #

REFRESH_DELAY = 5  # 5 seconds
WEBSOCKET_PORT = 8000

Running = True

class WebSocketHandler(WebSocket):
  def handle(self):
    """Handle recieved messages
    """
    data = json.loads(self.data)

    if data['type'] == 'position' and data['method'] == 'POST':
        print(f'Recieved position data in POST method:\n {data}\n')
     
    def _send_positions(self):
        try:
            while running:
                nodes = []
                for car in cars:
                    x, y, z = car.getxyz()
                    ip = car.IP()
                    nodes.append({'type': car, 'x': x, 'y': y, 'ip': ip})
                self.send_message(
                    json.dumps({'nodes': nodes})
                )
                time.sleep(REFRESH_DELAY)
        except:
            pass

    def connected(self):
        print(f'client {self.address} has connected\n')
        #threading.Thread(target=self._send_positions, daemon=True).start()

    def handle_close(self):
        print(f'client {self.address} has closed\n')
        Running = False

    @classmethod
    def serve(cls):
        print('Running Socket Server!!!')
        WebSocketServer('', WEBSOCKET_PORT, WebSocketHandler).serve_forever()


if __name__ == '__main__':
  print('Running Socket Server!!!')
  WebSocketServer('', WEBSOCKET_PORT, WebSocketHandler).serve_forever()