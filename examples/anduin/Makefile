push:
	oras push --username admin --password password --insecure localhost:5000/hello:v1 ./config.json

clean:
	docker-compose down -v && \
		rm -rf /tmp/zot && \
		docker-compose up -d
