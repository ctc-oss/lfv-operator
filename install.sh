#!/usr/bin/env bash

kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/crds/com.github.jw3_datavolumes_crd.yaml
kubectl apply -f deploy/operator.yaml
kubectl apply -f deploy/crds/com.github.jw3_v1alpha1_datavolume_cr.yaml
