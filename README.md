example k8s operator
===

operator example demonstrating managing a PVC according to a CRD

https://github.com/operator-framework/operator-sdk#create-and-deploy-an-app-operator

### init
- `operator-sdk new example-operator --repo github.com/jw3/example-operator`
- `operator-sdk add api --api-version=github.com/jw3/example-operator --kind=DataVolume`
- `operator-sdk add controller --api-version=github.com/jw3/example-operator --kind=DataVolume`

### build
- `operator-sdk build jwiii/example-operator`
- `docker push jwiii/example-operator`

### run
- `./please.sh install`
- `./please.sh uninstall`

### addd-on
- `./please.sh s3`
