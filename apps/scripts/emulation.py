import http
from numpy import stack
import socket
import json
from mn_wifi.wmediumdConnector import interference
from mn_wifi.vanet import vanet
from mn_wifi.net import Mininet_wifi
from mn_wifi.cli import CLI
from mn_wifi.link import wmediumd, adhoc
from mininet.log import setLogLevel, info
from flask import Flask
from flask_socketio import SocketIO, emit
import threading
from turtle import position
import sys
import os
from mininet.node import Controller
from mininet.node import UserSwitch
from mn_wifi.node import UserAP

import time
import shutil
import math
from http import client
import glob

from engineio.payload import Payload

Payload.max_decode_packets = 50

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
rsus_pool = []
stations_car = {}
ap_rsus = {}

running_threads = []

STATIONS_COUNT = 5
RSU_COUNT = 5

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind((HOST, PORT))


sio = SocketIO(APP, cors_allowed_origins="*")
"Create a network."
net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference,
                   controller=Controller, accessPoint=UserAP, autoAssociation=True, ac_method='ssf')  # ssf or llf
stations = {}
accessPoints = {}
car_kwargs = dict(
    wlans=2,
    bgscan_threshold=-60,
    s_inverval=5,
    l_interval=10
)


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


@sio.on('add-rsu')
def add_rsu(message):
    if len(rsus_pool) == 0:
        raise Exception("Pool ran out of stations")

    id = message['id']
    if id not in ap_rsus:
        rsu = rsus_pool.pop(0)
        ap_rsus[id] = rsu
    else:
        rsu = ap_rsus[id]

    coordinates = message["coordinates"]
    rsu_range = message["range"]

    position = to_grid(coordinates)
    try:
        rsu.setRange(rsu_range)
    except:
        #! IMPORTANT
        pass
    try:
        rsu.setPosition(position)
    except:
        #! IMPORTANT
        pass
    print(position)

    rsu.cmd(f"sudo dist/apps/network -id {id} -name rsu -debug &")


@sio.on('add-car')
def add_car(message):
    print("Received Add Car ...")
    if len(stations_pool) == 0:
        raise Exception("Pool ran out of stations")

    id = message['id']
    if id not in stations_car:
        st = stations_pool.pop(0)
        stations_car[id] = st
    else:
        st = stations_car[id]

    coordinates = message["coordinates"]
    position = to_grid(coordinates)
    st.setPosition(position)
    print(position)

    st.cmd(f"sudo dist/apps/network -id {id} -name car -debug &")
    st.cmd(f"sudo dist/apps/car -id {id} -debug &")

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
    thread = threading.Thread(target=recieve_from_car, args=(
        f"/tmp/car{id}write.socket",), daemon=True)

    running_threads.append(thread)

    thread.start()


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


def recieve_from_car(car_socket):
    global running
    try:
        server = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
        server.bind(car_socket)
        print(f"Listening on {car_socket}")
        while running:
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


def run_socket():
    sio.run(APP, host=HOST, port=IO_PORT)


def save(path, content):
    with open(path, 'w') as f:
        f.write(content)


EARTH_RAD = 6371 * 1000


def to_grid(coordinates):
    lng, lat = coordinates["lng"], coordinates["lat"]
    lng += 360 if lng < 0 else 0
    lat += 360 if lat < 0 else 0

    x = EARTH_RAD * lng * math.pi/180
    y = EARTH_RAD * lat * math.pi/180

    return f"{x},{y},0"


info("*** Creating nodes\n")


c1 = net.addController('c1')
s0 = net.addSwitch("s0", cls=UserSwitch, inNamespace=True)


def topology(args):

    info("*** Configuring Propagation Model\n")
    net.setPropagationModel(model="logDistance", exp=4)

    for i in range(STATIONS_COUNT):
        stations_pool.append(net.addStation(
            f'car{i + 1}', position="0,0,0", **car_kwargs))

    for i in range(RSU_COUNT):
        rsus_pool.append(net.addAccessPoint(f'rsu{i + 1}', ssid='VANESSA', mode='g', channel='1',
                                            failMode="standalone", position='0,0,0', range=100,
                                            ip=f'10.1.0.{i + 1}/16', cls=UserAP, inNamespace=True))

    net.configureWifiNodes()

    for i, st in enumerate(stations_pool):
        net.addLink(st, cls=adhoc,
                    intf=f'car{i + 1}-wlan0', **LINK_CONFIG)
        st.setIP(ip=f'10.0.1.{i + 1}/24', intf=f'car{i + 1}-wlan1')

    for i, rsu in enumerate(rsus_pool):
        net.addLink(s0, rsu)

    info("*** Starting network\n")
    net.build()
    c1.start()
    s0.start([c1])
    for i, rsu in enumerate(rsus_pool):
        rsu.start([c1])

    s0.cmd('sudo apps/scripts/switch -debug &')

    info("\n*** Establishing socket connections\n")
    for f in glob.glob('/tmp/car*.socket'):
        try:
            os.remove(f)
        except:
            pass
    thread = threading.Thread(target=run_socket, daemon=True)
    running_threads.append(thread)
    thread.start()

    info("*** Running CLI\n")
    CLI(net)
    info("*** Stopping network\n")
    global running
    running = False
    time.sleep(0.5)
    for thread in running_threads:
        thread.join(0.1)

    try:
        net.stop()
    except:
        pass


if __name__ == '__main__':
    setLogLevel('info')
    shutil.rmtree('/var/log/vanessa', ignore_errors=True)
    topology(sys.argv)
