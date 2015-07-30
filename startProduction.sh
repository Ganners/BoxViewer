#!/bin/bash

go install
killall -9 BoxViewer
nohup ~/go/bin/BoxViewer\
  --key=rnii4dbpu16254m844ikzx4jxi2u8sj4\
  --location=/home/cyber/BoxViewerFiles\
  --keyFile=/root/privkey.pem\
  --certFile=/root/STAR_onmojo_co_uk/certs.crt\
  --port=8050 &
