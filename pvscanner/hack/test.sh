#! /bin/bash

# docker load -i pvscanner.img

docker run  --privileged -v /:/host/root -v /proc:/host/proc pvdf/pvscanner:latest \
/pvscanner \
--logLevel=INFO \
--rootFsPath=/host/root \
--period=10s \
--lvmdConfigPath=/host/root/etc/topolvm/lvmd.yaml \
--nodeName=kspray2 \
--kubeconfig=/host/root/root/.kube/config \
--containerized \
--topolvm \
#--nsenter=/host/root/usr/bin/nsenter \
#--lvm=/host/root/sbin/lvm


# --procPath=/host/proc \

# nsenter --root=/host/root  lvm vgs --reportformat json --unbuffered --unit b