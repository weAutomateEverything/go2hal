kubectl create configmap hal-mongo --from-file=hal-mongo.properties
kubectl create secret generic hal-secrets --from-file=hal-secrets.properties
