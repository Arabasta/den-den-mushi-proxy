#!/bin/bash

echo "Removing old frontend angular..."
rm -r ../frontend/angular

npx @openapitools/openapi-generator-cli generate \
  -i ./openapi.yaml \
  -g typescript-angular \
  -o ../frontend/angular \
  --additional-properties=npmName=@yourorg/control-api-client,npmVersion=1.0.0,providedInRoot=true
