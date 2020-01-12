#include <linux/kernel.h>
#include <linux/sched.h>
#include <linux/list.h>
//https://stackoverflow.com/questions/32189782/using-for-each-process-in-my-program-prevents-me-from-compiling-compiler-says
#include <linux/sched/signal.h>
#include <linux/syscalls.h>

//////////////////////////////
struct stack_node {
 struct task_struct * info;
 struct stack_node * next;
};

struct stack_node * headptr = 0;

void stack_push(struct task_struct *  a){
    struct stack_node * tmp = (struct stack_node *)vmalloc(sizeof(struct stack_node));
    tmp->info = a;
    tmp->next = headptr;
    headptr = tmp;
}

struct task_struct *  stack_pop(void){
    struct stack_node * tmp = headptr;
    headptr = headptr->next;
    struct task_struct *  ret = tmp->info;
    vfree(tmp);
    return ret;
}

int stack_empty(void){
 return headptr == 0;   
}
/////////////////////////////

SYSCALL_DEFINE1(processdfs, int, arg_pid){
 struct task_struct * task;
 for_each_process(task){
     if(task->pid == arg_pid){
      //dfs obilazak stabla
         stack_push(task);
         while(!stack_empty()){
            struct task_struct * tmp = stack_pop();
            printk("%d  %s", tmp->pid, tmp->comm);
            struct task_struct * tmp1;
            struct list_head * list;
            list_for_each(list, &tmp->children){
                tmp1 = list_entry(list, struct task_struct, sibling);
                stack_push(tmp1);
            }
         }
         break;
     }
 }
 return 0;
}
