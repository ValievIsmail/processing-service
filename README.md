### Processing-service

Service for working with transactions with implementations of logic on the database side

To work, you need to install docker-compose

1. Run the database with the command 'make up-db' 

2. Start the service with the command 'make go'

### Requests:

1. Send Transaction - POST http://localhost:8080/v1/process BODY
```
{
	"amount": 15,
	"state": "lost",
	"transaction_id": 123
}
```