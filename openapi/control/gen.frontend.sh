#!/bin/bash

echo "Removing old control api client..."
rm -r ../frontend/control-api-client

npx @openapitools/openapi-generator-cli generate \
  -i ./openapi.yaml \
  -g typescript-axios \
  -o ../frontend/control-api-client \
  --additional-properties=useSingleRequestParameter=true,withSeparateModelsAndApi=true,modelPropertyNaming=original,typescriptThreePlus=true,apiPackage=api,modelPackage=models \
