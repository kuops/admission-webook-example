# admission-webook-example

This is simple MutatingAdmissionWebhook example, 
auto injects to pod Env :

```yaml
env:
  name: TZ
  value: Asia/Shanghai
```

## Usage 

push image

```bash
make build-image
make push-image
```

Create namespace:

```bash
kubectl create ns sidecar-injector
```


Create a signed cert/key pair and store it in a Kubernetes secret that will be consumed by sidecar injector deployment:

```bash
./deployment/webhook-create-signed-cert.sh \
    --service sidecar-injector-webhook-svc \
    --secret sidecar-injector-webhook-certs \
    --namespace sidecar-injector
```

Patch the MutatingWebhookConfiguration by set caBundle with correct value from Kubernetes cluster:

```bash
cat deployment/mutatingwebhook.yaml | \
    deployment/webhook-patch-ca-bundle.sh > \
    deployment/mutatingwebhook-ca-bundle.yaml
```

deploy:

```bash
kubectl create -f deployment/deployment.yaml
kubectl create -f deployment/service.yaml
kubectl create -f deployment/mutatingwebhook-ca-bundle.yaml
```

Get sidecar inject webhook state:

```bash
# kubectl get po  -n sidecar-injector
NAME                                                  READY   STATUS        RESTARTS   AGE
sidecar-injector-webhook-deployment-847c47cbb-ktns5   1/1     Running       0          6s
```

Deploy nginx test it :

```bash
kubectl label namespace default  sidecar-injection=enabled
kubectl create deployment my-dep --image=nginx --replicas=1
```

Get The Env field:

```bash
# kubectl get po -l app=my-dep -o=jsonpath='{.items[0].spec.containers[0].env}'
[{"name":"CLUSTER_NAME","value":"aks-test-01"},{"name":"TZ","value":"Asia/Shanghai"}]%  
```

## references

morvencao/kube-mutating-webhook-tutorial: https://github.com/morvencao/kube-mutating-webhook-tutorial

