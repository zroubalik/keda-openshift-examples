# Example Kafka Consumer and Producer applications

## Kafka Producer

### Relevant Environment Variables
```yaml
BOOTSTRAP_SERVERS=kafka-xxxxxx.com:443
SASL_USER=<Client ID>
SASL_PASSWORD=<Client Secret>
TOPIC=my-topic
GROUP_ID=my-group
SASL=enabled
MESSAGE_COUNT=100
DELAY_MS=100
```

### Build the Go app
```bash
make producer
```

### Test it!
Fill [set-env.sh](set-env.sh) appropriatelly and source it:
```bash
source ./set-env.sh
```
Run the app:
```bash
./kafkaproducerapp
2022/02/09 16:45:45 Go producer starting with config=&{BootstrapServers:kafka-xxxxxxx.kafka.rhcloud.com:443 Topic:my-topic Delay:100 Message:Hello from Go Kafka Sarama MessageCount:100 ProducerAcks:1 SaslEnabled:true SaslUser:xxxxxxx SaslPassword:xxxxxxx}
2022/02/09 16:45:46 Sending message: value=Hello from Go Kafka Sarama-0
2022/02/09 16:45:47 Message sent: partition=0, offset=174
2022/02/09 16:45:47 Sending message: value=Hello from Go Kafka Sarama-1
2022/02/09 16:45:47 Message sent: partition=0, offset=175
2022/02/09 16:45:47 Sending message: value=Hello from Go Kafka Sarama-2
2022/02/09 16:45:48 Message sent: partition=8, offset=160
2022/02/09 16:45:48 Sending message: value=Hello from Go Kafka Sarama-3
...
```

## Kafka Consumer

### Relevant Environment Variables
```yaml
BOOTSTRAP_SERVERS=kafka-xxxxxx.com:443
SASL_USER=<Client ID>
SASL_PASSWORD=<Client Secret>
TOPIC=my-topic
GROUP_ID=my-group
SASL=enabled
```

### Build the Go app
```bash
make consumer
```

### Test it!
Fill [set-env.sh](set-env.sh) appropriatelly and source it:
```bash
source ./set-env.sh
```
Run the app:
```bash
./kafkaconsumerapp
2022/02/09 16:47:48 Go consumer starting with config=&{BootstrapServers:kafka-xxxxxx.bf2.kafka.rhcloud.com:443 Topic:my-topic GroupID:my-group SaslEnabled:true SaslUser:xxxxxx SaslPassword:xxxxxx}
2022/02/09 16:47:55 Consumer group handler setup
2022/02/09 16:47:55 Sarama consumer up and running!...
2022/02/09 16:47:56 Message received: value=Hello from Go Kafka Sarama-4, topic=my-topic, partition=7, offset=179
2022/02/09 16:47:56 Message received: value=Hello from Go Kafka Sarama-5, topic=my-topic, partition=1, offset=194
2022/02/09 16:47:56 Message received: value=Hello from Go Kafka Sarama-3, topic=my-topic, partition=4, offset=200
2022/02/09 16:47:56 Message received: value=Hello from Go Kafka Sarama-6, topic=my-topic, partition=4, offset=201
2022/02/09 16:47:57 Message received: value=Hello from Go Kafka Sarama-0, topic=my-topic, partition=0, offset=174
2022/02/09 16:47:57 Message received: value=Hello from Go Kafka Sarama-1, topic=my-topic, partition=0, offset=175
...
```
