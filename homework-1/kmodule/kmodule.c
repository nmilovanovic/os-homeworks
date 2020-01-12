#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/sched.h>
#include <linux/list.h>
//https://stackoverflow.com/questions/32189782/using-for-each-process-in-my-program-prevents-me-from-compiling-compiler-says
#include <linux/sched/signal.h>

//////////////////////////////
struct node {
 struct task_struct * info;
 struct node * next;
};

struct node * headptr = 0;

void stack_push(struct task_struct *  a){
    struct node * tmp = (struct node *)vmalloc(sizeof(struct node));
    tmp->info = a;
    tmp->next = headptr;
    headptr = tmp;
}

struct task_struct *  stack_pop(void){
    struct node * tmp = headptr;
    headptr = headptr->next;
    struct task_struct *  ret = tmp->info;
    vfree(tmp);
    return ret;
}

int stack_empty(void){
 return headptr == 0;   
}
/////////////////////////////

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Nikola Milovanovic");
MODULE_DESCRIPTION("Linux kernel module for listing process DFS tree");
MODULE_VERSION("0.05");

static int arg_pid = 1;
module_param(arg_pid, int, 0);

static int __init kmodule_enter(void){
 printk(KERN_INFO "Module is loaded %d\n", arg_pid);
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

static void __exit kmodule_exit(void){
 printk(KERN_INFO "Module is unloaded %d\n", arg_pid);
}

module_init(kmodule_enter);
module_exit(kmodule_exit);
