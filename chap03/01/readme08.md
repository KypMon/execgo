# For assignment week8 

```shell
# build the image
docker build . -t puyuyang/golearning:chap08

# push to docker hub
docker push puyuyang/golearning:chap08

# start the server in container
docker run -d puyuyang/golearning:chap08
67f6af1e628b23dcb26d02ac0ca7428f1ece39ebbdbeb7b591cbae87a4ac78ab

# check network config
PID=$(docker inspect --format {{.State.Pid}} 67f6af1e628b23dcb26d02ac0ca7428f1ece39ebbdbeb7b591cbae87a4ac78ab)
nsenter --target $PID --mount --uts --ipc --net --pid

# === Result === 
# sudo nsenter -t $PID -n ip addr
# 1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
#    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
#    inet 127.0.0.1/8 scope host lo
#       valid_lft forever preferred_lft forever
# 23: eth0@if24: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
#    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
#    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
#       valid_lft forever preferred_lft forever
```

