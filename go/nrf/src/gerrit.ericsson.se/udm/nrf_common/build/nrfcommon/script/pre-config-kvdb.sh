#!/bin/bash

function IsJarExit() {
    listDeployedCommand=`curl -s -XPOST -H "Content-Type: application/json" -d'{"command":"list deployed"}' ${KVDB_ADMIN_MGR}/gfsh-commands`
    echo "list deployed command id = [$listDeployedCommand]"
    listDeployedCommandId=`echo ${listDeployedCommand} | jq -r '."commandId"'`
    sleep 2
    jars=`curl -s ${KVDB_ADMIN_MGR}/gfsh-commands/${listDeployedCommandId}`
    isJarExit=`echo ${jars} | jq -r '."output"' | grep "kvdb-1.0.jar"`
    if [[ "$isJarExit" != "" ]]; then
        return 1
    else
        return 0
    fi
}

function IsCacheRegionExit() {
    listRegionCommand=`curl -s -XPOST -H "Content-Type: application/json" -d'{"command":"list regions"}' ${KVDB_ADMIN_MGR}/gfsh-commands`
    echo "list region command id = [$listRegionCommand]"
    listRegionCommandId=`echo ${listRegionCommand} | jq -r '."commandId"'`
    sleep 2
    regions=`curl -s ${KVDB_ADMIN_MGR}/gfsh-commands/${listRegionCommandId}`
    isCacheRegionExit=`echo ${regions} | jq -r '."output"' | grep ericsson-nrf-cachenfprofiles`
    if [[ "$isCacheRegionExit" != "" ]]; then
        return 1
    else
        return 0
    fi
}

function IsDistributedRegionExist(){
    listRegionCommand=`curl -s -XPOST -H "Content-Type: application/json" -d'{"command":"list regions"}' ${KVDB_ADMIN_MGR}/gfsh-commands`
    echo "list region command id = [$listRegionCommand]"
    listRegionCommandId=`echo ${listRegionCommand} | jq -r '."commandId"'`
    sleep 2
    regions=`curl -s ${KVDB_ADMIN_MGR}/gfsh-commands/${listRegionCommandId}`
    isCacheRegionExit=`echo ${regions} | jq -r '."output"' | grep ericsson-nrf-distributedlock`
    if [[ "$isCacheRegionExit" != "" ]]; then
        return 1
    else
        return 0
    fi
}

function AlterDisbutedRegion(){
    createRegionResult=`curl -s -XPOST -H "Content-Type: application/json" -d'{"command":"alter region --name=/ericsson-nrf-distributedlock --entry-time-to-live-expiration-action=DESTROY --entry-time-to-live-custom-expiry=com.ericsson.geode.expiry.MyCustomExpiry"}' ${KVDB_ADMIN_MGR}/gfsh-commands`
    echo "create region command id = [$createRegionResult]"
    createRegionCommandId=`echo ${createRegionResult} | jq -r '."commandId"'`
    sleep 2
    result=`curl -s ${KVDB_ADMIN_MGR}/gfsh-commands/${createRegionCommandId}`
    echo "create distributedlock region result = [$result]"
    IsDistributedRegionExist
    if [[ $? == "1" ]]; then
        echo "create distributedlock region sussess"
        return 1
    else
        echo "create distributedlock region fail"
        return 0
    fi
}

function DeployJar(){
    jarFile="kvdb-1.0.jar"
    echo "curl -s -H \"Content-Type: multipart/form-data\"  -F\"file=@/bin/$jarFile\"  -XPUT ${KVDB_ADMIN_MGR}/app-jars/${jarFile}"
    postJarResult=`curl -s -H "Content-Type: multipart/form-data"  -F"file=@/bin/$jarFile"  -XPUT ${KVDB_ADMIN_MGR}/app-jars/${jarFile}`
    echo "post jar to kvdb admin mgr, msg=[$postJarResult]"
    deployJarCommand=`curl -s -XPOST -H "Content-Type: application/json" -d'{"command":"deploy --jar=/opt/dbservice/data/app-jars/'${jarFile}'"}' ${KVDB_ADMIN_MGR}/gfsh-commands`
    echo "deploy jar command id = [$deployJarCommand]"
    deployJarCommandId=`echo ${deployJarCommand} | jq -r '."commandId"'`
    sleep 2
    deployResult=`curl -s ${KVDB_ADMIN_MGR}/gfsh-commands/${deployJarCommandId}`
    echo "deploy jar result = [$deployResult]"
}

function DeployCacheRegion() {
    IsCacheRegionExit
    if [[ $? == "1" ]]; then
        return 1
    else
        createRegionResult=`curl -s -XPOST -H "Content-Type: application/json" -d'{"command":"create region --name=/ericsson-nrf-cachenfprofiles --type=REPLICATE --enable-statistics=true --eviction-entry-count=20000 --eviction-action=local-destroy"}' ${KVDB_ADMIN_MGR}/gfsh-commands`
        echo "create region command id = [$createRegionResult]"
        createRegionCommandId=`echo ${createRegionResult} | jq -r '."commandId"'`
        sleep 2
        result=`curl -s ${KVDB_ADMIN_MGR}/gfsh-commands/${createRegionCommandId}`
        echo "create cache region result = [$result]"
        IsCacheRegionExit
        if [[ $? == "1" ]]; then
            echo "create cache region sussess"
            return 1
        else
            echo "create cache region fail"
            return 0
        fi
    fi
}

##########start#########

##########1. configure pdx disk store, workaround to fix unknown pdx type#########
echo "start configure pdx disk store"
while true
do
    configurePdxResult=`curl -s -XPOST -H "Content-Type: application/json" -d'{"command":"configure pdx --disk-store=DEFAULT --read-serialized=true"}' ${KVDB_ADMIN_MGR}/gfsh-commands`
    echo "configure pdx command id = [$configurePdxResult]"
    configurePdxCommandId=`echo ${configurePdxResult} | jq -r '."commandId"'`
    if [[ ${configurePdxCommandId} != null && "${configurePdxCommandId}" != "" ]]; then
        sleep 2
        pdxResult=`curl -s ${KVDB_ADMIN_MGR}/gfsh-commands/${configurePdxCommandId}`
        echo "configure pdx result = [$pdxResult]"
        break
    fi
done


##########2. configure cache region#########
cache_region_enabled=$(echo ${CACHE_REGION_ENABLED} | tr '[A-Z]' '[a-z]')

echo "cache_region_enabled:$cache_region_enabled"

if [[ "${cache_region_enabled}" == "false" ]];then
      echo " skip create cache region"
      exit 0
fi

IsCacheRegionExit
if [[ $? == "1" ]]; then
    echo "cache region already exit"
    exit 0
fi

isContinue=true
while ${isContinue}
do
    IsJarExit
    if [[ $? != "1" ]]; then    
        DeployJar
        sleep 2
    else
        break
    fi 
done

while ${isContinue}
do
    DeployCacheRegion
    if [[ $? != "1" ]]; then
            sleep 5
    else
            break
    fi

done

exit 0
