#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>

int main(){
 fork();
 fork();
 fork();
 sleep(200);
}
