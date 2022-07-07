import sys
import socket

HOST = "127.0.0.1"  
PORT = 65432 

def controller():
    message = ''
    while message != 'q' and message != 'exit':
        s = socket.socket()
        s.connect((HOST, PORT))
        print(sys.version_info.major)
        message = input('-> ')
        s.send(str(message).encode('utf-8'))
        data = s.recv(1024).decode('utf-8')
        print('Received from server: ' + data)
        s.close()


if __name__ == '__main__':
    controller()