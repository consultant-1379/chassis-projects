#!/bin/bash

nrfagent_conf="ericsson-nrf-agent"
function CmProviderIsReady() {
    isSuccess=`curl -i -s -X GET -H "Content-Type:application/json" ${CM_URI_PREFIX}configurations/${nrfagent_conf}|grep "200 OK"|wc -l`
    if [ $isSuccess -le 0 ]; then
      echo "nrfagent configuration(${nrfagent_conf}) can't be got from CM, waiting..."
      return 1
    else
      echo "nrfagent configuration(${nrfagent_conf}) is got from CM successfully"
      return 0
    fi
}

##########start#########

cm_check_enabled=$(echo ${CM_CHECK_ENABLED} | tr '[A-Z]' '[a-z]')

echo "cm_check_enabled:$cm_check_enabled  cm_service_name:${CM_URI_PREFIX}configurations/${nrfagent_conf}  cm-conf-name:${nrfagent_conf}"

if [ "${cm_check_enabled}" == "false" ];then
      echo " skip cm configuration check"
      exit 0
fi

isContinue=true
while $isContinue
do
	CmProviderIsReady
    if [ $? != "0" ]; then
        sleep 1
    else
          exit 0
    fi
done
