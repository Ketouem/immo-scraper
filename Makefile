AWS_REGION ?= eu-west-1
LEBONCOIN_START_URL ?= https://www.leboncoin.fr/ventes_immobilieres/offres/ile_de_france/hauts_de_seine/?th=1&location=Boulogne-Billancourt%2092100&pe=24&sqs=9&ret=2

.PHONY: test
test:
	go test ./...

.PHONY: run
run:
	docker-compose up -d dynamodb
	AWS_REGION=${AWS_REGION} go run *.go -mode all -leboncoin-start-url "${LEBONCOIN_START_URL}" -dynamodb-endpoint-url http://localhost:8000

.PHONY: run-collect
run-collect:
	docker-compose up -d dynamodb
	AWS_REGION=${AWS_REGION} go run *.go -mode collect -leboncoin-start-url "${LEBONCOIN_START_URL}" -dynamodb-endpoint-url http://localhost:8000

.PHONY: run-notify
run-notify:
	docker-compose up -d dynamodb
	AWS_REGION=${AWS_REGION} go run *.go -mode notify -dynamodb-endpoint-url http://localhost:8000 -verbose

clean:
	docker-compose kill dynamodb
	docker-compose rm -fv dynamodb
