# Autoscaling of Kafka Consumer application connected to Kafka instance created by Red Hat Integration - AMQ Streams operator
The following guide describes the way how can be Kafka Consumer application autoscaled by KEDA on Openshift. The application is being scaled based on lag in the Kafka topic. If there isn't any traffic the application is autoscaled to 0 replicas, if there's some load the number of replicas is being scaled up to 5 replicas.

Kafka KEDA scaler is being used for this setup, for details please refer to [documentation](https://keda.sh/docs/latest/scalers/apache-kafka/).

[Appendix](#appendix-fallback) section describes `Fallback` functionality, useful when the external service that is being used to poll for a metrics (in this example a Kafka cluster) is unavailable.

## 0. Install KEDA
In `OperatorHub` locate and install KEDA, follow the instuctions to create `KedaController` instance in `keda` namespace.

## 1. Prepare Kafka Instance
 1. In `OperatorHub` locate and install `Red Hat Integration - AMQ Streams` operator
 2. Create a Kafka instance by running following command, this will create Kafka instance `my-cluster` in `openshift-operators` namespace. You can run the following command to do so:
 ```bash
 oc apply -f kafka.yaml
 ```
 3. Wait until the Kafka cluster is ready, you can monitor the `READY` status field on `my-cluster` kafka resource until it is `True`.
 ```bash
 watch oc get kafka.kafka.strimzi.io/my-cluster -n openshift-operators
 ``` 
 You should see the similar ouput, not that `READY` is `True`:
 ```bash
 Every 2,0s: oc get kafka.kafka.strimzi.io/my-cluster -n openshift-operators

 NAME         DESIRED KAFKA REPLICAS   DESIRED ZK REPLICAS   READY   WARNINGS
 my-cluster   3                        3                     True
 ```
 4. Create Kafka Topic `my-topic`, in order to allow Kafka Consumer autoscaling, you need to set the number of partitions to number greater than `1`. The number of partitions equals the maximum number of Kafka Consumer instances. In our example, set it to `5` partitions. Following command creates the topic with proper partitions count:
 ```bash
 oc apply -f topic.yaml
 ```

## 3. Deploy Kafka Consumer application
Deploy Kafka Consumer application with the following command:
 ```bash
oc apply -f deployment.yaml
 ```
Verify the consumer has been able to connect to Kafka instance, run following command:
 ```bash
oc logs deployment.apps/kafka-amqstreams-consumer
 ```
You should see similar output:
 ```bash
2022/02/09 20:36:00 Go consumer starting with config=&{BootstrapServers:my-cluster-kafka-bootstrap.openshift-operators.svc:9092 Topic:my-topic GroupID:my-group SaslEnabled:false SaslUser:user SaslPassword:password}
2022/02/09 20:36:00 Consumer group handler setup
2022/02/09 20:36:00 Sarama consumer up and running!...
 ```

## 4. Send messages to Kafka to test the Kafka Consumer application
To generate some load create this Kubernetes Job with the following command:
```bash
oc create -f load.yaml
```
Chech the logs of the Kafka Consumer application:
```bash
oc logs deployment.apps/kafka-amqstreams-consumer
```
There should be 15 messages generated and consumed:
```bash
2022/02/09 20:36:00 Go consumer starting with config=&{BootstrapServers:my-cluster-kafka-bootstrap.openshift-operators.svc:9092 Topic:my-topic GroupID:my-group SaslEnabled:false SaslUser:user SaslPassword:password}
2022/02/09 20:36:00 Consumer group handler setup
2022/02/09 20:36:00 Sarama consumer up and running!...
2022/02/09 21:02:13 Message received: value=Hello from Go Kafka Sarama-0, topic=my-topic, partition=4, offset=14
2022/02/09 21:02:13 Message received: value=Hello from Go Kafka Sarama-1, topic=my-topic, partition=2, offset=10
2022/02/09 21:02:13 Message received: value=Hello from Go Kafka Sarama-2, topic=my-topic, partition=1, offset=11
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-3, topic=my-topic, partition=3, offset=7
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-4, topic=my-topic, partition=0, offset=18
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-5, topic=my-topic, partition=1, offset=12
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-6, topic=my-topic, partition=1, offset=13
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-7, topic=my-topic, partition=1, offset=14
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-8, topic=my-topic, partition=0, offset=19
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-9, topic=my-topic, partition=3, offset=8
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-10, topic=my-topic, partition=1, offset=15
2022/02/09 21:02:14 Message received: value=Hello from Go Kafka Sarama-11, topic=my-topic, partition=4, offset=15
2022/02/09 21:02:15 Message received: value=Hello from Go Kafka Sarama-12, topic=my-topic, partition=2, offset=11
2022/02/09 21:02:15 Message received: value=Hello from Go Kafka Sarama-13, topic=my-topic, partition=2, offset=12
2022/02/09 21:02:15 Message received: value=Hello from Go Kafka Sarama-14, topic=my-topic, partition=2, offset=13
```

## 5. Deploy ScaledObject to enable Kafka Consumer application autoscaling
Then deploy a ScaledObject with the following command:
```bash
oc apply -f scaledobject.yaml
```
Check that KEDA has been able to access metrics and is correctly defined for autoscaling:
```bash
oc get scaledobject kafka-amqstreams-consumer-scaledobject
```
You should see similar output, `READY` should be `True`:
```bash
NAME                                     SCALETARGETKIND      SCALETARGETNAME             MIN   MAX   TRIGGERS   AUTHENTICATION   READY   ACTIVE   FALLBACK   AGE
kafka-amqstreams-consumer-scaledobject   apps/v1.Deployment   kafka-amqstreams-consumer   0     5     kafka                       True    False    False      17s
```
Because there aren't any messages in the Kafka topic, the Kafka Consumer application should be autoscaled to zero, run the following command:
```bash
oc get deployment.apps/kafka-amqstreams-consumer
```
You should see a similar output, `kafka-amqstreams-consumer` has been autoscaled to 0 replicas:
```bash
NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
kafka-amqstreams-consumer   0/0     0            0           11m
```

## 6. Send more messages to Kafka to test the Kafka Consumer application autoscaling
Update `MESSAGE_COUNT` environment variable in [load.yaml](load.yaml) file, increase the value from `15` to at least `500` to generate more load. Then create this Kubernetes Job with the following command:
```bash
oc create -f load.yaml
```
You should see created an increased nubmer replicas of the Kafka Consumer application until all sent messages are processed. And the the application will be again autoscaled down to zero. You can check the changing number of replicas by running the following command:
```bash
watch oc get deployment.apps/kafka-amqstreams-consumer
```

The output should be similar:
```bash
Every 2,0s: oc get deployment.apps/kafka-amqstreams-consumer
NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
kafka-amqstreams-consumer   5/5     5            5           21m

### After some time the application should be autoscaled back to 0

Every 2,0s: oc get deployment.apps/kafka-amqstreams-consumer
NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
kafka-amqstreams-consumer   0/0     0            0           23m
```

## 7. Clean up (skip this if you want to continue with APPENDIX)
Run the following commands to remove all resources created in the namespace
```bash
oc delete jobs --field-selector status.successful=1 
oc delete -f scaledobject.yaml
oc delete -f deployment.yaml
oc delete -f topic.yaml
oc delete -f kafka.yaml
``` 

## APPENDIX: Fallback
The `fallback` defines the number of replicas to fallback to if a scaler is in an error state. KEDA keeps track of the number of consecutive times each scaler has failed to get metrics from its source. Once that value passes the `fallback.failureThreshold` it scales to `fallback.replicas`.
In this example the application will be scaled to `2` replicas if there are `3` consecutive failures:
```yaml
  fallback:
    failureThreshold: 3
    replicas: 2
```

Let's introduce a failure to previously defined `kafka` scaler to see `fallback` functionality in action.

0. Delete the previously created ScaledObject:
```bash
oc delete -f scaledobject.yaml
```

1. Create a new ScaledObject with `fallback` specified:
```bash
oc apply -f scaledobject-fallback.yaml
```

2. Introduce a failure in the external service, ie. delete Kafka topic and cluster:
```bash
oc delete -f topic.yaml
oc delete -f kafka.yaml
```

3. Watch the Kafka Consumer application to check that number of replicas, it should start with 0 replicas if there is no load:
```bash
watch oc get deployment.apps/kafka-amqstreams-consumer
```

4. After some time it should scale to 2 replicas (fallback). In this specific example those replicas aren't indeed fully ready because they are not able to connect to non existent Kafka cluster.
```bash
Every 2,0s: oc get deployment.apps/kafka-amqstreams-consumer

NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
kafka-amqstreams-consumer   0/2     2            0           21m
```

5. Confirm the status of ScaleObject:
```bash 
oc get scaledobject kafka-amqstreams-consumer-scaledobject
```

6. Note the `FALLBACK` status set to `True`:
```bash
NAME                                     SCALETARGETKIND      SCALETARGETNAME             MIN   MAX   TRIGGERS   AUTHENTICATION   READY   ACTIVE   FALLBACK   AGE
kafka-amqstreams-consumer-scaledobject   apps/v1.Deployment   kafka-amqstreams-consumer   0     5     kafka                       True    False    True       2m53s
```

7. Final cleanup:
```bash
oc delete jobs --field-selector status.successful=1 
oc delete -f scaledobject-fallback.yaml
oc delete -f deployment.yaml
```
