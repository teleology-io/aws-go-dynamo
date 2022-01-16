#! /bin/bash

docker-compose up -d --force-recreate

aws dynamodb create-table --cli-input-json file://$PWD/invitations-table.json --endpoint-url http://localhost:8000