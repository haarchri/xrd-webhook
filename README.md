# xrd-webhook
test implementation for:
https://github.com/crossplane/crossplane/pull/3940
https://github.com/crossplane/crossplane/pull/4310

## PreRequisits

Create the crossplane deployment with cert-manager

```
./setup.sh
```

## Deploy the Webhook

```
kubectl apply -f examples/webhook/webhook.yaml
```

## Simulate the upgrade process

### Create v1beta1
Create the initial v1beta1 `CompositeResourceDefinition`
```
kubectl apply -f examples/v1beta1/definition.yaml
```

Create the the v1beta1 `Claim`
```
kubectl apply -f examples/v1beta1/claim.yaml
```

Confirm you can geth the v1beta `Claim`
```
kubectl get testconversion.conversion.haarchri.io -A
NAMESPACE           NAME                    SYNCED   READY   CONNECTION-SECRET   AGE
crossplane-system   test-conversion-beta    True     False                       12s
crossplane-system   test-conversion-gamma   True     False                       12s

kubectl get --raw /apis/conversion.haarchri.io/v1beta1/namespaces/crossplane-system/testconversions/test-conversion-beta | jq 
[...]

kubectl get testconversion.conversion.haarchri.io -n crossplane-system   test-conversion-beta -o yaml
apiVersion: conversion.haarchri.io/v1beta1
kind: TestConversion
metadata:
  annotations:
  name: test-conversion-beta
  namespace: crossplane-system
spec:
  compositeDeletePolicy: Background
  hostPort: localhost:8080
  resourceRef:
    apiVersion: conversion.haarchri.io/v1beta1
    kind: XTestConversion
    name: test-conversion-beta-8s7dq
```

### bump to v1

Apply our upgrade with v1beta1 and v1 `CompositeResourceDefinition`
```
kubectl apply -f examples/v1/definition.yaml
```

Claim is stored as `v1beta1`, but it can now serve as `v1` and reflect as `v1` schema

```
kubectl get --raw /apis/conversion.haarchri.io/v1beta1/namespaces/crossplane-system/testconversions/test-conversion-beta | jq 
[...]
kubectl get --raw /apis/conversion.haarchri.io/v1/namespaces/crossplane-system/testconversions/test-conversion-beta | jq 
[...]

kubectl get testconversion.conversion.haarchri.io -n crossplane-system   test-conversion-beta -o yaml
apiVersion: conversion.haarchri.io/v1
kind: TestConversion
metadata:
  annotations:
  name: test-conversion-beta
  namespace: crossplane-system
spec:
  compositeDeletePolicy: Background
  host: localhost
  port: "8080"
  resourceRef:
    apiVersion: conversion.haarchri.io/v1
    kind: XTestConversion
    name: test-conversion-beta-llxss
```

Create a `v1` Claim that will be stored in etcd as `v1`
```
kubectl apply -f examples/v1/claim.yaml
```


**Note**: At this point etcd is storing all claims from /examples/v1beta1/claim.yaml (`v1beta1`) and all claims from /examples/v1/claim.yaml (`v1`). 
When you are ready, you can bump all existing stored `v1beta1` versions to `v1` by following the instructions
https://github.com/kubernetes-sigs/kube-storage-version-migrator


# Licensing
xrd-webhook is under the Apache 2.0 license.
The material in webhook/ in xrd-webhool is partially derived from https://github.com/madorn/crd-conversion-webhook