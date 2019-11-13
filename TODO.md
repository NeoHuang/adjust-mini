# Todo
- add kafka, aero, redis, psql
- use config map
- PVC and stateful for csv_collector/csv_uploader

# Done
- helm restart (no need)
- rolling update readness liveness check
- add prometheus

# Notes
curl in mac doesn't work well with the apiserver keys
generate p12 first

openssl pkcs12 -export -in ~/.minikube/apiserver.crt -inkey ~/.minikube/apiserver.key -out apiserver.p12

then
curl  -k -H --cacert ~/.minikube/ca.cert -E apiserver.p12:i23456 "https://192.168.99.100:8443/api/v1/nodes"


## Prometheus Operator
helm install stable/prometheus-operator  --name=monitoring --namespace=monitoring
check match rules for selecting service monitor
```
kubectl get prometheus -o yaml -n monitoring
```
create service monitor to match the rules
service_monitor.selector (matches->) service (matches->) pods

check prometheus
kubectl port-forward -n monitoring prometheus-monitoring-prometheus-oper-prometheus-0 9090





