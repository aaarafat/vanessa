
from multiprocessing import Process
import threading
from flask_socketio import SocketIO,emit
from flask import Flask

HOST = "127.0.0.1"  
PORT = 65432 
APP = Flask("mininet")


sio = SocketIO(APP)

@sio.on('connect')
def connected():
    print('Connected')

@sio.on('disconnect')
def disconnected():
    print('Disconnected')

@sio.on('test')
def test(message):
    if(message["data"] == "stop"):
        stop()
    print(message["data"])
    emit('testResponse', {'data': message["data"] + " recieved"})


def run():
    """
    This function to run configured uvicorn server.
    """
    sio.run(APP, host=HOST, port=PORT)


def start():
    """
    This function to start a new process (start the UDI Publisher server).
    """
    threading.Thread(target=run, daemon=True).start()






if __name__ == '__main__':
    start()
    print("is it working ?")
    while True:
        x = input("->")