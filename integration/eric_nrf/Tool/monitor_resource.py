#!/usr/bin/python
import xlwt
import xlrd
import subprocess
import os
import threading
#import xlutils;
from subprocess import Popen, PIPE

from xlutils.copy import copy;
kafka_nodelist=[]
disc_nodelist=[]
colums =0


#usage of scprit:
#1. execute command ssh-keygen on client which will run this scprit.
#2. cat ~/.ssh/id_rsa.pub
#3.1 copy above content in file ~/.ssh/authorized_keys on every worker node ( for eccd env).
#3.2 execute ssh-copy-id eccd@node( for udm006,043)
#4. add the name and ip of every worker node in /etc/hosts on this client.
#5. export KUBECONFOG=admin.conf(test env)
#6. modify the define of process_name based on your requirement.
#eg:process_name = {'pod1 keyword': ['container name1','container name2'], 'pod2 keyword': ['container name']}
#7. you can start use this scprit now.

process_name = {'discovery': ['eric-nrf-discovery','eric-nrf-disc-dbproxy'], 'kvdb-ag-server': ['eric-nrf-kvdb-ag-server:']}
#process_name = {'discovery': ['nrfdisc','dbproxy.DBProxyServer'], 'kvdb-ag-server': ['eric-nrf-kvdb-ag-server'], 'message-bus-kf': ['SupportedKafka']}
#process_name = {'discovery': ['nrfdisc','dbproxy.DBProxyServer']}
#pid = {'discovery': {}, 'kvdb-ag-server': {}, 'message-bus-kf': {}}
pid = {}
cpu_memory = {}
execel_colum = {}
thread_list = []

#kafka_node=subprocess.Popen('kubectl get pod -o wide | egrep \"message-bus-kf\" | awk \'{print $7}\'', shell = True, stdout = subprocess.PIPE, stderr = subprocess.STDOUT)
def get_node():

#    for k in process_name:
#        print(k)
#        pid[k] = {}
#        disc_node = subprocess.Popen('kubectl get pod -o wide | egrep ' + k + ' | awk \'{print $1,$7}\'', shell = True, stdout = subprocess.PIPE, stderr = subprocess.STDOUT)
#        for line in disc_node.stdout.readlines() :
#            pod, node = str(line,'utf-8').strip('\r\n').split(' ',1)
#            for process in process_name[k]:
#                pid_temp = subprocess.Popen('ssh eccd@' + node + ' \"ps -ef|grep ' + process + '|grep -v grep\" | awk -F\' \' \'{print $2}\'',shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
#                pid_str = str(pid_temp.stdout.readlines()[-1],'utf-8').strip('\r\n')
#                if pod not in pid[k]:
#                    pid[k][pod] = []
#                pid[k][pod].append([process, pid_str, node])
    #print(pid)


    for k in process_name:
        print(k)
        pid[k] = {}
        disc_node = subprocess.Popen('kubectl get pod -o wide | egrep ' + k + ' | awk \'{print $1,$7}\'', shell = True, stdout = subprocess.PIPE, stderr = subprocess.STDOUT)
        for line in disc_node.stdout.readlines() :
            pod, node = str(line,'utf-8').strip('\r\n').split(' ',1)
            for process in process_name[k]:
                container_id_temp = subprocess.Popen('kubectl get pod ' + pod + ' -o jsonpath=\'{range .status.containerStatuses[*]}{.name}{":\\t"}{.containerID}{"\\n"}{end}\' | grep ' + process + ' |cut -f3 -d "/"', shell = True, stdout = subprocess.PIPE, stderr = subprocess.STDOUT)
                container_id = str(container_id_temp.stdout.readlines()[-1],'utf-8')
                pid_temp = subprocess.Popen('ssh eccd@' + node + ' sudo docker inspect --format \'{{.State.Pid}}\' ' + container_id, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
                pid_str = str(pid_temp.stdout.readlines()[-1],'utf-8').strip('\r\n')
                if pod not in pid[k]:
                    pid[k][pod] = []
                pid[k][pod].append([process, pid_str, node])



def get_Remote_CPU_Memory(pod, process, cpu_memory):

     processid = process[1]
     node = process[2]
     cpu_temp = subprocess.Popen('ssh eccd@' + node + ' \"top -n 2 -d 0.5 -b -p ' + processid + '|grep ' + processid + ' | tail -1 \" | awk -F\' \' \'{print $9}\'',shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
     cpu = str(cpu_temp.stdout.readlines()[-1],'utf-8').strip('\r\n')
     memory_temp = subprocess.Popen('ssh eccd@' + node + ' \"top -n 2 -d 0.5 -b -p ' + processid + '|grep ' + processid + ' | tail -1 \" | awk -F\' \' \'{print $6}\'',shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
     memory = str(memory_temp.stdout.readlines()[-1],'utf-8').strip('\r\n')
     #print(pod,process[0],cpu,memory)
     pod_process = pod + '_' + process[0]
     cpu_memory[pod_process] = [cpu, memory]


def get_CPU_Memory():

    global colums
    cpu_memory.clear()
    thread_list = [] 
    for k in sorted(pid):
        for pod in sorted(pid[k]):
            for process in pid[k][pod]:
                t = threading.Thread( None, get_Remote_CPU_Memory, None, (pod, process, cpu_memory,) )
                thread_list.append(t)
    for i in range(len(thread_list)):
        thread_list[i].start()
    for i in range(len(thread_list)):
        thread_list[i].join()
    for pod_process in cpu_memory: 
        if pod_process not in execel_colum:
            execel_colum[pod_process] = colums
            colums = colums + 2

    print(444444444444444444444444444)
    for k in sorted(cpu_memory):
        print(k,cpu_memory[k])
    print(555555555555555555555555555)


if __name__ == '__main__':
    get_node()

# ################################record data to excel start ####################
gConst = {'xls':{'sheetName':'result', 'fileName':'cpu_memory_monitor.xls' }}
styleBoldRed = xlwt.easyxf('font: color-index red, bold on')
headerStyle = styleBoldRed
wb = xlwt.Workbook()
ws = wb.add_sheet(gConst['xls']['sheetName'])

wb.save(gConst['xls']['fileName'])
# #
# # # open existed xls file

oldWb = xlrd.open_workbook(gConst['xls']['fileName'], formatting_info=True)
print(oldWb)  # <xlrd.book.Book object at 0x000000000315C940>
newWb = copy(oldWb)
print(newWb)
newWs = newWb.get_sheet(0)
for i in range(1,180):
    get_CPU_Memory()
    #print(cpu_memory)
    for pod_process in sorted(cpu_memory):
        #print(pod_process)
        newWs.write(i, execel_colum[pod_process], cpu_memory[pod_process][0])
        newWs.write(i, execel_colum[pod_process]+1, cpu_memory[pod_process][1])
for pod_process in sorted(cpu_memory):
    #print(pod_process)
    newWs.write(0, execel_colum[pod_process], pod_process + "_CPU")
    newWs.write(0, execel_colum[pod_process]+1, pod_process + "_Memory")
newWb.save(gConst['xls']['fileName'])
print("############## end ###################")

# ################################record data to excel finish ####################
