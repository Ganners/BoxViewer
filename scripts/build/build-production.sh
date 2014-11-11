#!/bin/bash
cd /var/www/apps/boxviewer/
export GOPATH=$PWD
export GOBIN=$GOPATH/bin

cd src/boxviewer/
/usr/local/go/bin test
/usr/local/go/bin install

killall -9 /var/www/apps/boxviewer/bin/boxviewer
/var/www/apps/boxviewer/bin/boxviewer --key="rnii4dbpu16254m844ikzx4jxi2u8sj4" --location="files"