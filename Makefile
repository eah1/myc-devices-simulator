
include deployment/mk/test-go.mk
include deployment/mk/database.mk
include deployment/mk/run.mk

swagger-v1:
	swag init --instanceName v1 --dir "business/infra/handlers/v1"  --output "business/infra/handlers/v1/docs"  --generalInfo group.go  --parseDependency true