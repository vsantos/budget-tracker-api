#!/bin/bash

curl -XPOST http://localhost:8080/api/v1/user \
    -H 'content-type: application/json' \
    -d '{"login": "vsantos", "firstname":"Victor", "lastname":"Santos", "email": "vsantos.py@gmail.com", "password": "randompass"}'


curl -XPOST http://localhost:8080/api/v1/user \
    -H 'content-type: application/json' \
    -d '{"login": "apinheiro", "firstname":"Alberto", "lastname":"Pinheiros", "email": "apinheiro@hotmail.com.br", "password": "myawesome"}'
