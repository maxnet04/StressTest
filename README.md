# stress-test-cli

The technical challenger of my master degree in golang.

## How to run

Build the docker image:
```bash
docker build -t stress-test-cli .
```

After that, you can run the following command to execute the stress test:
```bash
docker run stress-test-cli --url=http://google.com --requests=10 --concurrency=5
```

## About the outputs:

Total execution time: 3.163979825s
Total requests: 250

Summary:
Status codes 200 | Count 250 (100%)