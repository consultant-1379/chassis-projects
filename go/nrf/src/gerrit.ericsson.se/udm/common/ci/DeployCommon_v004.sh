#!/bin/sh

helmName=$1
port=$2

helm init --client-only
helm repo add proj-5gc-udm-helm https://armdocker.rnd.ericsson.se/artifactory/proj-5gc-udm-helm
helm repo update

helm delete --purge --timeout=0 $helmName
sleep 10

kubectl delete pods --all --grace-period=0 --force -n $helmName
sleep 10

helm install proj-5gc-udm-helm/eric-adp-5g-udm --name $helmName --namespace $helmName --set eric-data-coordinator-zk.cpu=200m --set eric-data-message-bus-kf.cpu=200m --set eric-data-coordinator-zk.memory=512Mi  --set eric-data-message-bus-kf.memory=512Mi --set tags.eric-data-document-database-pg=true --set tags.eric-cm-mediator=true --set tags.eric-pm-server=true --set tags.eric-data-coordinator-zk=true --set tags.eric-data-message-bus-kf=true --version=0.0.4 --set eric-cm-mediator.service.nodePort=$port
