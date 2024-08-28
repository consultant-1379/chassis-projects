VERSION?=lsv5
DEVVERSION?=latest
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/nrfdbproxy

release: clean relpush
devrelease: devpush


devbuild:
	docker rmi -f ${REGISTRY}/${PROJECT}:${DEVVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${DEVVERSION} -f ./Dockerfile.runtime .

devpush: login_check devbuild
	docker push ${REGISTRY}/${PROJECT}:${DEVVERSION}


relbuild:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${VERSION} -f ./Dockerfile.runtime .

relpush: login_check relbuild
	docker push ${REGISTRY}/${PROJECT}:${VERSION}


login_check:
	@LOGIN=`sed -n '/${REGISTRY}/{n;/auth/p}' ~/.docker/config.json | wc -l`; \
	if [ $${LOGIN} -eq 0 ]; then \
	    docker login ${REGISTRY}; \
	else \
	    echo "alreay login. run 'make login' to relogin"; \
	fi; \

login:	
	docker login ${REGISTRY}; \

clean:
	rm -f app.exe
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true

