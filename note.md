# Bridge

## Problem
we run a minikube in virtual box with network `vboxnet1` with network 192.168.99.0/24, we run adjust-mini services within minikube. we try to emulate kafka as an external service. so we run it with docker-compose in localhost and map to host network(docker-compose-kafka.yaml). but the host network has dynamic ip address depends on your WiFi environment. so from within adjust-mini service, we need to always change the ip address of the kafka service.

## Solution
we create a bridge between wifi interface(en0) and vbox(vboxnet1). 

```sh
sudo ifconfig bridge1 create
sudo ifconfig bridge1 addm en0 addm vboxnet1
sudo ifconfig bridge1 up
```

and we assign an IP for the bridge
```sh
sudo ifconfig bridge1 192.168.2.1
```

now we can update `docker-compose-kafka.yaml` to let kafka advertise this ip
address and also in adjust-mini service use this ip address as configuration

