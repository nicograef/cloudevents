# Queue
A simple implementation of a message queue.

## build & run
```sh
docker build -t github.com/nicograef/cloudevents/queue .

docker run -p 3000:3000 github.com/nicograef/cloudevents/queue
```

## test api
Add a new message to the queue.
```sh
curl -X POST localhost:3000/enqueue -d '{"payload": "this is some data"}'
```

Get the current/next queue message.
```sh
curl -i -X POST localhost:3000/peek
```

Remove a processed message from the queue
```sh
curl -i -X POST localhost:3000/dequeue -d '{"messageId":"fe5f38f3-1067-4c9b-a542-2bf9f16ca08d"}'
```