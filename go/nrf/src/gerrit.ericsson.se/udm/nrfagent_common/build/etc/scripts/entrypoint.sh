#!/bin/bash

Usage () {
	echo "entrypoint.sh in image nrfagent startup with wrong args"
	echo "Usage: entrypoint.sh pre-install/init-cm/start_reg/start_ntf/start_disc"
    exit 1	
}

if [ $# -eq 0 ]; then
	Usage
fi

case $1 in 
	pre-install)
		exec /etc/nrfagent/scripts/pre-install.sh "$@"
		;;
	init-cm)
		exec /etc/nrfagent/scripts/init-cm.sh "$@"
		;;
	start_reg)
		exec /bin/nrfagentreg "$@"
		;;
	start_ntf)
		exec /bin/nrfagentntf "$@"
		;;
	start_disc)
		exec /bin/nrfagentdisc "$@"
		;;
	*)
		Usage
		;;
esac
