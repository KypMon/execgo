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
k apply -f ./istiospec/ingress.yaml
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

## Generate the Cert
```shell
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.com' -keyout cncamp.com.key -out cncamp.com.crt

kubectl create -n istio-ingress secret tls cncamp-istio-credential --key=cncamp.com.key --cert=cncamp.com.crt
```

## apply the new ssl enabled ingress spec
```
k apply -f ./istiospec/ingress-with-cert.yml

curl --resolve cncamp.com:443:10.97.74.3 https://cncamp.com/healthz -v -k
```

```
* Added cncamp.com:443:10.97.74.3 to DNS cache
* Hostname cncamp.com was found in DNS cache
*   Trying 10.97.74.3:443...
* TCP_NODELAY set
* Connected to cncamp.com (10.97.74.3) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/certs/ca-certificates.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_256_GCM_SHA384
* ALPN, server accepted to use h2
* Server certificate:
*  subject: O=cncamp Inc.; CN=*.cncamp.com
*  start date: Dec 20 08:48:24 2021 GMT
*  expire date: Dec 20 08:48:24 2022 GMT
*  issuer: O=cncamp Inc.; CN=*.cncamp.com
*  SSL certificate verify result: self signed certificate (18), continuing anyway.
* Using HTTP2, server supports multi-use
* Connection state changed (HTTP/2 confirmed)
* Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
* Using Stream ID: 1 (easy handle 0x555b52d6de10)
> GET /healthz HTTP/2
> Host: cncamp.com
> user-agent: curl/7.68.0
> accept: */*
>
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* old SSL session ID is stale, removing
* Connection state changed (MAX_CONCURRENT_STREAMS == 2147483647)!
< HTTP/2 200
< accept: */*
< user-agent: curl/7.68.0
< version: 1.3.0
< x-b3-sampled: 0
< x-b3-spanid: b3ad3e72c40ca1b2
< x-b3-traceid: 6fb0254853bf4cf8b3ad3e72c40ca1b2
< x-envoy-attempt-count: 1
< x-envoy-internal: true
< x-forwarded-for: 10.0.2.15
< x-forwarded-proto: https
< x-request-id: e3ce2fa3-372f-4dd1-a97b-da58bf0534d5
< date: Mon, 20 Dec 2021 09:24:46 GMT
< content-length: 19
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 5
< server: istio-envoy
<
* Connection #0 to host cncamp.com left intact
welcome to healthz!
```

## HTTP based canary 
```shell
# deploy v2 app and v1 app with label
k apply -f ./k8sspec/spec.yaml
k apply -f ./k8sspec/spec-v2.yaml

# add canary virtual service
k apply -f ./istiospec/advance-route/virtualService-routing.yml

# start a common debug pod
k run doks-debug --image=digitalocean/doks-debug --rm -it --restart=Never bash

# inside of doks-debug
curl cncamp

# > welcome to home page

curl cncamp -H "user: cncamp"

# > welcome to home page V2!
```

## Telemetry
```shell
k apply -f ./istiospec/jaeger.yml
k edit configmap istio -n istio-system
# and set tracing.sampling=100

ki edit svc tracing
# from ClusterIP to NodePort

# in browser go http://192.168.34.2:31303/
```

# Jaeger Image here
