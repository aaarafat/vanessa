
import socket

HOST = "127.0.0.1"  
PORT = 65432 

def server():
    s = socket.socket()
    s.bind((HOST, PORT))
    s.listen()
    conn, address = s.accept() 
    print("Connection from: " + str(address))
    while True:
        # receive data stream. it won't accept data packet greater than 1024 bytes
        data = conn.recv(1024).decode()
        print("from connected user: " + str(data))
        data = input(' -> ')
        conn.send(data.encode())  # send data to the client

if __name__ == '__main__':
    server()