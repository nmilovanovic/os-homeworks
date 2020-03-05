from socket import *
import random

s = []
k = [1, 2, 3, 4, 5, 6, 7, 8]
for i in range(0,256):
    #s[i] = i
    s.append(i)
j = 0
for i in range(0,256):
    j = (j+s[i]+k[i % 8]) % 256
    p = s[i]
    s[i] = s[j]
    s[j] = p
    j = (j + s[i] + k[i%8])%256

def prng():
    i=0#treba da bude staticka promenljiva
    j=0#treba da bude staticka promenljiva
    i = (i+1)%256
    j = (j+s[i]) % 256
    p = s[i]
    s[i] = s[j]
    s[j] = p
    t = s[i]+s[j]
    t = t % 256
    return s[t]

def get_temp():
    return str(-20+prng()%60)
               
def get_hum():
    return str(80+prng()%18)

serverSocket = socket(AF_INET, SOCK_STREAM)
serverSocket.bind(('', 8080))
serverSocket.listen(1)
while True:
    print('Prihvatam konekciju...')
    connectionSocket, addr = serverSocket.accept()
    message = connectionSocket.recv(1024)
    connectionSocket.send('HTTP/1.0 200 OK\r\n\r\n')
    connectionSocket.send('<html><body><p>Temperature: '+get_temp()+'</p><p>Humidity: '+get_hum()+'</p></body></html>')
    connectionSocket.close()
serverSocket.close() 
