#!/bin/bash
cm_schema_directory="/etc/nrfagent/cm/schemas"
cm_init_config_directory="/etc/nrfagent/cm/configurations"
cm_yang_directory="/etc/nrfagent/cm/yang"
schema_list=`ls ${cm_schema_directory}`
init_config_list=`ls ${cm_init_config_directory}`

pre_hook_enabled=$(echo ${PRE_HOOK_ENABLED} | tr '[A-Z]' '[a-z]')

echo "pre_hook_enabled:${pre_hook_enabled}  cm_uri_prefix:${CM_URI_PREFIX}  cm_confname_prefix:${NF_INSTANCE_NAME}"

if [ "${pre_hook_enabled}" == "false" ];then
      echo " skip cm schema and config importing"
      exit 0
fi

for schema in ${schema_list}
do
    retCode="201"
    if [ ! -f "${cm_yang_directory}/${schema%.json*}.yang.tar.gz" ]; then
        echo "schema: ${cm_schema_directory}/${schema}"
        sed -i "s/ericmprefix0011/${NF_INSTANCE_NAME}/g" ${cm_schema_directory}/${schema}
        retCode=`curl -i -s -X POST --data "$(cat ${cm_schema_directory}/${schema})" -H "Content-Type:application/json" ${CM_URI_PREFIX}schemas|grep -E "201 CREATED|409 CONFLICT"`
    else
        echo "schema: ${cm_schema_directory}/${schema}  yang: ${cm_yang_directory}/${schema%.json*}.yang.tar.gz"
        retCode=`curl -i -s -X POST -H "Content-Type: multipart/form-data" -F name="${schema%.json*}" -F title="ericsson nrf agent cm schema" -F file=@${cm_schema_directory}/${schema} -F yangArchive=@${cm_yang_directory}/${schema%.json*}.yang.tar.gz  ${CM_URI_PREFIX}schemas|grep -E "201 CREATED|409 CONFLICT"`
    fi
    
    if [[ "${retCode}" =~ "201" ]]; then
        echo "create cm schema ${schema} successfully"
    elif [[ "${retCode}" =~ "409" ]];then
        echo "cm schema ${schema} have existed"
    else
        echo "create cm schema ${schema} failure, ${retCode}"
        exit 1
    fi
done

for config in ${init_config_list}
do
    sed -i "s/ericmprefix0011/${NF_INSTANCE_NAME}/g" ${cm_init_config_directory}/${config}

    retCode=`curl  -i -s -X POST --data "$(cat ${cm_init_config_directory}/${config})" -H "Content-Type:application/json" ${CM_URI_PREFIX}configurations|grep -E "201 CREATED|409 CONFLICT"`
    if [[ "${retCode}" =~ "201" ]]; then
        echo "create cm configuration ${config} successfully"
    elif [[ "${retCode}" =~ "409" ]];then
        echo "cm configuration ${config} have existed"
    else
        echo "create cm configuration ${config} failure, ${retCode}"
        exit 1
    fi
done

echo " create cm schema and init config successfully"
exit 0