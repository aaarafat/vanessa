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
import sys
import os
from mininet.node import Controller
from mininet.node import UserSwitch
from mn_wifi.node import UserAP

import time
import math
import glob
import base64
import subprocess

from engineio.payload import Payload


Payload.max_decode_packets = 1024

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
cmds = []

STATIONS_COUNT = 5
RSU_COUNT = 5

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind((HOST, PORT))


sio = SocketIO(APP, cors_allowed_origins="*", async_handlers=False)
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


key_bytes = os.urandom(16)
key = base64.b64encode(key_bytes).decode('utf-8')


def cmd(fn, c, bg=True):
    cmds.append(c.replace('sudo', '').strip())
    if bg:
        fn(c + ' &')
    else:
        fn(c)


@sio.on('connect')
def connected():
    print('Connected')


@sio.on('disconnect')
def disconnected():
    print('Disconnected')


@sio.on('clear')
def clear(is_mn_down=False):
    global cmds
    global stations_car, stations_pool, rsus_pool, ap_rsus

    print('Begin Clearing... ')
    if not is_mn_down:
        for st in stations_car.values():
            st.setPosition("0,0,0")
            stations_pool.append(st)
        for rsu in ap_rsus.values():
            try:
                rsu.setPosition("0,0,0")
            except:
                #! IMPORTANT
                pass
            rsus_pool.append(rsu)
        stations_car = {}
        ap_rsus = {}

    lines = {
        l[0]: l[1] for l in (line.split(maxsplit=1)
                             for c in cmds
                             for line in subprocess.check_output(f'pgrep -u root -a | grep -w \'{c}\'', shell=True).decode('utf-8').splitlines())
    }
    cmds = []
    for pid, cmd in lines.items():
        if 'grep' in cmd:
            continue
        print(f'killing {pid}:{cmd}')
        os.system(f'sudo kill -9 {pid} > /dev/null 2>&1')
    clear_unix_sockets()
    if sio:
        sio.emit('cleared')
    print('Done Clearing.')


@sio.on('destination-reached')
def destination_reached(message):
    try:
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
    except Exception as e:
        if running:
            print(e)


@sio.on('obstacle-detected')
def obstacle_detected(message):
    try:
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
    except Exception as e:
        if running:
            print(e)


@sio.on('add-rsu')
def add_rsu(message):
    try:
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

        port = message['port']

        cmd(rsu.cmd,
            f"sudo dist/apps/rsu -id {id} -key {key} -debug")
        cmd(os.system,
            f"socat TCP4-LISTEN:{port},fork,reuseaddr UNIX-CONNECT:/tmp/rsu{id}.ui.socket")
    except Exception as e:
        if running:
            print(e)


@sio.on('add-car')
def add_car(message):
    try:
        id = message['id']
        if id in stations_car:
            print("Received Update Car ...")
            update_car(message)
            return

        print("Received Add Car ...")
        if len(stations_pool) == 0:
            raise Exception("Pool ran out of stations")

        st = stations_pool.pop(0)
        stations_car[id] = st

        port = message['port']

        coordinates = message["coordinates"]
        position = to_grid(coordinates)
        st.setPosition(position)
        print(position)

        cmd(st.cmd, f"sudo dist/apps/network -id {id} -debug")
        cmd(st.cmd, f"sudo dist/apps/car -id {id} -key {key} -debug")
        cmd(os.system,
            f"socat TCP4-LISTEN:{port},fork,reuseaddr UNIX-CONNECT:/tmp/car{id}.ui.socket")

        payload = {
            'type': 'add-car',
            'data': {
                'coordinates': coordinates,
                'speed': message['speed'],
                'route': message['route'],
            }
        }
        time.sleep(0.5)
        send_to_car(f"/tmp/car{id}.socket", payload)

        # run in a new thread
        time.sleep(0.01)
        thread = threading.Thread(target=recieve_from_car, args=(
            f"/tmp/car{id}write.socket", id), daemon=True)

        running_threads.append(thread)

        thread.start()
    except Exception as e:
        if running:
            print(e)


def update_car(message):
    try:
        id = message['id']
        st = stations_car[id]

        coordinates = message["coordinates"]
        position = to_grid(coordinates)
        st.setPosition(position)
        print(position)

        payload = {
            'type': 'add-car',
            'data': {
                'coordinates': coordinates,
                'speed': message['speed'],
                'route': message['route'],
            }
        }
        send_to_car(f"/tmp/car{id}.socket", payload)
    except Exception as e:
        if running:
            print(e)


@sio.on('update-location')
def update_locations(message):
    try:
        id = message['id']
        if id not in stations_car:
            raise Exception("Car not found")

        coordinates = message["coordinates"]
        position = to_grid(coordinates)
        stations_car[id].setPosition(position)

        # lng, lat = coordinates["lng"], coordinates["lat"]
        # print(f"car {id} moved to {position}, lng: {lng} lat: {lat}")

        payload = {
            'type': 'update-location',
            'data': {
                'coordinates': coordinates,
            }
        }
        send_to_car(f"/tmp/car{id}.socket", payload)
    except Exception as e:
        if running:
            print(e)


def recieve_from_car(car_socket, id):
    global running
    try:
        server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        server.bind(car_socket)
        server.listen(1)
        conn, _ = server.accept()
        print(f"Listening on {car_socket}")
        while running and id in stations_car:
            data = conn.recv(1024)
            if not data:
                continue
            data_json = json.loads(data)
            sio.emit(data_json['type'], data_json)

    except Exception as e:
        if not running or id not in stations_car:
            return
        print(f'recieve_from_car error: {e}')
        pass


def send_to_car(car_socket, payload):
    try:
        client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        client.connect(car_socket)
        client.send(json.dumps(payload).encode('ASCII'))
    except Exception as e:
        if not running:
            return
        print(f'retrying to send car {car_socket}')
        time.sleep(0.5)
        send_to_car(car_socket, payload)
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

    info("*** Clearing unix sockets\n")
    clear_unix_sockets()

    info("\n*** Establishing socket connections\n")
    thread = threading.Thread(target=run_socket, daemon=True)
    running_threads.append(thread)
    thread.start()

    info("*** Running CLI\n")
    CLI(net)

    info("*** Stopping network\n")
    global running, sio
    running = False
    time.sleep(0.5)
    print(f"*** Stopping {len(running_threads)} threads")
    for thread in running_threads:
        thread.join(1)
        info('.')
    info('.\n')
    sio = None

    try:
        net.stop()
    except:
        pass
    clear(is_mn_down=True)
    os.system(f"kill $(lsof -t -i:{IO_PORT}) > /dev/null 2>&1 &")


def clear_unix_sockets():
    remove_file_safe('/tmp/car*.socket')
    remove_file_safe('/tmp/rsu*.socket')


def remove_file_safe(path):
    for f in glob.glob(path):
        try:
            os.remove(f)
        except:
            pass


if __name__ == '__main__':
    setLogLevel('info')
    remove_file_safe('/var/log/vanessa/*')
    topology(sys.argv)
