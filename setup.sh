#!/usr/bin/env bash
set -aeuo pipefail

echo "Running setup.sh"

kind create cluster --name=xrd-webhook
kubectx kind-xrd-webhook
kubectl create ns crossplane-system

helm install crossplane --namespace crossplane-system crossplane-stable/crossplane --version 1.13.1 --wait

echo "Waiting for all pods to come online..."
kubectl -n crossplane-system wait --for=condition=Available deployment --all --timeout=5m

echo "Creating provider..."
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-nop
spec:
  package: xpkg.upbound.io/upbound/upjet-provider-template:v0.0.0-159.gf689969
EOF

echo "Creating cloud credential secret..."
kubectl -n crossplane-system create secret generic provider-secret --from-literal=credentials="{\"username\": \"admin\"}" --dry-run=client -o yaml | kubectl apply -f -

echo "Waiting until provider is healthy..."
kubectl wait provider.pkg --all --for condition=Healthy --timeout 5m

echo "Waiting for all pods to come online..."
kubectl -n crossplane-system wait --for=condition=Available deployment --all --timeout=5m

echo "Creating a default provider config..."
cat <<EOF | kubectl apply -f -
apiVersion: template.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: provider-secret
      namespace: crossplane-system
      key: credentials
EOF

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml

./examples/webhook/tls.sh
