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
helm install mysql-user -f conf/helm/mysql/uservalues.yaml bitnami/mysql --version 8.8.12

kubectl apply -f ./conf/k8s/base/.

# Запуск тестов
newman run ./data/lab-5.tests.postman_collection.json
```

```shell
Stream processing

Цель:
В этом ДЗ вы научитесь реализовывать сервис заказа.

Реализовать сервис заказа. Сервис биллинга. Сервис нотификаций. 

При создании пользователя, необходимо создавать аккаунт в сервисе биллинга. В сервисе биллинга должна быть возможность положить деньги на аккаунт и снять деньги. 

Сервис нотификаций позволяет отправить сообщение на email. И позволяет получить список сообщений по методу API. 

Пользователь может создать заказ. У заказа есть  параметр - цена заказа. 
Заказ происходит в 2 этапа:

сначала снимаем деньги с пользователя с помощью сервиса биллинга
отсылаем пользователю сообщение на почту с результатами оформления заказа. Если биллинг подтвердил платеж, должно отослаться письмо счастья. Если нет, то письмо горя.
Упрощаем и считаем, что ничего плохого с сервисами происходить не может (они не могут падать и т.д.). Сервис нотификаций на самом деле не отправляет, а просто сохраняет в БД.

ТЕОРЕТИЧЕСКАЯ ЧАСТЬ (5 баллов):
0) Спроектировать взаимодействие сервисов при создании заказов. Предоставить варианты взаимодействий в следующих стилях в виде sequence диаграммы с описанием API на IDL:

только HTTP взаимодействие
событийное взаимодействие с использование брокера сообщений для нотификаций (уведомлений)
Event Collaboration cтиль взаимодействия с использованием брокера сообщений
вариант, который вам кажется наиболее адекватным для решения данной задачи. Если он совпадает одним из вариантов выше - просто отметить это.
ПРАКТИЧЕСКАЯ ЧАСТЬ (5 баллов):
Выбрать один из вариантов и реализовать его. 
На выходе должны быть
0) описание архитектурного решения и схема взаимодействия сервисов (в виде картинки)

команда установки приложения (из helm-а или из манифестов). Обязательно указать в каком namespace нужно устанавливать.
тесты постмана, которые прогоняют сценарий:
Создать пользователя. Должен создаться аккаунт в биллинге.
Положить деньги на счет пользователя через сервис биллинга.
Сделать заказ, на который хватает денег.
Посмотреть деньги на счету пользователя и убедиться, что их сняли.
Посмотреть в сервисе нотификаций отправленные сообщения и убедиться, что сообщение отправилось
Сделать заказ, на который не хватает денег.
Посмотреть деньги на счету пользователя и убедиться, что их количество не поменялось.
Посмотреть в сервисе нотификаций отправленные сообщения и убедиться, что сообщение отправилось.
  В тестах обязательно 

наличие {{baseUrl}} для урла
использование домена arch.homework в качестве initial значения {{baseUrl}}
отображение данных запроса и данных ответа при запуске из командной строки с помощью newman.
```

Домашнее задание 5:
```shell
# Настройка окружения
kubectl apply -f ./conf/k8s/base/namespace.yaml
kubens otus

# app
helm install mysql-auth -f conf/helm/mysql/authvalues.yaml bitnami/mysql --version 8.8.12
helm install mysql-user -f conf/helm/mysql/uservalues.yaml bitnami/mysql --version 8.8.12
helm install mysql-order -f conf/helm/mysql/ordervalues.yaml bitnami/mysql --version 8.8.12
helm install mysql-billing -f conf/helm/mysql/billingvalues.yaml bitnami/mysql --version 8.8.12
helm install mysql-notification -f conf/helm/mysql/notificationvalues.yaml bitnami/mysql --version 8.8.12

# monitoring
helm install prom prometheus-community/kube-prometheus-stack -f ./conf/helm/prometheus/values.yaml --atomic
helm install nginx ingress-nginx/ingress-nginx -f ./conf/helm/ingress-nginx/values.yaml --atomic

#app
kubectl apply -f ./conf/k8s/base/. --recursive

#monitoring
kubectl apply -f ./conf/k8s/monitoring/. --recursive
```