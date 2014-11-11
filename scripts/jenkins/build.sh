export GOPATH=$PWD
export GOBIN=$GOPATH/bin

cd src/boxviewer/
go test
go install