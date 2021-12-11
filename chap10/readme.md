# CHAP 08

## deploy nginx ingress controller
Replace all the k8s.gcr.io/ingress-nginx image to docker hub
```
# FROM 
image: k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0@sha256:f3b6b39a6062328c095337b4cadcefd1612348fdd5190b1dcbcb9b9e90bd8068
# TO
image: liangjw/kube-webhook-certgen:v1.0

# FROM
image: k8s.gcr.io/ingress-nginx/controller:v1.0.0@sha256:0851b34f69f69352bf168e6ccf30e1e20714a264ab1ecd1933e4d8c0fc3215c6
# TO
image: liangjw/ingress-nginx-controller:v1.0.0

# FROM
image: k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0@sha256:f3b6b39a6062328c095337b4cadcefd1612348fdd5190b1dcbcb9b9e90bd8068
# TO
image: liangjw/kube-webhook-certgen:v1.0
```

Deploy NGINX Ingress controller
```
kubectl apply -f ./nginx-ingress-deployment.yaml
```

## Deploy Cert Manager
```
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.4/cert-manager.yaml
kubectl apply -f ./cert-manager/clusterissuer.yaml
```

## Deploy App and Service
```
kubectl apply -f ./spec.yaml
```

## Deploy Ingress
```
kubectl apply -f ./ingress.yaml
```

## check if the ingress is up and ready
```
kubectl get ing 
```
should have an IP address in Address column

## get the IP for the host
```
k get svc ingress-nginx-controller  -n ingress-nginx
```
returns 
```
NAME                       TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller   NodePort   10.107.78.92   <none>        80:32510/TCP,443:31154/TCP   123m
```

## send the request to the host
```
curl -H "Host: cncamp.com" https://10.107.78.92/healthz -v -k
```

returns 

```
welcome to healthz!
```

# Chap 10

```shell
k apply -f spec.yaml
```

install loki-prometheus stack: https://github.com/cncamp/101/tree/master/module10/loki-stack
```shell
helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```



```shell
k get secrets loki-grafana  -o jsonpath="{.data.admin-user}" | base64 -d

k get secrets loki-grafana  -o jsonpath="{.data.admin-password}" | base64 -d
```

import the json to the grafana dashboard: https://github.com/cncamp/101/raw/master/module10/httpserver/grafana-dashboard/httpserver-latency.json