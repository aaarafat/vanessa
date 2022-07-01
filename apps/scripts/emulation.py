#!/usr/bin/python


import glob
from http import client
import math
import shutil
import time
from mn_wifi.node import UserAP
from mininet.node import Controller
import os
import sys
from turtle import position
import threading
from flask_socketio import SocketIO, emit
from flask import Flask

from mininet.log import setLogLevel, info
from mn_wifi.link import wmediumd, adhoc
from mn_wifi.cli import CLI
from mn_wifi.net import Mininet_wifi
from mn_wifi.vanet import vanet
from mn_wifi.wmediumdConnector import interference
import json

import socket

from numpy import stack


HOST = "127.0.0.1"
IO_PORT = 65432
PORT = 65433
APP = Flask("mininet")

LINK_CONFIG = {
    "ssid": 'adhocNet',
    "mode": 'g',
    "channel": 5,
    "ht_cap": 'HT40+'
}

running = True

stations_pool = []
stations_car = {}

STATIONS_COUNT = 5

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind((HOST, PORT))


sio = SocketIO(APP, cors_allowed_origins="*")
"Create a network."
net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference,
                   autoAssociation=True, ac_method='llf')  # ssf or llf


@sio.on('connect')
def connected():
    print('Connected')


@sio.on('disconnect')
def disconnected():
    print('Disconnected')


@sio.on('destination-reached')
def destination_reached(message):
    id = message['id']
    st = stations_car[id]

    coordinates = message['coordinates']
    position = to_grid(coordinates)
    st.setPosition(position)

    payload = {
        'type': 'destination-reached',
        'data': {
            'coordinates': coordinates,
        }
    }
    send_to_car(f"/tmp/car{id}.socket", payload)


@sio.on('obstacle-detected')
def obstacle_detected(message):
    id = message['id']
    st = stations_car[id]

    coordinates = message['coordinates']
    position = to_grid(coordinates)
    st.setPosition(position)

    obstacle_coordinates = message['obstacle_coordinates']
    payload = {
        'type': 'obstacle-detected',
        'data': {
            'coordinates': coordinates,
            'obstacle_coordinates': obstacle_coordinates
        }
    }
    send_to_car(f"/tmp/car{id}.socket", payload)


@sio.on('add-car')
def add_car(message):
    if len(stations_pool) == 0:
        raise Exception("Pool ran out of stations")

    st = stations_pool.pop(0)
    id = message['id']
    stations_car[id] = st

    coordinates = message["coordinates"]
    position = to_grid(coordinates)
    st.setPosition(position)
    print(position)

    st.cmd(f"sudo apps/scripts/car-unix -id {id} -debug &")

    payload = {
        'type': 'add-car',
        'data': {
            'coordinates': coordinates,
        }
    }
    time.sleep(0.01)
    send_to_car(f"/tmp/car{id}.socket", payload)

    # run in a new thread
    time.sleep(0.01)
    threading.Thread(target=recieve_from_car, args=(
        f"/tmp/car{id}write.socket",)).start()


@sio.on('update-location')
def update_locations(message):
    id = message['id']
    if id not in stations_car:
        raise Exception("Car not found")

    coordinates = message["coordinates"]
    position = to_grid(coordinates)
    stations_car[id].setPosition(position)

    lng, lat = coordinates["lng"], coordinates["lat"]
    # print(f"car {id} moved to {position}, lng: {lng} lat: {lat}")

    payload = {
        'type': 'update-location',
        'data': {
            'coordinates': coordinates,
        }
    }
    send_to_car(f"/tmp/car{id}.socket", payload)


sockets = {}


def recieve_from_car(car_socket):
    try:
        server = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
        server.bind(car_socket)
        print(f"Listening on {car_socket}")
        while True:
            data = server.recv(1024)
            sio.emit('change', data)

    except Exception as e:
        print(f'recieve_from_car error: {e}')
        pass


def send_to_car(car_socket, payload):
    try:
        client = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
        client.connect(car_socket)
        client.send(json.dumps(payload).encode('ASCII'))
    except Exception as e:
        print(f'send_to_car error: {e}')
        pass


info("*** Creating nodes\n")


def topology(args):
    info("*** Configuring Propagation Model\n")
    net.setPropagationModel(model="logDistance", exp=4)

    for i in range(STATIONS_COUNT):
        stations_pool.append(net.addStation(
            f'car{i + 1}', position="0,0,0", wlans=2))

    net.configureWifiNodes()

    for i, st in enumerate(stations_pool):
        net.addLink(st, cls=adhoc,
                    intf=f'car{i + 1}-wlan0', **LINK_CONFIG)

    # info("*** Configuring wifi nodes\n")
    # net.configureWifiNodes()

    # info("*** Plotting network\n")
    # net.plotGraph(max_x=500, max_y=500)
    # net.setMobilityModel(time=0, model='RandomDirection', max_x=1000, max_y=1000,
    #                      min_v=10, max_v=100, seed=20)

    # 802.11b standard defines 13 channels on the 2.4 GHz band at 2.4835 Ghz,
    # allocating 22 MHz for each channel, with a spacing of 5 MHz among them.
    # With this arrangement, only channels 1, 6 and 11 can operate without band
    # overlap.
    info("*** Starting network\n")
    net.build()

    info("*** Running CLI\n")
    CLI(net)
    info("*** Stopping network\n")
    net.stop()


def run_socket():
    sio.run(APP, host=HOST, port=IO_PORT)


EARTH_RAD = 6371 * 1000


def to_grid(coordinates):
    lng, lat = coordinates["lng"], coordinates["lat"]
    lng += 360 if lng < 0 else 0
    lat += 360 if lat < 0 else 0

    x = EARTH_RAD * lng * math.pi/180
    y = EARTH_RAD * lat * math.pi/180

    return f"{x},{y},0"


if __name__ == '__main__':
    setLogLevel('info')
    shutil.rmtree('/logs', ignore_errors=True)
    for f in glob.glob('/tmp/car*.socket'):
        try:
            os.remove(f)
        except:
            pass
    threading.Thread(target=run_socket, daemon=True).start()
    topology(sys.argv)
