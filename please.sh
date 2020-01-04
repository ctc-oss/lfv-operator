#!/usr/bin/env bash

readonly operator="deploy/example-operator"

install() {
  kubectl apply -f deploy/service_account.yaml
  kubectl apply -f deploy/role.yaml
  kubectl apply -f deploy/role_binding.yaml
  kubectl apply -f deploy/crds/com.github.jw3_datavolumes_crd.yaml
  kubectl apply -f deploy/operator.yaml
  kubectl apply -f deploy/crds/com.github.jw3_v1alpha1_datavolume_cr.yaml
}

uninstall() {
  kubectl delete datavolume --all
  kubectl delete pod example-datavolume-cloner
  kubectl delete pvc --all
  kubectl delete ${operator}
}

show() {
  kubectl get deploy
  kubectl get datavolume
  kubectl get pvc
}

watch() {
  /usr/bin/watch ${0} show
}

"$@"
