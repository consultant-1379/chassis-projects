VERSION?=lsv5-0.5.10
DEVVERSION?=dev
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/nrfagentinitcm

BUILDPATH?=build
DEPLOYPATH?=helm/eric-nrf-register-agent

devbuild:
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${DEVVERSION} -f ./${BUILDPATH}/Dockerfile.runtime .

devpush: devbuild
	docker push ${REGISTRY}/${PROJECT}:${DEVVERSION}

relbuild:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${VERSION} -f ./${BUILDPATH}/Dockerfile.runtime .

relpush: relbuild
	docker push ${REGISTRY}/${PROJECT}:${VERSION}

relclean:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true

devclean:
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true

clean:
	rm -f nrfagentinitcm || true
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true

