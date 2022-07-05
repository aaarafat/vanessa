#!/usr/bin/python


import sys


from mininet.node import Controller
from mn_wifi.node import UserAP
from mininet.log import setLogLevel, info
from mn_wifi.link import wmediumd, adhoc
from mn_wifi.cli import CLI
from mn_wifi.net import Mininet_wifi
from mn_wifi.vanet import vanet
from mn_wifi.wmediumdConnector import interference

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


"Create a network."
net = Mininet_wifi(link=wmediumd, wmediumd_mode=interference,
                   autoAssociation=True, controller=Controller, accessPoint=UserAP)
stations = {}
kwargs = dict()


info("*** Creating nodes\n")

#kwargs['range'] = 100
stations["car1"] = net.addStation('car1', position="250, 300, 0", wlans=2, bgscan_threshold=-60,
                                  s_inverval=5, l_interval=10)
stations["car2"] = net.addStation('car2', position="150, 300, 0", wlans=2, bgscan_threshold=-60,
                                  s_inverval=5, l_interval=10)
stations["car3"] = net.addStation('car3', position="140, 100, 0", wlans=2, bgscan_threshold=-60,
                                  s_inverval=5, l_interval=10)


def topology(args):
    kwargs = dict(wlans=1)

    net.setPropagationModel(model="logDistance", exp=4)
    ap1 = net.addAccessPoint('rsu1', ssid='VANESSA', mode='g', channel='1',
                             failMode="standalone", position='200,300,0', range=100,
                             ip=f'10.1.0.1/16', cls=UserAP, inNamespace=True)

    ap2 = net.addAccessPoint('rsu2', ssid='VANESSA', mode='g', channel='1',
                             failMode="standalone", position='0,50,0', range=100,
                             ip=f'10.1.0.2/16', cls=UserAP, inNamespace=True)

    ap3 = net.addAccessPoint('rsu3', ssid='VANESSA', mode='g', channel='1',
                             failMode="standalone", position='75,100,0', range=100,
                             ip=f'10.1.0.3/16', cls=UserAP, inNamespace=True)
    c1 = net.addController('c1')
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

    s0 = net.addSwitch("s0")

    net.addLink(s0, ap1)
    net.addLink(s0, ap2)
    net.addLink(s0, ap3)

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
    ap1.start([c1])
    ap2.start([c1])
    ap3.start([c1])

    # info("\n*** Addressing...\n")
    # if 'proto' not in kwargs:
    #     sta1.setIP6('2001::1/64', intf="sta1-wlan0")
    #     sta2.setIP6('2001::2/64', intf="sta2-wlan0")
    #     sta3.setIP6('2001::3/64', intf="sta3-wlan0")
    # sta1.cmd("ip route add 10.0.0.3 via 10.0.0.2")
    # sta3.cmd("ip route add 10.0.0.1 via 10.0.0.2")
    # for id in stations:
    #     f.write(stations[id].wintfs[0].ip+"\n")
    # f.close()
    # metadata = {'mac':dict(),
    #         'mac2ip':dict()}
    # for id in stations:
    #     metadata['mac'][stations[id].name] = stations[id].wintfs[0].mac
    #     metadata['mac2ip'][stations[id].wintfs[0].mac] = stations[id].wintfs[0].ip
    # save("/tmp/mn.metadata.json", json.dumps(metadata))
    # for id in stations:
    #     #TODO : remove number of stations so that we can run it on new added cars
    #     stations[id].cmd(f'sudo python3 aodv.py {stations[id].name} {len(stations)} &')
    info("*** Running CLI\n")
    CLI(net)
    info("*** Stopping network\n")
    net.stop()


if __name__ == '__main__':
    setLogLevel('info')
    topology(sys.argv)
