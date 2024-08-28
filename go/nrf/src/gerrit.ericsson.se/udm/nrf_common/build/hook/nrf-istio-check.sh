#!/bin/sh

export TOKEN=/var/run/secrets/kubernetes.io/serviceaccount/token
export CERT=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
 
isExist=`curl -s --cacert $CERT -H "Authorization: Bearer $(cat $TOKEN)" -G https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_PORT_443_TCP_PORT/api/v1/namespaces/istio-system/services/istio-ingress | grep "selfLink" |wc -l`      


if [ $isExist -gt 0 ]; then
  echo " istio is running"
  exit 0

else
  echo "istio is not running"
  exit 1

fi
