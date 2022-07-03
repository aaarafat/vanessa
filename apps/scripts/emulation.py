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


stations["car1"] = net.addStation('car1', position="50, 100, 0",
                                  **car_kwargs)
stations["car2"] = net.addStation('car2', position="100, 100, 0",
                                  **car_kwargs)
stations["car3"] = net.addStation('car3', position="140, 100, 0",
                                  **car_kwargs)


rsu1 = net.addAccessPoint('rsu1', ssid='VANESSA', mode='g', channel='1',
                          failMode="standalone", position='15,70,0', range=100,
                          ip='10.1.0.1/16', cls=UserAP, inNamespace=True)
rsu2 = net.addAccessPoint('rsu2', ssid='VANESSA', mode='g', channel='1',
                          failMode="standalone", position='45,70,0', range=100,
                          ip='10.1.0.2/16', cls=UserAP, inNamespace=True)
rsu3 = net.addAccessPoint('rsu3', ssid='VANESSA', mode='g', channel='1',
                          failMode="standalone", position='75,70,0', range=100,
                          ip='10.1.0.3/16', cls=UserAP, inNamespace=True)

c1 = net.addController('c1')
s0 = net.addSwitch("s0", cls=UserSwitch, inNamespace=True)


def topology(args):

    net.setPropagationModel(model="logDistance", exp=4)

    info("*** Configuring wifi nodes\n")
    net.configureWifiNodes()

    info("*** Creating links\n")
    # MANET routing protocols supported by proto:
    # babel, batman_adv, batmand and olsr
    # WARNING: we may need to stop Network Manager if you want
    # to work with babel
    protocols = ['babel', 'batman_adv', 'batmand', 'olsrd', 'olsrd2']
    kwargs = dict()

    net.addLink(stations["car1"], cls=adhoc, intf='car1-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+',  **kwargs)
    stations["car1"].setIP('10.0.1.1/24',
                           intf='car1-wlan1')
    net.addLink(stations["car2"], cls=adhoc, intf='car2-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+',  **kwargs)
    stations["car2"].setIP('10.0.1.2/24',
                           intf='car2-wlan1')
    net.addLink(stations["car3"], cls=adhoc, intf='car3-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+', **kwargs)
    stations["car3"].setIP('10.0.1.3/24',
                           intf='car3-wlan1')

    net.addLink(s0, rsu1)
    net.addLink(s0, rsu2)
    net.addLink(s0, rsu3)

    # net.addLink(ap1, ap2)
    info("*** Plotting network\n")
    net.plotGraph(max_x=500, max_y=500)
    # net.setMobilityModel(time=0, model='RandomDirection', max_x=1000, max_y=1000,
    #                      min_v=10, max_v=100, seed=20)

    # 802.11b standard defines 13 channels on the 2.4 GHz band at 2.4835 Ghz,
    # allocating 22 MHz for each channel, with a spacing of 5 MHz among them.
    # With this arrangement, only channels 1, 6 and 11 can operate without band
    # overlap.
    info("*** Starting network\n")
    net.build()
    c1.start()
    s0.start([c1])
    rsu1.start([c1])
    rsu2.start([c1])
    rsu3.start([c1])

    # stations["car1"].cmd('sysctl net.ipv4.ip_forward=1')
    # stations["car2"].cmd('echo 1 > /proc/sys/net/ipv4/ip_forward')
    # stations["car3"].cmd('sysctl net.ipv4.ip_forward=1')

    metadata = {'mac': dict(),
                'mac2ip': dict()}
    for id in stations:
        metadata['mac'][stations[id].name] = stations[id].wintfs[0].mac
        metadata['mac2ip'][stations[id].wintfs[0].mac] = stations[id].wintfs[0].ip
    save("/tmp/mn.metadata.json", json.dumps(metadata))
    for id in stations:
        stations[id].cmd(
            f'sudo {os.path.join(os.path.dirname(__file__), "../../dist/apps/router")} {stations[id].name} {len(stations)} &')
    # ap2.cmd(f"./rsuWatcher {ap2.wintfs[0].ssid}")
    # ap1.cmd(f"./rsuWatcher {ap1.wintfs[0].ssid}")

    info("*** Running CLI\n")
    CLI(net)
    info("*** Stopping network\n")
    net.stop()


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
