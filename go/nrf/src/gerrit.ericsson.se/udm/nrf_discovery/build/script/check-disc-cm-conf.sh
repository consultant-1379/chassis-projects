#!/bin/bash

function CmProviderIsReady() {
    isSuccess=`curl -i -s -X GET -H "Content-Type:application/json" ${CM_SERVICE_NAME}/configurations/${CM_SCHEMA_PREFIXNAME}|grep "200 OK"|wc -l`
    if [ $isSuccess -le 0 ]; then
      echo "nrf configuration(${CM_SCHEMA_PREFIXNAME}) can't be got from CM, waiting..."
      return 1
    else
      echo "nrf configuration(${CM_SCHEMA_PREFIXNAME}) is got from CM successfully"
      return 0
    fi
}

##########start#########

cm_check_enabled=$(echo ${CM_CHECK_ENABLED} | tr '[A-Z]' '[a-z]')

echo "cm_hook_enabled:$cm_check_enabled  cm_service_name:${CM_SERVICE_NAME}  cm-schema-prefixname:${CM_SCHEMA_PREFIXNAME}"

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
