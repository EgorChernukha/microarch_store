#!/bin/bash

#curl -L -X POST 'arch.homework/api/v1/user' \-H 'Content-Type: application/json' --data-raw '{"username": "echernukha","firstname": "egor","lastname": "chernukha","email": "echernukha@mail.com","phone": "+78005553535"}'

while true; do
    # get existent user
    ab -n 1000 -c 50 http://arch.homework/api/v1/user/42b96cc5-4c74-11ec-bc9b-0242ac110009
    # get non-existent user
    ERR_COUNT=$((1 + $RANDOM % 50))
    ab -n $ERR_COUNT -c $ERR_COUNT http://arch.homework/api/v1/user/00000000-0000-0000-0000-000000000000
    sleep 3
done