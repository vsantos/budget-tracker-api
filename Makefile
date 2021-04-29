swagger-spec:
	swagger generate spec --scan-models -o ./docs/swagger.yaml
run:
	docker-compose up -d
rebuild:
	docker-compose up -d --force-recreate --build budget-tracker
