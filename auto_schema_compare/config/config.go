package config

import "fmt"

var (
	SwaggerUrl             = "https://s3.dualstack.us-east-1.amazonaws.com/inin-prod-api/us-east-1/public-api-v2/swagger-schema/publicapi-v2-latest.json"
	DirName                = "swagger_files"
	NewSwaggerFile         = "new_swagger.json"
	CurrentSwaggerFile     = "current_swagger.json"
	NewSwaggerFilePath     = fmt.Sprintf("%s/%s", DirName, NewSwaggerFile)
	CurrentSwaggerFilePath = fmt.Sprintf("%s/%s", DirName, CurrentSwaggerFile)
)