Домашнее задание https://otus.ru/lessons/microservice-architecture

Настройка окружения:
```
https://kubernetes.io/ru/docs/tasks/tools/install-minikube/ # minikube
https://www.virtualbox.org/wiki/Downloads # virtualbox
https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/ # minikube ingress
```

Домашнее задание 1:
```shell
ip=$(minikube ip) && curl -H 'Host: arch.homework' http://$ip/health
# curl -H 'Host: arch.homework' http://192.168.99.100/health
```

Как запустить:
```shell
kubectl apply -f ./conf/k8s/base/namespace.yaml && kubectl apply -f ./conf/k8s/base/.
```


Домашнее задание 2:
```shell
# Установка mysql
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install mysql -f conf/helm/mysql/values.yaml bitnami/mysql --version 8.8.12

# Запуск
kubectl apply -f ./conf/k8s/base/namespace.yaml
kubectl apply -f ./conf/k8s/base/.
```


Домашнее задание 3:
```shell
# Установка prometheus и nginx
minikube addons disable ingress
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install prom prometheus-community/kube-prometheus-stack -f ./conf/helm/prometheus/values.yaml --atomic
helm install nginx ingress-nginx/ingress-nginx -f ./conf/helm/ingress-nginx/values.yaml --atomic

# Запуск grafana. Перейти на http://localhost:9000(admin: prom-operator)
kubectl port-forward service/prom-grafana 9000:80

# В grafana необходимо импортировать dashboard из /data/grafana/dashboard.json

# Запуск стресс-теста
./data/test.sh
# 
```


Домашнее задание 5:
```shell
# Настройка окружения
kubectl apply -f ./conf/k8s/base/namespace.yaml
kubens otus
helm install prom prometheus-community/kube-prometheus-stack -f ./conf/helm/prometheus/values.yaml --atomic
helm install nginx ingress-nginx/ingress-nginx -f ./conf/helm/ingress-nginx/values.yaml --atomic
helm install mysql-auth -f conf/helm/mysql/authvalues.yaml bitnami/mysql --version 8.8.12
helm install mysql-store -f conf/helm/mysql/storevalues.yaml bitnami/mysql --version 8.8.12

kubectl apply -f ./conf/k8s/base/.


```