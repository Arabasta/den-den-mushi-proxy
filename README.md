# Den Den Mushi Proxy

Proxy server thing

## Build all
```bash
make
```
Outputs binaries in `bin/` directory.

## Run without Building

From root run:
```bash
make run CMD=proxy
```
```bash
make run CMD=control
```

## OpenAPI

Docs: http://localhost:55007/swagger/control <br>
JSON: http://localhost:55007/swagger-spec/control.json

### Generate Golang Server Stub
[gen.go](openapi/control/gen.go)

```bash
make generate
```

### Generate Angular code
```bash
make generate-client
```
## Load Test (switch to load-test branch)

Apache JMeter 5.6.3

1. Download Jmeter
   https://jmeter.apache.org/download_jmeter.cgi

2. Install Plugin Manager
   https://jmeter-plugins.org/install/Install/https://jmeter-plugins.org/install/Install/

3. Install WebSocket Samplers by Peter Doornbosch

   a. Open JMeter

   b. Options > Plugins Manager > Available Plugins

   c. Apply changes and restart

4. Run JMeter

5. Open [simple_load_test.jmx](simple_load_test.jmx)