#!/bin/sh
export PROJECT="nrf"
export PROVPROJECT="nrfprov"
export PMPROJECT="nrfpmjobloader"
export COMMONCHART="eric-nrf-common"
export REPOPATH="proj-ipworks"
export VERSION="0.28.1"
# export USERNAME="ipworks1"
# export PASSWORD="De5PAchanedUb4eW"

if [ -n "$GERRIT_PATCHSET_REVISION" ] ; then
  export IMAGEVERSION=`echo $GERRIT_PATCHSET_REVISION | cut -c1-7`
  export ARMDOCKER="armdocker.rnd.ericsson.se/proj-ipworks/precommit"
  export REPOPATH="proj-ipworks/precommit"
  export DOCKERREG="https://arm.epk.ericsson.se/artifactory/docker-v2-global-local/proj-ipworks/precommit"
  export EMAIL=$GERRIT_CHANGE_OWNER_EMAIL
  export CHANGE="https://gerrit.ericsson.se/$GERRIT_CHANGE_NUMBER"
else
  if [ "$IMAGETAG" != "" ] ; then
    export IMAGEVERSION=$IMAGETAG
  else
    export IMAGEVERSION="DDATE"
  fi
  export ARMDOCKER="armdocker.rnd.ericsson.se/proj-ipworks"
  export REPOPATH="proj-ipworks"
  export DOCKERREG="https://arm.epk.ericsson.se/artifactory/docker-v2-global-local/proj-ipworks"
  export EMAIL="PDLIPW5GPR@pdl.internal.ericsson.com"
fi

