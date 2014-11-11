export GOPATH=$PWD
export GOBIN=$GOPATH/bin

cd src/boxviewer/
# CGO_ENABLED=0 go install --ldflags '-extldflags "-static"'
go install -compiler gccgo -gccgoflags '-static-libgo'