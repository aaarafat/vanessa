#!/usr/bin/python

# autor: Ramon dos Reis Fontes
# livro: Emulando Redes sem Fio com Mininet-WiFi
# github: https://github.com/ramonfontes/mn-wifi-book-pt

import os

from mininet.node import Controller
from mininet.log import setLogLevel, info
from mn_wifi.node import UserAP
from mn_wifi.cli import CLI
from mn_wifi.net import Mininet_wifi
from mn_wifi.sumo.runner import sumo
from mn_wifi.link import wmediumd, mesh
from mn_wifi.wmediumdConnector import interference


def topology():
    "Create a network."
    net = Mininet_wifi(controller=Controller, accessPoint=UserAP,
                       link=wmediumd, wmediumd_mode=interference)

    info("*** Creating nodes\n")
    cars = []
    for id in range(0, 10):
        cars.append(net.addCar('car%s' % (id+1), wlans=2, bgscan_threshold=-45, 
                               s_inverval=5, l_interval=10, bgscan_module="simple", position=str(500*id)+','+str(500*id)+','+str(0)))
    
    rsu1 = net.addAccessPoint('rsu1', ssid='vanet-ssid', mac='00:00:00:11:00:01',
                            mode='g', channel='1', passwd='123456789a',
                            encrypt='wpa2', position='155,100,0')
    rsu2 = net.addAccessPoint('rsu2', ssid='vanet-ssid', mac='00:00:00:11:00:02',
                            mode='g', channel='6', passwd='123456789a',
                            encrypt='wpa2', position='2320.82,3565.75,0')
    rsu3 = net.addAccessPoint('rsu3', ssid='vanet-ssid', mac='00:00:00:11:00:03',
                            mode='g', channel='11', passwd='123456789a',
                            encrypt='wpa2', position='2806.42,3395.22,0')
    rsu4 = net.addAccessPoint('rsu4', ssid='vanet-ssid', mac='00:00:00:11:00:04',
                            mode='g', channel='1', passwd='123456789a',
                            encrypt='wpa2', position='3332.62,3253.92,0')
    rsu5 = net.addAccessPoint('rsu5', ssid='vanet-ssid', mac='00:00:00:11:00:05',
                            mode='g', channel='6', passwd='123456789a',
                            encrypt='wpa2', position='2887.62,2935.61,0')
    rsu6 = net.addAccessPoint('rsu6', ssid='vanet-ssid', mac='00:00:00:11:00:06',
                            mode='g', channel='11', passwd='123456789a',
                            encrypt='wpa2', position='2351.68,3083.40,0')
    c1 = net.addController('c1')

    info("*** Configuring Propagation Model\n")
    net.setPropagationModel(model="logDistance", exp=2)

    info("*** Configuring wifi nodes\n")
    net.configureWifiNodes()

    net.addLink(rsu1, rsu2)
    net.addLink(rsu2, rsu3)
    net.addLink(rsu3, rsu4)
    net.addLink(rsu4, rsu5)
    net.addLink(rsu5, rsu6)
    for car in cars:
        net.addLink(car, intf=car.params['wlan'][1],
                    cls=mesh, ssid='mesh-ssid', channel=5)

    # net.useExternalProgram(program=sumo, port=8813,
    #                        config_file='map.sumocfg')
    info("*** plotting Network\n")
    net.plotGraph(max_x=1000, max_y=1000)

    info("*** Starting network\n")
    net.build()
    net.addNAT().configDefault()
    c1.start()
    rsu1.start([c1])
    rsu2.start([c1])
    rsu3.start([c1])
    rsu4.start([c1])
    rsu5.start([c1])
    rsu6.start([c1])
    #each car has two interfaces each has an IP address 
    for car in cars:
        car.setIP('192.168.0.%s/24' % (int(cars.index(car))+1),
                  intf='%s-wlan0' % car)
        car.setIP('192.168.1.%s/24' % (int(cars.index(car))+1),
                  intf='%s-mp1' % car)
        car.setRange(10, intf="car%s-mp1"%(int(cars.index(car))+1))

    

    info("*** Running CLI\n")
    CLI(net)

    info("*** Stopping network\n")
    net.stop()


if __name__ == '__main__':
    setLogLevel('info')
    topology()