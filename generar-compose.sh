#!/bin/bash

NOMBRE_ARCHIVO=$1
CLIENTES=$2

if [[ $# -ne 2 ]]; then
    echo "ERROR: Nombre del archivo de salida o número de clientes no proporcionados."
    echo "Uso: $0 <nombre_archivo> <numero_clientes>"
    exit 1
fi

echo "name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - CLIENTES=$CLIENTES
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini" > $NOMBRE_ARCHIVO

for i in $(seq 1 $CLIENTES)
do
    echo "
  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-$i.csv:/data/data.csv" >> $NOMBRE_ARCHIVO
done

echo "
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24" >> $NOMBRE_ARCHIVO

echo "Archivo $NOMBRE_ARCHIVO generado con éxito."
