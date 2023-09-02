# goq
Learning project that implements distributed-message queue like [Celery]

## Todo
- [ ] Use wire in worker initialization to reduce boilerplate

## Concepts

- pubsub - implements method to communicate with the message broker
- task - push and pull task from pubsub
- worker - run task to handle the message
