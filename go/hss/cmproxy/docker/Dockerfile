FROM            scratch

EXPOSE          9080

ARG             BUILD_TYPE
ADD             ${BUILD_TYPE}/bin/cmproxy /root/
ENTRYPOINT      ["/root/cmproxy"]
