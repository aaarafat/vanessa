#!/usr/bin/python


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


HOST = "127.0.0.1"
PORT = 65432
APP = Flask("mininet")

running = True

LINK_CONFIG = {
    "ssid": 'adhocNet',
    "mode": 'g',
    "channel": 5,
    "ht_cap": 'HT40+'
}

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind((HOST, PORT))


sio = SocketIO(APP, cors_allowed_origins="*")
"Create a network."
net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference,
                   autoAssociation=True, ac_method='llf')  # ssf or llf
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


@sio.on('add-car')
def add_car(message):
    id = f"car{message['id']}"
    position = message["coordinates"]
    # TODO: convert coordinates to mn position
    position = "0,0,0"
    stations[id] = net.addStation(id, position=position,
                                  **kwargs)
    net.addLink(stations[id], cls=adhoc,
                intf=f'{id}-wlan0', **LINK_CONFIG)

    lng = lat = 0  # TODO
    send_location_to_car(f"/tmp/car{id}.socket", lng, lat)


@sio.on('update-locations')
def update_locations(message):
    data = message['data']
    for car_info in data:
        id = car_info['id']
        position = car_info['coordinates']
        # TODO: convert coordinates to mn position
        position = "0,0,0"
        stations[f'car{id}'].setPosition(position)
        lng = lat = 0  # TODO
        send_location_to_car(f"/tmp/car{id}.socket", lng, lat)


def send_location_to_car(car_socket, lng, lat):
    try:
        client = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)
        client.connect(car_socket)
        client.send(json.dumps({'lng': lng, 'lat': lat}).encode('ASCII'))
    except:
        pass


info("*** Creating nodes\n")


def topology(args):
    kwargs = dict(wlans=1)

    net.setPropagationModel(model="logDistance", exp=4)

    info("*** Configuring wifi nodes\n")
    net.configureWifiNodes()

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

    # stations["car1"].cmd('sysctl net.ipv4.ip_forward=1')
    # stations["car2"].cmd('echo 1 > /proc/sys/net/ipv4/ip_forward')
    # stations["car3"].cmd('sysctl net.ipv4.ip_forward=1')

    info("*** Running CLI\n")
    CLI(net)
    info("*** Stopping network\n")
    net.stop()


def run_socket():
    sio.run(APP, host=HOST, port=PORT)


if __name__ == '__main__':
    setLogLevel('info')
    threading.Thread(target=run_socket, daemon=True).start()
    topology(sys.argv)
