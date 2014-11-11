#!/bin/bash
cd /var/www/apps/boxviewer/
export GOPATH=$PWD
export GOBIN=$GOPATH/bin

cd src/boxviewer/
/usr/local/go/bin/go test
/usr/local/go/bin/go install

killall -9 boxviewer
nohup /var/www/apps/boxviewer/bin/boxviewer --key="rnii4dbpu16254m844ikzx4jxi2u8sj4" --location="/var/www/apps/boxviewer_files/" &