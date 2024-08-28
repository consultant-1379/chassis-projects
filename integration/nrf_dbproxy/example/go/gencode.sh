find ../../java/src/main/proto -type f | xargs -t -I {} protoc --proto_path=../../java/src/main/proto --go_out=plugins=grpc:. {}
VERSION=`git log --oneline | head -1 | awk '{print $1}'`
echo $VERSION > com/dbproxy/.git_version
cp -r com ../../../vendor
