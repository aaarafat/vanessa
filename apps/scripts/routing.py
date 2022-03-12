import sys
import os
from turtle import position
import time
from mininet.log import setLogLevel, info
from mn_wifi.link import wmediumd, adhoc
from mn_wifi.cli import CLI
from mn_wifi.net import Mininet_wifi
from mn_wifi.vanet import vanet
from mn_wifi.wmediumdConnector import interference

# print(sys.argv[1])
# info(sys.argv[0])
myIP = sys.argv[1]

while True:
    time.sleep(4)
    f = open('ips') 
    netIPs = f.read().splitlines() 
    f.close()
    # print(len(netIPs))
    netIPs.remove(myIP)
    # print(len(netIPs))
    neighbours = []
    for i in range(len(netIPs)):
        response = os.system("ping -c 5 " + netIPs[i])
        if response==0:
            neighbours.append(netIPs[i])
    # print(myIP, 'neghbours:',neighbours)
    for i in range(len(neighbours)):
        netIPs.remove(neighbours[i])
    for i in netIPs:
        if len(neighbours) == 0:
            break
        os.system(f"ip route add {i} via {neighbours[0]}")
