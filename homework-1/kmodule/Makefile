obj-m += kmodule.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules
test:
	sudo insmod kmodule.ko arg_pid=$(PID)
	sudo rmmod kmodule
	dmesg
