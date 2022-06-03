#!/usr/bin/python


import os
import sys
from turtle import position
import threading
from flask_socketio import SocketIO, emit
from flask import Flask



HOST = "127.0.0.1"
PORT = 65432
APP = Flask("mininet")


sio = SocketIO(APP, cors_allowed_origins="*")
"Create a network."
stations = {}
kwargs = dict(wlans=2)


@sio.on('connect')
def connected():
    print('Connected')


@sio.on('disconnect')
def disconnected():
    print('Disconnected')


@sio.on('obstacle-detected')
def obstacle_detected(message):
    print(message)


def run_socket():
    sio.run(APP, host=HOST, port=PORT)


if __name__ == '__main__':
    threading.Thread(target=run_socket, daemon=False).start()
