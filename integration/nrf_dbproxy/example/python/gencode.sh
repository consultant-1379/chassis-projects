mkdir -p dbproxy
find ../../java/src/main/proto -type f | xargs -t -I {} python -m grpc_tools.protoc --proto_path=../../java/src/main/proto --python_out=dbproxy {}
python -m grpc_tools.protoc -I ../../java/src/main/proto --python_out=dbproxy --grpc_python_out=dbproxy NFDataManagementService.proto 
VERSION=`git log --oneline | head -1 | awk '{print $1}'`
echo $VERSION > dbproxy/nfmessage/.git_version
cp -r dbproxy/* ../../../NRF_TestSuite/Library/
