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
# ap_scan
# configureAdhoc
# configureMacAddr
# get_default_gw
# setAdhocMode
# get_pid_filename
# setAdhocMode
# setConnected


HOST = "127.0.0.1"
PORT = 65432
APP = Flask("mininet")


sio = SocketIO(APP, cors_allowed_origins="*")
"Create a network."
net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference,
                   autoAssociation=True)
stations = {}
kwargs = dict(wlans=2)


@sio.on('connect')
def connected():
    print('Connected')


@sio.on('disconnect')
def disconnected():
    print('Disconnected')


@sio.on('test')
def test(message):
    print(message["data"])

    emit('testResponse', {'data': message["data"] + " recieved"})


@sio.on('position')
def position(message):
    print("position updated")
    stations["car1"].setPosition("782,120,0")


@sio.on('add_car')
def add_car(data):
    print(f"card id : {data.id}")
    print(f"card position : {data.pos}")
    stations[data.id] = net.addStation(data.id, position=data.pos,
                                       **kwargs)
    net.addLink(stations[data.id], cls=adhoc, intf=f"{data.id}-wlan0",
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+',  **kwargs)
    # TODO save to json file
    # TODO the aodv protocol


def run_socket():
    sio.run(APP, host=HOST, port=PORT)


def save(path, content):
    with open(path, 'w') as f:
        f.write(content)


info("*** Creating nodes\n")

#kwargs['range'] = 100
stations["car1"] = net.addStation('car1', position="50, 100, 0",
                                  **kwargs)
stations["car2"] = net.addStation('car2', position="100, 100, 0",
                                  **kwargs)
stations["car3"] = net.addStation('car3', position="140, 100, 0",
                                  **kwargs)


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
    for proto in args:
        if proto in protocols:
            kwargs['proto'] = proto

    net.addLink(stations["car1"], cls=adhoc, intf='car1-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+',  **kwargs)
    net.addLink(stations["car2"], cls=adhoc, intf='car2-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+',  **kwargs)
    net.addLink(stations["car3"], cls=adhoc, intf='car3-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+', **kwargs)

    #info("*** Plotting network\n")
    #net.plotGraph(max_x=1000, max_y=1000)
    # net.setMobilityModel(time=0, model='RandomDirection', max_x=1000, max_y=1000,
    #                      min_v=10, max_v=100, seed=20)

    info("*** Starting network\n")
    net.build()

    # info("\n*** Addressing...\n")
    # if 'proto' not in kwargs:
    #     sta1.setIP6('2001::1/64', intf="sta1-wlan0")
    #     sta2.setIP6('2001::2/64', intf="sta2-wlan0")
    #     sta3.setIP6('2001::3/64', intf="sta3-wlan0")
    # sta1.cmd("ip route add 10.0.0.3 via 10.0.0.2")
    # sta3.cmd("ip route add 10.0.0.1 via 10.0.0.2")
    # stations["car1"].cmd('sysctl net.ipv4.ip_forward=1')
    # stations["car2"].cmd('echo 1 > /proc/sys/net/ipv4/ip_forward')
    # stations["car3"].cmd('sysctl net.ipv4.ip_forward=1')
    f = open("ips", "w")
    for id in stations:
        f.write(stations[id].wintfs[0].ip+"\n")
    f.close()
    metadata = {'mac': dict(),
                'mac2ip': dict()}
    for id in stations:
        metadata['mac'][stations[id].name] = stations[id].wintfs[0].mac
        metadata['mac2ip'][stations[id].wintfs[0].mac] = stations[id].wintfs[0].ip
    save("/tmp/mn.metadata.json", json.dumps(metadata))
    for id in stations:
        # TODO : remove number of stations so that we can run it on new added cars
        stations[id].cmd(
            f'sudo {os.path.join(os.path.dirname(__file__), "../../dist/apps/router")} {stations[id].name} {len(stations)} &')
            # f'sudo ./nt {stations[id].name} {stations[id].wintfs[0].ip} &')
    info("*** Running CLI\n")
    CLI(net)
    info("*** Stopping network\n")
    net.stop()


if __name__ == '__main__':
    setLogLevel('info')
    threading.Thread(target=run_socket, daemon=True).start()
    topology(sys.argv)
