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
  kubectl delete all -l app=minio
  kubectl delete job --all
  kubectl delete pvc --all
  kubectl delete datavolume --all
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

s3() {
  kubectl apply -f - <<Y
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: minio
  labels:
    app: minio
spec:
  selector:
    matchLabels:
      app: minio
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: minio
    spec:
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: example-datavolume
      containers:
      - name: minio
        volumeMounts:
        - name: data
          mountPath: "/data"
        image: minio/minio
        args:
        - server
        - /data
        env:
        - name: MINIO_ACCESS_KEY
          value: "minio"
        - name: MINIO_SECRET_KEY
          value: "minio123"
        ports:
        - containerPort: 9000
---
apiVersion: v1
kind: Service
metadata:
  name: minio
  labels:
    app: minio
spec:
  type: LoadBalancer
  ports:
    - port: 9000
      targetPort: 9000
      protocol: TCP
  selector:
    app: minio
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: minio
  labels:
    app: minio
spec:
  rules:
  - host: lnx-d4025
    http:
      paths:
      - path: /minio
        backend:
          serviceName: minio
          servicePort: 9000
Y
}

"$@"
