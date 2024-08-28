VERSION?=v3
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/rhel-tools
TOOLVERSION?=v1
TOOLPROJECT?=proj-ipworks/buildtool

build:
	sudo docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true
	sudo docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${VERSION} -f ./Dockerfile ../..	

push: login_check build
	docker push ${REGISTRY}/${PROJECT}:${VERSION}
	
toolbuild:
	sudo docker rmi -f ${REGISTRY}/${TOOLPROJECT}:${TOOLVERSION} || true
	sudo docker build --no-cache=true -t ${REGISTRY}/${TOOLPROJECT}:${TOOLVERSION} -f ./Dockerfile.buildtool ../.. 
	
toolpush:toolbuild
	docker push ${REGISTRY}/${TOOLPROJECT}:${TOOLVERSION}	
	

clean:
	sudo docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true

login_check:
	@LOGIN=`sed -n '/${REGISTRY}/{n;/auth/p}' ~/.docker/config.json | wc -l`; \
	if [ $${LOGIN} -eq 0 ]; then \
	    docker login ${REGISTRY}; \
	else \
	    echo "alreay login. run 'make login' to relogin"; \
	fi; \

login:	
	docker login ${REGISTRY}; \

