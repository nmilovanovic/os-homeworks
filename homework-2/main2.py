import Adafruit_DHT
from machine import Pin
from time import sleep
import pyb


sensor = dht.DHT22(Pin(14))
led = pyb.LED(4)
servo = pyb.Servo(1)

while True:
  try:
    sleep(2)
    hum, temp = Adafruit_DHT.read_retry(Adafruit.DHT.DHT22, 4)
    print('Temperature: %3.1f C' %temp)
    print('Humidity: %3.1f %%' %hum)
    if(temp > 45.0):
        led.on()
        servo.angle(90,5000)
    else:
        led.off()
        servo.angle(0,5000)
        
  except OSError as e:
    print('Failed to read sensor.')
