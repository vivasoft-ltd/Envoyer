#!/bin/bash

trap "rabbitmqctl stop; exit" SIGHUP SIGINT SIGTERM
rabbitmq-server &
sleep 10
service mysql start
./app migrate
./app server &
yarn --cwd envoyer_frontend start &
wait
