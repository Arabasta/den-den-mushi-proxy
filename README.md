# Den Den Mushi Proxy

Proxy server thing

1. Build all

```bash
make
```

2. Run

```bash
make run CMD=proxy
```

3. Clean

```bash
make clean
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