#!/bin/bash

docker compose -f docker-compose-dev.yaml up -d server

MENSAJE="mensaje"
RED="tp0_testing_net"

RESPUESTA=$(docker run --rm --network "$RED" busybox  \
            sh -c "echo '$MENSAJE' | nc server 12345")

if [ "$RESPUESTA" = "$MENSAJE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi

docker compose -f docker-compose-dev.yaml stop server -t 1
docker compose -f docker-compose-dev.yaml down
