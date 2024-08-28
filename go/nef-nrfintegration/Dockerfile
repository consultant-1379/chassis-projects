FROM  selidocker.lmera.ericsson.se/proj-nef/base-image-redhat-golang-builder:1.12.5 as go-builder
ENV GO111MODULE=on
RUN mkdir -p /go/src/gerrit.ericsson.se/nef/
COPY ./src/ /go/src/gerrit.ericsson.se/nef/
COPY ./go.mod /go/src/gerrit.ericsson.se/nef/
WORKDIR /go/src/gerrit.ericsson.se/nef/
RUN mkdir -p /root/.ssh/
ADD ssh/ /root/.ssh/
RUN chmod 600 -R /root/.ssh/
RUN git config --global url."ssh://eupgbld@gerrit.ericsson.se:29418/nef".insteadOf https://gerrit.ericsson.se/a/nef
RUN git config --global url."ssh://eupgbld@gerrit.ericsson.se:29418/nef/nef-golangcommon".insteadOf https://gerrit.ericsson.se/a/nef/nef-golangcommon

#COPY ./go.sum /go/s/gerrit.ericsson.se/nef/
RUN cd /go/src/gerrit.ericsson.se/nef/ ; \
    go test -v -cover ./...; \
    # the -o path is going to be /go/src/gerrit.ericsson.se/nef/ somewhere then
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /bin/nef-nrfintegration . ;

FROM armdocker.rnd.ericsson.se/proj-ldc/lc_sles15:latest
COPY --from=go-builder /bin/nef-nrfintegration  /bin/nef-nrfintegration
CMD ["/bin/nef-nrfintegration"]