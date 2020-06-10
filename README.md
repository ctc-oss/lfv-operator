Git LFS Volume Operator (lfv-operator)
===

Operator responsible for managing PVCs according to a CRD storage specification.

### init
- `operator-sdk new example-operator --repo github.com/jw3/example-operator`
- `operator-sdk add api --api-version=github.com/jw3/example-operator --kind=DataVolume`
- `operator-sdk add controller --api-version=github.com/jw3/example-operator --kind=DataVolume`

### build
- `operator-sdk build ctcoss/example-operator`
- `docker push ctcoss/example-operator`

### run
- `./please.sh install`
- `./please.sh uninstall`

### addd-on
- `./please.sh s3`

### reference
- https://github.com/operator-framework/operator-sdk#create-and-deploy-an-app-operator