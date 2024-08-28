#!/bin/bash

MODE="-server -XX:MetaspaceSize=64M"
HEAP_PARAMETERS="-Xms500m -Xmx500m"
GC_PARAMETERS="-XX:NewRatio=1 -XX:MaxTenuringThreshold=15 -XX:MaxGCPauseMillis=10 -XX:+UseFastAccessorMethods"
PRINT_GC_INFO="-XX:+PrintGCDetails -XX:+PrintHeapAtGC -XX:+PrintGCDateStamps -XX:+PrintGCTimeStamps -Xloggc:/home/dbproxy/gc_log"
SCRIPT_PARAMETERS=$@

GB_IN_BYTES=1073741824
if [ -n "${MEMORY_LIMIT}" ];then
  if [[ ${MEMORY_LIMIT} -gt ${GB_IN_BYTES} ]];then
    memory_limit=$((${MEMORY_LIMIT} / ${GB_IN_BYTES}))
    if [[ ${memory_limit} -gt 4 ]];then
        memory_limit=`expr $memory_limit - 2`
    else
        memory_limit=`expr $memory_limit - 1`
    fi
	if [[ ${memory_limit} -eq 0 ]];then 
	   memory_limit=1
	fi
    HEAP_PARAMETERS="-Xms${memory_limit}g -Xmx${memory_limit}g"
  fi
fi

if [ -n "${CPU_LIMIT}" ];then
  if [[ ${CPU_LIMIT} -gt 1 ]]; then
    GC_PARAMETERS="-XX:+UseParallelOldGC ${GC_PARAMETERS}"
  else
    GC_PARAMETERS="-XX:+UseSerialGC ${GC_PARAMETERS}"
  fi
fi

if [[ ${GC_DEBUG} = "true" ]];then
  JVM_PARAMETERS="${MODE} ${HEAP_PARAMETERS} ${GC_PARAMETERS} ${PRINT_GC_INFO} ${SCRIPT_PARAMETERS}"
else
  JVM_PARAMETERS="${MODE} ${HEAP_PARAMETERS} ${GC_PARAMETERS} ${SCRIPT_PARAMETERS}"
fi

echo "${JVM_PARAMETERS}"
exec java ${JVM_PARAMETERS}
