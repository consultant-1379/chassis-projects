#!/bin/sh
set -x 
NRFVERSION=$1
NRFREPO=$2
NRFCHART=$3

export ENCRYPTED=$(echo "admin" | openssl passwd -1 -stdin)
TESTREPODIR="/git/chassis-projects/integration/eric_nrf"
source ${TESTREPODIR}/common.sh

if [ -d /secrets/git/ ] ; then
  USERNAME=$(cat /secrets/git/user)
  PASSWORD=$(cat /secrets/git/pwd)
fi
helm init --upgrade
helm init --client-only
kubectl get secret rbd-client-secret -n default --export -o yaml | kubectl apply -n $NAMESPACE -f -

kubectl  create secret docker-registry ipworks --docker-username=${USERNAME} --docker-password=${PASSWORD} --docker-server=armdocker.rnd.ericsson.se --namespace $NAMESPACE 

helm repo add nrf  ${NRFREPO}   --username  ${USERNAME}  --password ${PASSWORD}
helm repo add arm-eric-adp-5g-udm https://arm.lmera.ericsson.se/artifactory/proj-5g-udm-release-helm  --username  ${USERNAME}  --password ${PASSWORD}  
helm repo update

cd ${TESTREPODIR}/Tool/certs
kubectl create secret generic dsa-tls-secret --from-file=ca.crt=2_intermediate/certs/ca-chain.cert.pem --from-file=tls.crt=5_application_ecdsa/certs/eccd-udm00027.seli.gic.ericsson.se.cert.pem --from-file=tls.key=5_application_ecdsa/private/eccd-udm00027.seli.gic.ericsson.se.key.pem -n istio-system 
kubectl create -n $NAMESPACE secret generic eric-ccrc-sbi-client-ca-certs --from-file=ca-chain.cert.pem=2_intermediate/certs/ca-chain.cert.pem 
kubectl create -n $NAMESPACE secret tls eric-ccrc-sbi-client-certs --cert 4_client/certs/eccd-udm00027.seli.gic.ericsson.se.cert.pem --key 4_client/private/eccd-udm00027.seli.gic.ericsson.se.key.pem 
     
helm install arm-eric-adp-5g-udm/eric-adp-5g-udm --name="$HELMNAME-adp" -f ${TESTREPODIR}/values_adp_minikube.yaml --version $VERSION --set global.registry.pullSecret=ipworks --namespace $NAMESPACE  \
      --set eric-cm-yang-provider.users.user=admin \
      --set eric-cm-yang-provider.users.groups='[system-admin\, system-security-admin]' \
      --set eric-pm-server.rbac.appMonitoring.enabled=false \
      --set eric-pm-bulk-reporter.users.user=admin \
      --set eric-pm-bulk-reporter.users.groups='[pm-sftp-users]' \
      --set eric-pm-server.rbac.clusterMonitoring.enabled=false \
      --set eric-cm-yang-provider.users.encryptedPass=${ENCRYPTED} \
      --set eric-pm-bulk-reporter.users.encryptedPass=${ENCRYPTED} \
      --set tags.eric-log-transformer=false \
      --set tags.eric-log-shipper=false \
      --set eric-pm-bulk-reporter.env.nodeName=seliius22081,eric-pm-bulk-reporter.env.nodeType=5G_UDM \
      --set eric-cm-yang-provider.persistentVolumeClaim.storageClassName=$STORAGECLASS \
      --set tags.eric-data-search-engine=false \
      --set tags.eric-data-visualizer-gf=false \
      --set tags.eric-data-visualizer-kb=false \
      --set tags.eric-log-shipper=false \
      --set eric-cm-yang-provider.persistentVolumeClaim.storageClassName=$STORAGECLASS \
      --set eric-data-document-database-pg.persistence.storageClass=$STORAGECLASS \
      --set eric-data-coordinator-zk.persistantVolumeClaim.storageClassName=$STORAGECLASS \
      --set eric-data-message-bus-kf.persistentVolumeClaim.storageClassName=$STORAGECLASS \
      --set eric-pm-bulk-reporter.persistentVolumeClaim.storageClassName=$STORAGECLASS

