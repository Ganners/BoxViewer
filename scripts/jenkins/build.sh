export GOPATH=$PWD
export GOBIN=$GOPATH/bin

cd src/boxviewer/
CGO_ENABLED=0 go install