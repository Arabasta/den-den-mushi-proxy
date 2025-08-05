#!/bin/bash

echo "Removing old frontend control angular..."
rm -r ../frontend/control/angular

npx @openapitools/openapi-generator-cli generate \
  -i ./openapi.yaml \
  -g typescript-angular \
  -o ../frontend/control/angular \
  --additional-properties=npmName=@kei/control-api-client,npmVersion=1.0.0,providedInRoot=true
