## Setup
```
$ sudo ip tuntap add mode tap name tap0
$ sudo ip addr add 192.0.2.1/24 dev tap0
$ sudo ip link set tap0 up
```