# Chap 12

## install using helm

```shell
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update


kubectl create namespace istio-system

helm install istio-base istio/base -n istio-system

helm install istiod istio/istiod -n istio-system --wait

kubectl create namespace istio-ingress
kubectl label namespace istio-ingress istio-injection=enabled
helm install istio-ingress istio/gateway -n istio-ingress --wait
```
remember to edit the loadbalancer type to nodePort to finish the installation

## install


## inject the istio enabled
```shell
k label ns default istio-injection=enabled

# remember to restart the previous deploy to get the sidecar injected
k rollout restart deploy cncamp
```

## get the label of the istio-ingress (be aware of that i am using default helm config, so the ingress label is different)
```shell
alias king="kubectl -n=istio-ingress"

king get po --show-labels
# NAME                             READY   STATUS    RESTARTS   AGE   LABELS
# istio-ingress-69495c6667-n5f8t   1/1     Running   0          23h   app=istio-ingress,istio.io/rev=default,istio=ingress,pod-template-hash=69495c6667,service.istio.io/canonical-name=istio-ingress,service.istio.io/canonical-revision=latest,sidecar.istio.io/inject=true

# get the nodeport
king get svc
# NAME            TYPE       CLUSTER-IP   EXTERNAL-IP   PORT(S)                                      AGE
# istio-ingress   NodePort   10.97.74.3   <none>        15021:30845/TCP,80:32582/TCP,443:32332/TCP   23h
```

## Apply the istio specs
```
k apply -f ./istiospec/
```

try to get the page from MacBook
```
curl -v -H "Host: cncamp.com" 192.168.34.2:32582    
```

```
*   Trying 192.168.34.2...
* TCP_NODELAY set
* Connected to 192.168.34.2 (192.168.34.2) port 32582 (#0)
> GET / HTTP/1.1
> Host: cncamp.com
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< accept: */*
< user-agent: curl/7.64.1
< version: 1.3.0
< x-b3-sampled: 0
< x-b3-spanid: 8aa250fa0a2b3aa0
< x-b3-traceid: ec797ae1653e128f8aa250fa0a2b3aa0
< x-envoy-attempt-count: 1
< x-envoy-internal: true
< x-forwarded-for: 10.0.2.15
< x-forwarded-proto: http
< x-request-id: 94194831-f62e-4d16-8edd-0ace68fa253e
< date: Sun, 19 Dec 2021 17:51:19 GMT
< content-length: 21
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 175
< server: istio-envoy
<
* Connection #0 to host 192.168.34.2 left intact
welcome to home page!* Closing connection 0
```