
import socketio

HOST = "127.0.0.1"  
PORT = 65432 

sockets = {"1" : "a7a","2" : "555555"}
x = "3"
sockets[x] = "zby"
print(f"{x}-asda")
for soc in sockets:
    print(soc)
    
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