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
kubectl apply -f ./conf/k8s/base/namespace.yml
kubectl apply -f ./conf/k8s/base/.
```

Домашнее задание 1:
```shell
ip=$(minikube ip) && curl -H 'Host: arch.homework' http://$ip/health
# curl -H 'Host: arch.homework' http://192.168.99.100/health
```

**hw2**:
`Инфраструктурные паттерны`

**Цель**:
`В этом ДЗ вы создадите простейший RESTful CRUD.`

1. Сделать простейший RESTful CRUD по созданию, удалению, просмотру и обновлению пользователей.
Пример API  - https://app.swaggerhub.com/apis/otus55/users/1.0.0

2. Добавить базу данных для приложения.
Конфигурация приложения должна хранится в Configmaps.
Доступы к БД должны храниться в Secrets.
Первоначальные миграции должны быть оформлены в качестве Job-ы, если это требуется.
Ingress-ы должны также вести на url arch.homework/ (как и в прошлом задании)

**На выходе должны быть предоставлена**

1. ссылка на директорию в github, где находится директория с манифестами кубернетеса
инструкция по запуску приложения.
2. команда установки БД из helm, вместе с файлом values.yaml.
3. команда применения первоначальных миграций
4. команда kubectl apply -f, которая запускает в правильном порядке манифесты кубернетеса
Postman коллекция, в которой будут представлены примеры запросов к сервису на создание, получение, изменение и удаление пользователя. Важно: в postman коллекции использовать базовый url - `arch.homework`.

Задание со звездочкой:
+5 балла за шаблонизацию приложения в helm чартах