# qugo

Qugo is an async message queue for event-driven applications.

Following the [CNCF Cloudevents specification](https://cloudevents.io), qugo is a message queue service written in go that receives event messages via http and pushes them to a configured webhook. Qugo uses a go channel to decouple producer and consumer. Qugo waits for the result of the consumer before deleting the event from its queue.

Qugo is packaged in a docker container that can be started by `docker run @nicograef/qugo`

## Configuration

Configuration is done via cli parameters or enviornment variables.

ENV:

- `QUEUE_SIZE=1000`
- `CONSUMER_URL=http://localhost:4000`

CLI:

- `--queue-size 1000`
- `--consumer-url http://localhost:4000`

## Usage

```bash
docker run \
-e QUEUE_SIZE='1000' \
-e CONSUMER_URL='http://localhost/4000' \
-p 3000:3000 \
--name qugo github.com/nicograef/qugo
```


## Development

### build & run

```sh
QUEUE_SIZE=1000 CONSUMER_URL=http://localhost/4000 go run . 
```

```sh
docker build -t github.com/nicograef/qugo .

docker run -p 3000:3000 github.com/nicograef/qugo
```

### test api

Add a new message to the queue.

```sh
curl -X POST localhost:3000 -d '{"payload": "this is some data"}'
```

## Related stuff

- [Cloudevens spec](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md)
- [@thenativeweb](https://github.com/thenativeweb)'s [EventsourcingDB](https://www.thenativeweb.io/products/eventsourcingdb)

## License

MIT
