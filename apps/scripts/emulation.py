#!/usr/bin/python


import sys
from turtle import position
import threading
from flask_socketio import SocketIO,emit
from flask import Flask

from mininet.log import setLogLevel, info
from mn_wifi.link import wmediumd, adhoc
from mn_wifi.cli import CLI
from mn_wifi.net import Mininet_wifi
from mn_wifi.vanet import vanet
from mn_wifi.wmediumdConnector import interference
import socketio
import json
#ap_scan
#configureAdhoc
#configureMacAddr
#get_default_gw
#setAdhocMode
#get_pid_filename
#setAdhocMode
#setConnected



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

def save(path, content):
        with open(path, 'w') as f:
            f.write(content)

"Create a network."
net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference, autoAssociation = True)
info("*** Creating nodes\n")
kwargs = dict()
#kwargs['range'] = 100
sta1 = net.addStation('sta1', position="50, 100, 0",
                          **kwargs)
sta2 = net.addStation('sta2', position="100, 100, 0",
                        **kwargs)
sta3 = net.addStation('sta3', position="140, 100, 0",
                        **kwargs)
stations = []
stations.append(sta1)
stations.append(sta2)
stations.append(sta3)


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
                ht_cap='HT40+',  **kwargs)
    net.addLink(sta2, cls=adhoc, intf='sta2-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+',  **kwargs)
    net.addLink(sta3, cls=adhoc, intf='sta3-wlan0',
                ssid='adhocNet', mode='g', channel=5,
                ht_cap='HT40+', **kwargs)
    
    info("*** Plotting network\n")
    net.plotGraph(max_x=1000, max_y=1000)
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
    sta1.cmd('sysctl net.ipv4.ip_forward=1')
    sta2.cmd('echo 1 > /proc/sys/net/ipv4/ip_forward')
    sta3.cmd('sysctl net.ipv4.ip_forward=1')
    f = open("ips", "w")
    for i in range(len(stations)):
        f.write(stations[i].wintfs[0].ip+"\n")
    f.close()
    metadata = {'mac':dict(),
            'mac2ip':dict()}
    for sta in stations:
        metadata['mac'][sta.name] = sta.wintfs[0].mac
        metadata['mac2ip'][sta.wintfs[0].mac] = sta.wintfs[0].ip
    save("/tmp/mn.metadata.json", json.dumps(metadata))
    for sta in stations:
        sta.cmd(f'sudo python3 aodv.py {sta.name} {len(stations)} &')
    info("*** Running CLI\n")
    CLI(net)
    info("*** Stopping network\n")
    net.stop()


if __name__ == '__main__':
    setLogLevel('info')
    threading.Thread(target=run_socket, daemon=True).start()
    topology(sys.argv)
