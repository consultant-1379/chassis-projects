#!/bin/bash

echo " starting post deleting....."

post_delete_enabled=$(echo ${POST_HOOK_DELETE_ENABLED} | tr '[A-Z]' '[a-z]')

echo "pre_hook_enabled:$post_delete_enabled  cm_service:${CM_SERVICE}  cm-schema-name:${CM_SCHEMA_NAME}"

if [ "${post_delete_enabled}" == "false" ];then
      echo " skip cm schema and config deleting"
      exit 0
fi


retCode=`curl -i -s -X DELETE   -H "Content-Type:application/json"  ${CM_SERVICE}/configurations/${CM_SCHEMA_NAME}|grep -E "200 OK|404 NOT FOUND"` 
    echo "retCode: $retCode"  
if [[ "${retCode}" =~ "200" ]]; then
    echo "Deleted cm configuration ${CM_SCHEMA_NAME} successfully"
elif [[ "${retCode}" =~ "404" ]];then
    echo "Cannot find cm configuration ${CM_SCHEMA_NAME}"
   exit 0		
else
    echo "delete cm configuration failure"
#    exit 1
fi

retCode=`curl -i -s -X DELETE   -H "Content-Type:application/json"  ${CM_SERVICE}/schemas/${CM_SCHEMA_NAME}|grep -E "200 OK|404 NOT FOUND"` 
    echo "retCode: $retCode"  
if [[ "${retCode}" =~ "200" ]]; then
    echo "Deleted cm schema ${CM_SCHEMA_NAME} successfully"
elif [[ "${retCode}" =~ "404" ]];then
    echo "Cannot find cm schema ${CM_SCHEMA_NAME}"
   exit 0		
else
    echo "delete cm schema failure"
#    exit 1
fi

exit 0
