#!/usr/bin/python


from ipaddress import IPv4Address
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
from mn_wifi.wmediumdConnector import interference

from mininet.node import Controller
from mn_wifi.node import UserAP
import json
# ap_scan
# configureAdhoc
# configureMacAddr
# get_default_gw
# setAdhocMode
# get_pid_filename
# setAdhocMode
# setConnected
# py car1.wintfs[1].apsInRange
# py car3.wintfs[1].ssid
# py car1.wintfs[1].associatedTo


HOST = "127.0.0.1"
PORT = 65432
APP = Flask("mininet")


sio = SocketIO(APP, cors_allowed_origins="*")
"Create a network."
net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference,
                   controller=Controller, accessPoint=UserAP, autoAssociation=True, ac_method='ssf') #ssf or llf
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
                                       **car_kwargs)
    net.addLink(stations[data.id], cls=adhoc, intf=f"{data.id}-wlan0",
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+')
    # TODO save to json file
    # TODO the aodv protocol


def run_socket():
    sio.run(APP, host=HOST, port=PORT)


def save(path, content):
    with open(path, 'w') as f:
        f.write(content)


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
    threading.Thread(target=run_socket, daemon=True).start()
    topology(sys.argv)
    