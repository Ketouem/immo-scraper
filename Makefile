AWS_REGION ?= eu-west-1

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	docker-compose up -d dynamodb
	AWS_REGION=${AWS_REGION} go run *.go -mode all -leboncoin-start-url https://www.leboncoin.fr/ventes_immobilieres/offres/\?th\=1\&q\=appartement\&sqs\=10 -dynamodb-endpoint-url http://localhost:8000

.PHONY: run-collect
run-collect:
	docker-compose up -d dynamodb
	AWS_REGION=${AWS_REGION} go run *.go -mode collect -leboncoin-start-url https://www.leboncoin.fr/ventes_immobilieres/offres/\?th\=1\&q\=appartement\&sqs\=10 -dynamodb-endpoint-url http://localhost:8000

.PHONY: run-notify
run-notify:
	docker-compose up -d dynamodb
	AWS_REGION=${AWS_REGION} go run *.go -mode notify -dynamodb-endpoint-url http://localhost:8000
