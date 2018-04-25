.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	docker-compose up -d dynamodb
	AWS_REGION=eu-west-1 go run *.go -leboncoin-start-url https://www.leboncoin.fr/ventes_immobilieres/offres/\?th\=1\&q\=appartement\&sqs\=10 -dynamodb-endpoint-url http://localhost:8000
