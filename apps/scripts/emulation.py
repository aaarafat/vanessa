#!/usr/bin/python


import sys
import threading
from flask_socketio import SocketIO,emit
from flask import Flask

from mininet.log import setLogLevel, info
from mn_wifi.link import wmediumd, adhoc
from mn_wifi.cli import CLI
from mn_wifi.net import Mininet_wifi
from mn_wifi.wmediumdConnector import interference




HOST = "127.0.0.1"  
PORT = 65432 
APP = Flask("mininet")


sio = SocketIO(APP,cors_allowed_origins="*")

@sio.on('connect')
def connected():
    print('Connected')

@sio.on('disconnect')
def disconnected():
    print('Disconnected')

@sio.on('test')
def test(message):
    print(message["data"])
    print(sta1.position)
    sta1.setPosition("400,890,0")
    emit('testResponse', {'data': message["data"] + " recieved"})

@sio.on('position')
def position(message):
    print("position updated")
    print(sta1.position)
    sta1.setPosition("782,120,0")


def run_socket():
    sio.run(APP, host=HOST, port=PORT)


"Create a network."
net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference)
info("*** Creating nodes\n")
kwargs = dict()
#kwargs['range'] = 100
sta1 = net.addStation('sta1', ip6='fe80::1',position="900,150,0",
                        **kwargs)
sta2 = net.addStation('sta2', ip6='fe80::2',position="150,50,0",
                        **kwargs)
sta3 = net.addStation('sta3', ip6='fe80::3',position="400,150,0",
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
    net.addLink(sta1, cls=adhoc, intf='sta1-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+', **kwargs)
    net.addLink(sta2, cls=adhoc, intf='sta2-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                **kwargs)
    net.addLink(sta3, cls=adhoc, intf='sta3-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+', **kwargs)

    info("*** Starting network\n")
    net.build()
    info("*** Running CLI\n")
    CLI(net)
    info("*** Stopping network\n")
    net.stop()


if __name__ == '__main__':
    setLogLevel('info')
    threading.Thread(target=run_socket, daemon=True).start()
    topology(sys.argv)
