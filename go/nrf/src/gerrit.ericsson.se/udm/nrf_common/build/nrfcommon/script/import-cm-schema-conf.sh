#!/bin/bash
cm_schema_directory="/cm-schema"
cm_init_config_directory="/cm-init-config"
#schema_list=`ls ${cm_schema_directory}`
init_config_list=`ls ${cm_init_config_directory}`

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

#for schema in ${schema_list}
#do
#    if [[ "${schema}" =~ ".json" ]]; then
#	    echo "schema: ${cm_schema_directory}/${schema}  yang: ${cm_schema_directory}/${schema%.json*}.yang.tar.gz"
#        retCode=`curl -i -s -X POST -H "Content-Type: multipart/form-data" -F name="${CM_SCHEMA_PREFIXNAME}-nrf-cm"  -F title="ericsson nrf CM schema"  -F file=@${cm_schema_directory}/${schema} -F yangArchive=@${cm_schema_directory}/${schema%.json*}.yang.tar.gz  ${CM_SERVICE_NAME}/schemas|grep -E "201 CREATED|409 CONFLICT"` 
# 	    echo "retCode: $retCode"  
#        if [[ "${retCode}" =~ "201" ]]; then
#            echo "create cm schema ${schema} successfully"
#	    elif [[ "${retCode}" =~ "409" ]];then
#            echo "cm schema have existed"
#            exit 0		
#	    else
#            echo "create cm schema failure"
#            exit 1
#        fi
#	fi
#done

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

# import cm init config
for config in ${init_config_list}
do
    echo "DEBUG: start to import cm init config ${CM_SCHEMA_NAME}"
    sed -i "s/ericCmTempName0011/${CM_SCHEMA_NAME}/" ${cm_init_config_directory}/${config}
    uuid=`uuidgen -t`
    sed -i "s/0c765084-9cc5-49c6-9876-ae2f5fa2a63f/${uuid}/" ${cm_init_config_directory}/${config}

    configRetCode=`curl -i -s -X POST --data "$(cat ${cm_init_config_directory}/${config})" -H "Content-Type:application/json" ${CM_SERVICE_NAME}/configurations|grep -E "201 CREATED|409 CONFLICT"`

    echo "retCode: $configRetCode"
    if [[ "${configRetCode}" =~ "201" ]]; then
        echo "DEBUG: import cm init config ${CM_SCHEMA_NAME} successfully"
        exit 0
    elif [[ "${configRetCode}" =~ "409" ]];then
        echo "DEBUG: cm init config ${CM_SCHEMA_NAME} already existed"
        exit 0
    else
        echo "ERROR: import cm init config ${CM_SCHEMA_NAME} failure!"
        exit 1
    fi
done

exit 0
