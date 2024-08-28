#!/bin/sh

ZOOKEEPER=$2
BROKER_LIST=$3
TOPIC_NAME=$4
TOPIC_PATITIONS=$5
TOPIC_REPLICATION=$6

# list topics in Kafka
topiclist=`/usr/bin/kafka-topics --zookeeper ${ZOOKEEPER} --list 2>/dev/null`
if [ $? -ne 0 ]; then
        echo "pre-install.sh check Kafka serivce failed"
        exit 1
else
	# check topic exist or not
	if echo ${topiclist} | grep -w ${TOPIC_NAME} > /dev/null; then
		# found the topic, then describe topic check partition and replication
		echo "pre-install.sh found topic in kafka"
	else
		/usr/bin/kafka-topics --zookeeper ${ZOOKEEPER} --create \
			--partitions ${TOPIC_PATITIONS} \
			--replication-factor ${TOPIC_REPLICATION} \
			--topic ${TOPIC_NAME} 2>/dev/null
		if [ $? -ne 0 ]; then
			echo "pre-install.sh create topic ${TOPIC_NAME}(p:${TOPIC_PATITIONS},r:${TOPIC_REPLICATION}) failed"
			exit 1
		fi
	fi
	exit 0
fi
