from machine import Pin
from time import sleep
import dht 
import pyb

#https://randomnerdtutorials.com/esp32-esp8266-dht11-dht22-micropython-temperature-humidity-sensor/

sensor = dht.DHT22(Pin(14))
led = pyb.LED(4)
servo = pyb.Servo(1)

while True:
  try:
    sleep(2)
    sensor.measure()
    temp = sensor.temperature()
    hum = sensor.humidity()
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
