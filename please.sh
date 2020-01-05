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
  kubectl delete job --all
  kubectl delete pvc --all
  kubectl delete ${operator}
}

show() {
  kubectl get deploy
  kubectl get pods
  kubectl get jobs
  kubectl get datavolume
  kubectl get pvc
}

watch() {
  /usr/bin/watch ${0} show
}

lsvol() {
  kubectl apply -f - <<Y
kind: Pod
apiVersion: v1
metadata:
  name: vdb
spec:
  volumes:
    - name: v
      persistentVolumeClaim:
        claimName: example-datavolume
  containers:
    - name: debugger
      image: busybox
      command: ['sleep', '3600']
      volumeMounts:
        - mountPath: "/data"
          name: v
Y
  sleep 3
  kubectl exec vdb -it -- ls /data
  kubectl delete pod vdb
}

"$@"
