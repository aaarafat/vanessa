#!/usr/bin/python



import sys
from turtle import position

from mininet.log import setLogLevel, info
from mn_wifi.link import wmediumd, adhoc
from mn_wifi.cli import CLI
from mn_wifi.net import Mininet_wifi
from mn_wifi.vanet import vanet
from mn_wifi.wmediumdConnector import interference

#ap_scan
#configureAdhoc
#configureMacAddr
#get_default_gw
#setAdhocMode
#get_pid_filename
#setAdhocMode
#setConnected
def topology(args):
    "Create a network."
    net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference, 
    autoAssociation = False)

    info("*** Creating nodes\n")
    kwargs = dict()
    # kwargs['range'] = 100
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
    net.setPropagationModel(model="logDistance", exp=4)

    info("*** Configuring wifi nodes\n")
    net.configureWifiNodes()

    info("*** Creating links\n")
    # MANET routing protocols supported by proto:
    # babel, batman_adv, batmand and olsr
    # WARNING: we may need to stop Network Manager if you want
    # to work with babel
    protocols = ['babel', 'batman_adv', 'batmand', 'olsrd', 'olsrd2']
    # kwargs = dict()
    # for proto in args:
    #     if proto in protocols:
    #         kwargs['proto'] = proto

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
    sta1.cmd(f'sudo python3 routing.py {sta1.wintfs[0].ip} &')
    sta2.cmd(f'sudo python3 routing.py {sta2.wintfs[0].ip} &')
    sta3.cmd(f'sudo python3 routing.py {sta3.wintfs[0].ip} &')
    info("*** Running CLI\n")
    CLI(net)

    info("*** Stopping network\n")
    net.stop()


if __name__ == '__main__':
    setLogLevel('info')
    topology(sys.argv)