if [ $? != 0 ] ; then
  echo "helm install adp failed!!!"
  exit 1
fi

 time=0
while [ $time -le 60 ] ; do
  kubectl get pod --namespace $NAMESPACE|grep -v Completed |while read line ; do
  if [ "$line" == "" ] ; then
    break
  else
    echo $line >> /tmp/test123
    fi
  done
  real=$(sed '/NAME/d'  /tmp/test123| sed '/^$/'d |grep -v NAME |awk '{print $2}' |awk -F "/" '{print $1}')
  expected=$(sed '/NAME/d'  /tmp/test123| sed '/^$/'d |grep -v NAME |awk '{print $2}'|awk -F "/" '{print $2}')
  if [ "$real" == "$expected" ] ; then
    echo
    echo "Pods of ADP is ready"
    cat /tmp/test123
    rm -rf /tmp/test123
    break
  else
    echo "Pods of ADP is not ready"
    cat /tmp/test123
    rm -rf /tmp/test123
  fi
  time=`expr $time + 1`
  sleep 10
done



cmIP=$(kubectl get svc -n ${NRFCOMMON}|grep eric-cm-mediator|grep cm-mediator |awk -F " " '{print $3}')
cmPort=$(kubectl get svc -n ${NRFCOMMON}|grep eric-cm-mediator|grep cm-mediator |awk -F " " '{print $5}'| cut -d / -f 1| cut -d : -f 1)



h2IP=`grep 192.168 /etc/hosts|awk -F ' ' '{print $1}'`



kubectl apply -n $NAMESPACE -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: sim-server
spec:
  hosts:
  - $h2IP
  ports:
  - number: 5443
    name: http-port-for-tls-origination
    protocol: HTTP2
---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: sim-server
spec:
  host: $h2IP
  trafficPolicy:
    loadBalancer:
      simple: ROUND_ROBIN
    portLevelSettings:
    - port:
        number: 5443
      tls:
        mode: MUTUAL
        clientCertificate: /etc/istio/egress-certs/tls.crt
        privateKey: /etc/istio/egress-certs/tls.key
        caCertificates: /etc/istio/egress-ca-cert/ca-chain.cert.pem
EOF



helm install --name="${HELMNAME}-release" nrf/${NRFCHART} --version=$NRFVERSION --set eric-nrf-common.rbac.istio_clusterrolebinding=nrf-istio-cb-${NAMESPACE} --set tags.eric-nrf-all=true --set eric-nrf-common.global.ingressSecret=istio-ingress-certs-${HELMNAME} --set global.adpNamespace=${NRFCOMMON} --values=${TESTREPODIR}/helm/eric-nrf/values/values_minikube.yaml --set eric-nrf-common.eric-data-kvdb-ag.persistence.data.storageClass=$STORAGECLASS --set global.registry.pullSecret=ipworks --namespace $NAMESPACE  || exit 1


if [ $? != 0 ] ; then
  echo "helm install nrf failed!!!"
  exit 1
fi

time=0
while [ $time -le 60 ] ; do
  status=$(kubectl get pod --namespace $NAMESPACE|grep eric-nrf-common-cm-init|awk '{print $3}')
   if [ "$status" == "Completed" ] ; then
    break
   fi
done
ingressIP=$(kubectl get nodes  -o jsonpath="{.items[0].status.addresses[0].address}")
ingressPort=$(kubectl get svc -n istio-system|grep nrf-traffic |awk '{print $5}'|awk -F "/" '{print $1}'|awk -F ":" '{print $2}'|head -1)


sed '2i "name": "ericsson-nrf",' ${TESTREPODIR}/cm_init_config.json > ${TESTREPODIR}/cm_init_config_post.json
echo "curl -i -X POST --data "$(cat ${TESTREPODIR}/cm_init_config_post.json)" -H "Content-Type:application/json" http://${cmIP}:${cmPort}/cm/api/v1.1/configurations"
curl -i -X POST --data "$(cat ${TESTREPODIR}/cm_init_config_post.json)" -H "Content-Type:application/json" http://${cmIP}:${cmPort}/cm/api/v1.1/configurations

