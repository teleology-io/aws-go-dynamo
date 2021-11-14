# How to test

In order to test locally we need to use docker and the docker-compose.yml file in the project. 

To start downloading and running dynamodb locally, run:
```bash
docker-compopse up -d --force-recreate
```

Now create the test table with:
```bash
aws dynamodb create-table --cli-input-json file://invitations-table.json --endpoint-url http://localhost:8000
```

And to verify the table was created:
```bash
aws dynamodb list-tables --endpoint-url http://localhost:8000
```

Your output should look similar to this:
```json
{
    "TableNames": [
        "foundation-local-invitations"
    ]
}
```

