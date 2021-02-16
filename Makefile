test:
	docker-compose up -d && docker-compose exec \
										-e MYSQL_HOST=mysql \
										-e POSTGRES_HOST=postgres \
										-e NEO4J_HOST=neo4j \
										golang go test -race ./...

lint:
	docker run -v ${PWD}:/app -w /app golangci/golangci-lint:v1.35.2 golangci-lint run -v --timeout 20m

aws-gateway-lambda-terraform-build:
	GOOS=linux GOARCH=amd64 go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix \
		netgo -o ./learning/other/aws_api_gateway_lambda_terraform/build/bin/app \
		./learning/other/aws_api_gateway_lambda_terraform

aws-gateway-lambda-terraform-init:
	docker-compose up -d && docker-compose exec terraform terraform init infra

aws-gateway-lambda-terraform-plan:
	docker-compose up -d && docker-compose exec terraform terraform plan -var-file=terraform.tvars infra

aws-gateway-lambda-terraform-apply:
	docker-compose up -d && docker-compose exec terraform terraform apply -var-file=terraform.tvars infra
