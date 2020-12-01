swagger-spec:
	swagger generate spec --scan-models -o ./swagger.yaml
run:
	go run main.go
swagger-ui:
	docker run --rm -it -p 8081:8080 swaggerapi/swagger-ui
