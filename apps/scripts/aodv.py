from ctypes import sizeof
from email import message
from posixpath import split
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
ROUTE_TTL = 200
HELLO_INTERVAL = 7
NEIGHBOR_TTL = 15

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
        self.node_id = int(station[3:])
        self.node_type = station[:3]
        self.num_nodes = number_of_nodes
        # self.broadcast_socket = socket.socket(socket.AF_PACKET, socket.SOCK_DGRAM)
        self.interface = station+'-wlan0'
        self.node_name = station
        self.broadcast_id = 0   #RREQ from Same Node with Same broadcast_id Will Not Be Broadcasted More than Once
        self.seq_num = 0
        self.mac_address = ''
        self.node_ip = ''
        self.neighbors = dict()
        self.port_number = None
        self.socket = None
        self.timer = None
        self.routing_table = dict()

    def get_mac_address(self):
        f= open('/tmp/mn.metadata.json')
        metadata = json.load(f)
        f.close()
        self.mac_address = metadata['mac'][self.node_name]
        self.node_ip = self.get_ip_from_mac(self.mac_address)
        return self.mac_address

    def get_ip_from_mac(self, mac):
        f= open('/tmp/mn.metadata.json')
        metadata = json.load(f)
        f.close()
        return metadata['mac2ip'][mac]
    
    def set_port(self):
        if self.node_type == "rsu":
            self.port_number = 3000 + self.node_id
        elif self.node_type == "sta":
            self.port_number = 4000 + self.node_id
        
    def get_neighbor_type(self, neighbor_ip):
        vals = neighbor_ip.split('.')
        # print(vals)
        if vals[2] == '0':
            return "sta"
        elif vals[2] == '1':
            return "rsu"
        else:
            return None

    def get_neighbor_id(sel, neighbor_ip):
        vals = neighbor_ip.split('.')
        return int(vals[-1])

    def get_port(self, neighbor_ip):
        neighbor_type = self.get_neighbor_type(neighbor_ip)
        neighbor_id = self.get_neighbor_id(neighbor_ip)
        if neighbor_type == "rsu":
            return 3000 + neighbor_id
        elif neighbor_type == "sta":
            return 4000 + neighbor_id

    def send(self, dst_port, msg):
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as client_socket:
                encoded_msg = bytes(msg, 'utf-8')
                client_socket.connect(('localhost', int(dst_port)))
                client_socket.sendall(encoded_msg)
                data = client_socket.recv(1024)
                print("client: ",data)
                # client_socket.send(encoded_msg, 0, 
                #                   ('localhost', int(dst_port)))
        except:
            print(f"failed to send message to {dst_port}")
    
    def create_route(self, route):
        timer = Timer(ROUTE_TTL, 
                      self.handle_dead_neighbor, [route])
        route['Lifetime'] = timer
        timer.start()

    def handle_dead_neighbor(self, neighbor_id):
        del self.neighbors[neighbor_id]
        #delete routes this neighbor is involved in
        
    def send_hello_neighbors(self):
        try:
            for key, neighbor in self.neighbors.items():
                message_type = "HELLO"
                sender_id = str(self.node_id)
                sender_ip = self.node_ip
                message_data = "Hello sent by " + str(sender_id)
                message = message_type + ":" + sender_id + ":" + sender_ip + ":" + message_data
                port = self.get_port(neighbor['ip'])
                print("sending hello")
                self.send(port, message)
                
        
            # Restart the timer
            self.timer.cancel()
            self.timer = Timer(HELLO_INTERVAL, self.send_hello_neighbors, ())
            self.timer.start()
            
        except:
            pass
    def handle_hello(self, hello):
    
        sender_id = hello[1]
        sender_ip = hello[2]
        # Get the sender's ID and restart its neighbor liveness timer
        try:
            if (sender_id in self.neighbors.keys()):
                neighbor = self.neighbors[sender_id]
                timer = neighbor['timer']
                timer.cancel()
                timer = Timer(NEIGHBOR_TTL, 
                              self.handle_dead_neighbor, [sender_id])
                self.neighbors[sender_id] = {'ip': sender_ip, 
                                          'timer': timer}
                timer.start()
            
                # Restart the lifetime timer
                # route = self.routing_table[sender_id]
                # refresh_route(route, False)

            else:
                timer = Timer(HELLO_INTERVAL, 
                              self.send_hello_neighbors, [sender_id])
                self.neighbors[sender_id] = {'ip':sender_ip, 
                                          'Timer-Callback': timer}
                timer.start()
            
                # Update the routing table as well
                if (sender_id in self.routing_table.keys()):
                    route = self.routing_table[sender_id]
                    self.aodv_refresh_route_timer(route, False)
                else:
                    self.routing_table[sender_id] = {'Destination': sender_id, 
                                                  'Destination-Port': self.get_port(sender_id), 
                                                  'Next-Hop': sender_id, 
                                                  'Next-Hop-Port': self.get_port(sender_id), 
                                                  'Seq-No': '1', 
                                                  'Hop-Count': '1'}
                    self.refresh_route(self.routing_table[sender_id], True)

        except KeyError:
            # This neighbor has not been added yet. Ignore the message.
            pass

    

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
    #route request   
    # def send_RREQ(self):
    #     self.broadcast_id +=1
    #     message_type = "RREQ"
    #     destination_type = "rsu"
    #     # Increment our sequence number
    #     self.seq_num += 1
    #     sender_id = self.node_id
    #     sender_ip = self.node_ip
    #     hop_count = 0
    #     rreq_id = self.rreq_id
    #     origin = self.node_id
    #     origin_seq_no = self.seq_num
    #     message = message_type + "-" + str(sender_id) + "-" + sender_ip + "-" + str(hop_count) + "-" + str(rreq_id) + "-" + str(origin) + "-" + str(origin_seq_no)
        
    #     # Broadcast the RREQ packet to all the neighbors
    #     for key, neighbor in self.neighbors.items():
    #         port = self.get_port(neighbor['ip'])
    #         self.send(int(port), message)
    #         logging.debug("['" + message_type + "', 'Broadcasting RREQ to rsu" + "']")
            
    #     # Buffer the RREQ_ID for PATH_DISCOVERY_TIME. This is used to discard duplicate RREQ messages
    #     if (self.node_id in self.rreq_id_list.keys()):
    #         per_node_list = self.rreq_id_list[self.node_id]
    #     else:
    #         per_node_list = dict()
    #     path_discovery_timer = Timer(AODV_PATH_DISCOVERY_TIME, 
    #                                  self.handle_dead_neighbor, 
    #                                  [self.node_id, rreq_id])
    #     per_node_list[rreq_id] = {'RREQ_ID': rreq_id, 
    #                               'Timer-Callback': path_discovery_timer}
    #     self.rreq_id_list[self.node_id] = {'Node': self.node_id, 
    #                                        'RREQ_ID_List': per_node_list}
    #     path_discovery_timer.start()
    #route reply
    def send_RREP(self):
        message_type = "RREP"

    # def broadcast(self, msg="hi", raw_socket=None):
    #     with socket.socket(socket.AF_PACKET, socket.SOCK_RAW) as client_socket:
    #     # Bind an interface
    #         client_socket.bind((self.interface, 0))
    #         # Send a frame
    #         print(self.node_name)
    #         dst = binascii.unhexlify(''.join(("02:00:00:00:01:00").split(':')))
    #         src = binascii.unhexlify(''.join(self.get_mac_address().split(':')))
    #         print(client_socket.sendall(
    #             # Pack in network byte order
    #             struct.pack('!6s6sH2s',
    #                         dst,             # Destination MAC address
    #                         src,    # Source MAC address
    #                         ETH_P_802_EX1,                      # Ethernet type
    #                         'Hi'.encode())))                     # Payload
    #         print('Sent!')
        # if raw_socket == None:
        #     print("please define a broadcasting socket socket")
        #     return
        # dst = b'\xff\xff\xff\xff\xff\xff'  # destination MAC address
        # src = binascii.unhexlify(''.join(self.get_mac_address().split(':')))  # source MAC address
        # payload = msg.encode()            # payload
        # print(self.node_name+' sending')
        # raw_socket.sendall(struct.pack('!6s6sH2s', dst , src , ETH_P_802_EX1 , payload))
    def run(self):
        # raw_socket = socket.socket(socket.AF_PACKET, socket.SOCK_RAW, socket.htons(ETH_P_IP))
        # raw_socket.bind((self.interface, 0))
        print(self.get_mac_address())
        # if self.node_name == "sta3":
        #     self.broadcast(raw_socket=raw_socket)
        # else:
        self.set_port()
        
        time.sleep(1)
        neighbours = self.update_neighbors()
        
        print(self.node_name+'\'s neighbours : ', neighbours)
        ips = []
        for neighbour in neighbours:
            ips.append(self.get_ip_from_mac(neighbour))
            neighbour_id = self.get_neighbor_id(self.get_ip_from_mac(neighbour))
            timer = Timer(NEIGHBOR_TTL, self.handle_dead_neighbor, [neighbour_id])
            self.neighbors[neighbour_id] = {'ip':self.get_ip_from_mac(neighbour), 'timer':timer}
        self.timer = Timer(HELLO_INTERVAL, self.send_hello_neighbors, ())
        self.timer.start()
        print(self.node_name+'\'s neighbours : ', ips)

        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.bind(('localhost', self.port_number))
            s.listen()
            conn, addr = s.accept()
            with conn:
                print(f"Connected by {addr}")
                while True:
                    data = conn.recv(1024)
                    print(data)
                    if not data:
                        break
                    conn.sendall(data)

        # self.socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        # self.socket.bind(('localhost', self.port_number))
        # print("listening to port", self.port_number)
        # self.socket.setblocking(0)
        # self.socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        # outputs = []
        # inputs = [self.socket]
        # # Run the main loop
        # while inputs:
        #     print("RR")
        #     readable, _, _ = select.select(inputs, outputs, inputs)
        #     for r in readable:
        #         print("recieving")
        #         msg, _ = self.socket.recvfrom(200)
        #         decoded = msg.decode('utf-8')
        #         message_lst = decoded.split('-')
        #         print("Recieved ")
        #         message_type = message_lst[0]
        #         if (message_type == "HELLO"):
        #             print("Recieved hello")
        #             self.handle_hello(message_lst)
                # elif (message_type == "RREQ"):
                #     self.aodv_process_rreq_message(message_lst)

            # elif (message_type == "RREP"):
            #     self.aodv_process_rrep_message(message_lst)
            # elif (message_type == "RERR"):
            #     self.aodv_process_rerr_message(message_lst)
            # elif (message_type == "USER_MESSAGE"):
            #     self.aodv_process_user_message(message_lst)
        # raw_socket.close()

aodv_thread = aodv()
aodv_thread.start() 
