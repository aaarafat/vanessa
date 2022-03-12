
import socketio

HOST = "127.0.0.1"  
PORT = 65432 

s = socketio.Client()
url='http://'+HOST+':'+str(PORT)
print(url)
s.connect(url)

@s.on('testResponse')
def testResponse(message):
    print(message["data"])


while True :
    i = input("->")
    if (i == "s") :
        m = input("Enter your message :")
        s.emit('test', {'data': m})