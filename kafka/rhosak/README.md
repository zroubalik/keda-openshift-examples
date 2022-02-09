# Autoscaling of Kafka Consumer application connected to Red Hat OpenShift Streams for Apache Kafka
The following guide describes the way how can be Kafka Consumer application autoscaled by KEDA on Openshift. The application is being scaled based on lag in the Kafka topic. If there isn't any traffic the application is autoscaled to 0 replicas, if there's some load the number of replicas is being scaled up to 5 replicas.

Kafka KEDA scaler is being used for this setup, for details please refer to [documentation](https://keda.sh/docs/latest/scalers/apache-kafka/).

## 0. Install KEDA
In `OperatorHub` locate and install KEDA, follow the instuctions to create `KedaController` instance in `keda` namespace.

## 1. Prepare Kafka Instance
 1. Create a Kafka instance at [Red Hat Hybrid Cloud Console](https://console.redhat.com/application-services/streams/kafkas) and obtain Boostrap Server address.
 2. Create new Service Account or reuse existing one to assing permissions to access Kafka Topic `my-topic` and Consumer Group `my-group` for the created Kafka instance

## 2. Create Secret with Kafka credentials
There are two ways how to create a secret:
 1. You can create a new secret with the following command, just replace `<Client ID>` and `<Client Secret>` with credentials assigned to the authorized Service Account.
 ```bash
oc create secret generic keda-kafka-secrets --from-literal=sasl='plaintext' --from-literal=tls='enable' --from-literal=username='<Client ID>' --from-literal=password='<Client Secret>' 
 ```
 2. Or you can update [secret.yaml](secret.yaml) file, specifying `username` and `password` properties. Use base64 vaules of `<Client ID>` and `<Client Secret>` credentials assigned to authorized the Service Account. Then create this secret with the following command:
 ```bash
 oc apply -f secret.yaml
 ```

## 3. Deploy Kafka Consumer application
Update `BOOTSTRAP_SERVERS` environment variable in [deployment.yaml](deployment.yaml) file to point to the `Bootstrap server` of the Kafka instance created earlier. Then deploy this application with the following command:
 ```bash
oc apply -f deployment.yaml
 ```
Verify the consumer has been able to connect to Kafka instance, run following command:
 ```bash
oc logs deployment.apps/kafka-consumer
 ```
You should see similar output:
 ```bash
2022/02/09 13:41:50 Go consumer starting with config=&{BootstrapServers:kafka-xxxxxx.kafka.rhcloud.com:443 Topic:my-topic GroupID:my-group SaslEnabled:true SaslUser:srvc-acct-xxxxxx SaslPassword:xxxxxx}
2022/02/09 13:41:56 Consumer group handler setup
2022/02/09 13:41:56 Sarama consumer up and running!...
 ```

## 4. Send messages to Kafka to test the Kafka Consumer application
Update `BOOTSTRAP_SERVERS` environment variable in [load.yaml](load.yaml) file to point to the `Bootstrap server` of the Kafka instance created earlier. Then create this Kubernetes Job with the following command:
```bash
oc create -f load.yaml
```
Chech the logs of the Kafka Consumer application:
```bash
oc logs deployment.apps/kafka-consumer
```
There should be 10 messages generated and consumed:
```bash
2022/02/09 14:15:41 Go consumer starting with config=&{BootstrapServers:kafka-xxxxxx.kafka.rhcloud.com:443 Topic:my-topic GroupID:my-group SaslEnabled:true SaslUser:xxxxxx SaslPassword:xxxxxx}
2022/02/09 14:15:47 Consumer group handler setup
2022/02/09 14:15:47 Sarama consumer up and running!...
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-1, topic=my-topic, partition=7, offset=120
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-0, topic=my-topic, partition=2, offset=116
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-2, topic=my-topic, partition=2, offset=117
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-3, topic=my-topic, partition=2, offset=118
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-5, topic=my-topic, partition=8, offset=111
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-4, topic=my-topic, partition=1, offset=142
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-6, topic=my-topic, partition=9, offset=127
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-7, topic=my-topic, partition=7, offset=121
2022/02/09 14:15:48 Message received: value=Hello from Go Kafka Sarama-8, topic=my-topic, partition=4, offset=130
2022/02/09 14:15:49 Message received: value=Hello from Go Kafka Sarama-9, topic=my-topic, partition=9, offset=128
```

## 5. Deploy ScaledObject to enable Kafka Consumer application autoscaling
Update `bootstrapServers` property in the trigger metadata in [scaledobject.yaml](scaledobject.yaml) file to point to the `Bootstrap server` of the Kafka instance created earlier. Then deploy this ScaledObject with the following command:
```bash
oc apply -f scaledobject.yaml
```
Check that KEDA has been able to access metrics and is correctly defined for autoscaling:
```bash
oc get scaledobject kafka-consumer-scaledobject
```
You should see similar output, `READY` should be `True`:
```bash
NAME                          SCALETARGETKIND      SCALETARGETNAME   MIN   MAX   TRIGGERS     AUTHENTICATION                      READY   ACTIVE   FALLBACK   AGE
kafka-consumer-scaledobject   apps/v1.Deployment   kafka-consumer    0     5     kafka        keda-trigger-auth-kafka-credential  True    False    False      1m10s
```
Because there aren't any messages in the Kafka topic, the Kafka Consumer application should be autoscaled to zero, run the following command:
```bash
oc get deployment.apps/kafka-consumer
```
You should see a similar output, `kafka-consumer` has been autoscaled to 0 replicas:
```bash
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
kafka-consumer   0/0     0            0           11m
```

## 6. Send more messages to Kafka to test the Kafka Consumer application autoscaling
Update `MESSAGE_COUNT` environment variable in [load.yaml](load.yaml) file, increase the value from `10` to at least `500` to generate more load. Then create this Kubernetes Job with the following command:
```bash
oc create -f load.yaml
```
You should see created an increased nubmer replicas of the Kafka Consumer application until all sent messages are processed. And the the application will be again autoscaled down to zero. You can check the changing number of replicas by running the following command:
```bash
watch oc get deployment.apps/kafka-consumer
```

The output should be similar:
```bash
Every 2,0s: oc get deployment.apps/kafka-consumer
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
kafka-consumer   5/5     5            5           21m

### After some time the application should be autoscaled back to 0

Every 2,0s: oc get deployment.apps/kafka-consumer
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
kafka-consumer   0/0     0            0           23m
```

## 7. Clean up
Run the following commands to remove all resources created in the namespace
```bash
oc delete jobs --field-selector status.successful=1 
oc delete -f scaledobject.yaml
oc delete -f deployment.yaml
oc delete -f secret.yaml
``` 
