# Task Queue with Kafka

API publishes tasks to Kafka, workers consume and process in parallel. 

Built this to understand message queues, consumer groups, and how services scale independently by decoupling input from execution.

```bash
docker-compose up
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"type":"report"}'
```
