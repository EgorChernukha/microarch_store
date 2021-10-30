Домашнее задание https://otus.ru/lessons/microservice-architecture

Как запустить:
```shell
make
kubectl create namespace microarch-test
kubectl apply -f ./conf/k8s/base/.
```

Домашнее задание 1:
```shell
ip=$(minikube ip)
curl -H 'Host: arch.homework' http://$ip/health
# 192.168.99.100
# curl -H 'Host: arch.homework' http://192.168.99.100/health
```