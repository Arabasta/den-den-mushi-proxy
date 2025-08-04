#!/bin/bash

echo "Removing old frontend admin angular..."
rm -r ../frontend/admin/angular

npx @openapitools/openapi-generator-cli generate \
  -i ./openapi.yaml \
  -g typescript-angular \
  -o ../frontend/admin/angular \
  --additional-properties=npmName=@kei/admin-api-client,npmVersion=1.0.0,providedInRoot=true
