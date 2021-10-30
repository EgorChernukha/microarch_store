Домашнее задание https://otus.ru/lessons/microservice-architecture

Настройка окружения:
```
https://kubernetes.io/ru/docs/tasks/tools/install-minikube/ # minikube
https://www.virtualbox.org/wiki/Downloads # virtualbox
https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/ # minikube ingress
```

Как запустить:
```shell
make
kubectl create namespace microarch-test
kubectl apply -f ./conf/k8s/base/.
```

Домашнее задание 1:
```shell
ip=$(minikube ip) && curl -H 'Host: arch.homework' http://$ip/health
# curl -H 'Host: arch.homework' http://192.168.99.100/health
```