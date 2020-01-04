example k8s operator
===

simple operator example demonstrates managing a PVC according to a CRD

https://github.com/operator-framework/operator-sdk#create-and-deploy-an-app-operator


### dev
- `operator-sdk build jwiii/example-operator`
- `docker push jwiii/example-operator`
- `./please.sh install`