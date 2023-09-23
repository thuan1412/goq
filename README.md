# goq
Learning project that implements distributed-message queue like [Celery]

## Todo
- [ ] Add retry for task
- [ ] Add task result for monitoring
- [ ] Use wire in worker initialization to reduce boilerplate
- [ ] Use github CI to check lint
- [x] Use testcontainer to implement integration test

## Concepts

- pubsub - implements method to communicate with the message broker
- task - push and pull task from pubsub
- worker - run task to handle the message
