from ctypes import sizeof
import threading
import logging
import socket
import os
import select
from threading import Timer
import sys
import json
import binascii
import struct
import time
import selectors
from tkinter.ttk import Separator
# Global definitions for the Ethernet IEEE 802.3 interface.
# Source: https://github.com/torvalds/linux/blob/master/include/uapi/linux/if_ether.h
ETH_ALEN = 6                # Octets in one ethernet addr
ETH_TLEN = 2                # Octets in ethernet type field
ETH_HLEN = 14               # Total octets in header.
ETH_ZLEN = 60               # Min. octets in frame sans FCS
ETH_DATA_LEN = 1500         # Max. octets in payload
ETH_FRAME_LEN = 1514        # Max. octets in frame sans FCS

ETH_P_ALL = 0x0003          # Every packet (be careful!!!)
ETH_P_IP = 0x0800           # Internet Protocol packet
ETH_P_ARP = 0x0806          # Address Resolution packet
ETH_P_802_EX1 = 0x88B5      # Local Experimental Ethertype 1
ETH_P_802_EX2 = 0x88B6      # Local Experimental Ethertype 2
station = sys.argv[1]
number_of_nodes= sys.argv[2]
class aodv(threading.Thread):

    # Constructor
    def __init__(self):
        threading.Thread.__init__(self)
        self.node_id = int(station[-1])
        self.num_nodes = number_of_nodes
        self.broadcast_socket = socket.socket(socket.AF_PACKET, socket.SOCK_DGRAM)
        self.interface = station+'-wlan0'
        self.node_name = station
        self.broadcast_id = 0
        self.seq_num = 0
        self.mac_address = ''
    def get_mac_address(self):
        f= open('/tmp/mn.metadata.json')
        metadata = json.load(f)
        f.close()
        self.mac_address = metadata['mac'][self.node_name]
        return self.mac_address
    def get_ip_from_mac(self, mac):
        f= open('/tmp/mn.metadata.json')
        metadata = json.load(f)
        f.close()
        return metadata['mac2ip'][mac]
    def update_neighbors(self):
        sources = []
        with socket.socket(socket.AF_PACKET, socket.SOCK_RAW, socket.htons(ETH_P_ALL)) as server_socket:
        # Bind the interface
            server_socket.bind((self.interface, 0))
            time1 = time.time()
            time2 = time.time()
            while time2-time1<10:
                    # Receive a frame
                print(self.node_name)
                frame = server_socket.recv(ETH_FRAME_LEN)
                # Extract a header
                header = frame[:ETH_HLEN]
                # Unpack an Ethernet header in network byte order
                dst, src, proto = struct.unpack('!6s6sH', header)
                destination = ':'.join('%02x' % octet for octet in dst)
                source = ':'.join('%02x' % octet for octet in src)
                # Extract a payload
                payload = frame[ETH_HLEN:]
                print(f'dst: {destination}, '
                    f'src: {source}, '
                    f'type: {hex(proto)}, '
                    f'payload: {payload[:4]}...')
                if source not in sources and source != self.mac_address:
                    sources.append(source)
                time.sleep(0.25)
                time2 = time.time()
            return sources      
    def broadcast(self, msg="hi", raw_socket=None):
        with socket.socket(socket.AF_PACKET, socket.SOCK_RAW) as client_socket:
        # Bind an interface
            client_socket.bind((self.interface, 0))
            # Send a frame
            print(self.node_name)
            dst = binascii.unhexlify(''.join(("02:00:00:00:01:00").split(':')))
            src = binascii.unhexlify(''.join(self.get_mac_address().split(':')))
            print(client_socket.sendall(
                # Pack in network byte order
                struct.pack('!6s6sH2s',
                            dst,             # Destination MAC address
                            src,    # Source MAC address
                            ETH_P_802_EX1,                      # Ethernet type
                            'Hi'.encode())))                     # Payload
            print('Sent!')
        # if raw_socket == None:
        #     print("please define a broadcasting socket socket")
        #     return
        # dst = b'\xff\xff\xff\xff\xff\xff'  # destination MAC address
        # src = binascii.unhexlify(''.join(self.get_mac_address().split(':')))  # source MAC address
        # payload = msg.encode()            # payload
        # print(self.node_name+' sending')
        # raw_socket.sendall(struct.pack('!6s6sH2s', dst , src , ETH_P_802_EX1 , payload))
    def run(self):
        raw_socket = socket.socket(socket.AF_PACKET, socket.SOCK_RAW, socket.htons(ETH_P_IP))
        raw_socket.bind((self.interface, 0))
        print(self.get_mac_address())
        # if self.node_name == "sta3":
        #     self.broadcast(raw_socket=raw_socket)
        # else:
        time.sleep(1)
        neighbours = self.update_neighbors()

        print(self.node_name+'\'s neighbours : ', neighbours)
        ips = []
        for neighbour in neighbours:
            ips.append(self.get_ip_from_mac(neighbour))
        print(self.node_name+'\'s neighbours : ', ips)
        raw_socket.close()

aodv_thread = aodv()
aodv_thread.start() 
