#!/bin/bash
cm_schema_directory="/cm-schema"
cm_init_enabled=$(echo ${CM_INIT_ENABLED} | tr '[A-Z]' '[a-z]')

echo "DEBUG: pre_hook_enabled:$cm_init_enabled  cm_service_name:${CM_SERVICE_NAME}  cm-schema-name:${CM_SCHEMA_NAME}"

if [ "${cm_init_enabled}" == "false" ];then
    echo "INFO: the hook is disabled, skip cm schema and config importing"
    exit 0
fi

while true; 
do 
    isAvaialable=`curl  -i -s -X GET -H "Content-Type:application/json" ${CM_SERVICE_NAME}/schemas|grep "200 OK"|wc -l` 
    if [ $isAvaialable -le 0 ]; then
        echo "WARNING: CM is not available! Waiting..."
    else
	break
    fi  
    sleep 5   
done

# import cm schema
echo "DEBUG: start to import cm schema ${CM_SCHEMA_NAME}"
retCode=`curl -i -s -X POST -H "Content-Type: multipart/form-data" -F name="${CM_SCHEMA_NAME}"  -F title="ericsson nrf CM schema"  -F file=@${cm_schema_directory}/${CM_SCHEMA_NAME}.json -F yangArchive=@${cm_schema_directory}/${CM_SCHEMA_NAME}.yang.tar.gz  ${CM_SERVICE_NAME}/schemas|grep -E "201 CREATED|409 CONFLICT"` 

echo "retCode: $retCode"
if [[ "${retCode}" =~ "201" ]]; then
    echo "DEBUG: import cm schema ${CM_SCHEMA_NAME} successfully"
elif [[ "${retCode}" =~ "409" ]];then
    echo "DEBUG: cm schema ${CM_SCHEMA_NAME} already existed"		
else
    echo "ERROR: import cm schema ${CM_SCHEMA_NAME} failure!"
    exit 1
fi

exit 0
