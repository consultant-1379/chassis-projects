VERSION?=v1
REGISTRY?=armdocker.rnd.ericsson.se
PROJECT?=proj-ipworks/sles-tools
SLESVERSION?=v2
SLESPROJECT?=proj-ipworks/sles-java-tools
TOOLVERSION?=v3
TOOLPROJECT?=proj-ipworks/buildtool
JAVATOOLVERSION?=v2
JAVATOOLPROJECT?=proj-ipworks/buildtool-java

build:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${PROJECT}:${VERSION} -f ./Dockerfile.gobase ../..

push: login_check build
	docker push ${REGISTRY}/${PROJECT}:${VERSION}

javabasebuild:
	docker rmi -f ${REGISTRY}/${SLESPROJECT}:${SLESVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${SLESPROJECT}:${SLESVERSION} -f ./Dockerfile.javabase ../..

javabasepush: javabasebuild
	docker push ${REGISTRY}/${SLESPROJECT}:${SLESVERSION}


gotoolbuild:
	docker rmi -f ${REGISTRY}/${TOOLPROJECT}:${TOOLVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${TOOLPROJECT}:${TOOLVERSION} -f ./Dockerfile.gobuild ../..

gotoolpush:gotoolbuild
	docker push ${REGISTRY}/${TOOLPROJECT}:${TOOLVERSION}

javatoolbuild:
	docker rmi -f ${REGISTRY}/${JAVATOOLPROJECT}:${JAVATOOLVERSION} || true
	docker build --no-cache=true -t ${REGISTRY}/${JAVATOOLPROJECT}:${JAVATOOLVERSION} -f ./Dockerfile.javabuild ../..

javatoolpush:javatoolbuild
	docker push ${REGISTRY}/${JAVATOOLPROJECT}:${JAVATOOLVERSION}


clean:
	docker rmi -f ${REGISTRY}/${PROJECT}:${VERSION} || true

login_check:
	@LOGIN=`sed -n '/${REGISTRY}/{n;/auth/p}' ~/.docker/config.json | wc -l`; \
	if [ $${LOGIN} -eq 0 ]; then \
	    docker login ${REGISTRY}; \
	else \
	    echo "alreay login. run 'make login' to relogin"; \
	fi; \

login:
	docker login ${REGISTRY}; \
