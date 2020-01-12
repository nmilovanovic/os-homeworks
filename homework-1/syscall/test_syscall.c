#include <stdio.h>
#include <linux/kernel.h>
#include <sys/syscall.h>
#include <unistd.h>
#include <stdlib.h>

int main(int argc, char ** argv){
    int pid = atoi(argv[1]);
    long int ret = syscall(335, pid);
    printf("Syscall returned %ld \n", ret);
    char str[512];
    perror(str);
    printf("%s\n",str);
    return 0;
}
