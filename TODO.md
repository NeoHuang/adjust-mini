todo
- helm restart (no need)
- rolling update readness liveness check
- add prometheus
- add kafka, aero, redis, psql
- use config map

done

notes
curl in mac doesn't work well with the apiserver keys
generate p12 first

openssl pkcs12 -export -in ~/.minikube/apiserver.crt -inkey ~/.minikube/apiserver.key -out apiserver.p12

then
curl  -k -H --cacert ~/.minikube/ca.cert -E apiserver.p12:i23456 "https://192.168.99.100:8443/api/v1/nodes"
