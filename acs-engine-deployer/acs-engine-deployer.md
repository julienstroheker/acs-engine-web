# Acs Engine ACI Container POC
The purpose of this POC is to attempt to encapsulate and simplify the logic of deploying a custom DC/OS, Docker or Kubernetes cluster.  This approach decouples the UI logic that gathers the cluster choices from the logic to deploy the cluster, encapsulating the latter in a Docker container.

![Architecture](https://bagbyimages.blob.core.windows.net/gitimages/Approach1.jpg "Architecture")

## UI
The role of the UI is to gather the cluster choices of the user, as well as to gather the Service Principal required to deploy the cluster.
The UI choices could be a web app, command line or an ARM template.  A web UI would encapsulate the choice logic.  In the case of an ARM template, the UI logic would be in the createUIDefinition.json.  Obviously, there would be no UI logic in the command line client.

## Docker Image
The image is relatively simple.  At a high level, it installs three applications, acs-engine-template-generator, acs-engin and az cli and sets an entrypoint to a shell script that runs them.  The following is deploy.sh:

```
acs-engine-template-generator
acs-engine generate --output-directory template /tmp/template.json
az login --service-principal -u ${SERVICE_PRINCIPLE_NAME} -p ${SERVICE_PRINCIPLE_PWD} --tenant ${TENANT_ID}
az group create -l "${LOCATION}" -n ${RESOURCE_GROUP}
az group deployment create -g ${RESOURCE_GROUP} --template-file template/azuredeploy.json --parameters @template/azuredeploy.parameters.json
```
### acs-engine-template-generator
acs-engine-template-generator is a simple, single file go application that, using the acs-engine object model, creates an ARM template, given the ENV variables passed in.

### acs-engine
acs-engine is installed in the Dockerfile.

### az cli
The az cli is installed in the Dockerfile and is used to deploy the generated ARM template.

## Deploying
### To Run the Docker Container via the command prompt
You will need to fill in the variables in the provided env.conf file first.
```
docker run --name test --env-file env.conf -ti rbagby/acs-engine-deployer /bin/bash
```

### To Run via an Azure Container Instance via the command prompt
```
az login
az group create -l "West US" -n aciresourcegroupname
az container create --name aciname \
  --image rbagby/acs-engine-deployer \
  --cpu 1 \
  --memory 1 \
  --ip-address public \
  -g aciresourcegroupname 
  -e ORCHESTRATOR_VERSION=1.10 \
  ORCHESTRATOR_TYPE=DCOS \
  MASTER_COUNT=3 \
  MASTER_DNS_PREFIX=bagbymaster \
  MASTER_VM_SIZE=Standard_D2_v2 \
  PUBLIC_POOL_COUNT=3 \
  PUBLIC_POOL_DNS_PREFIX=bagbyagent \
  PUBLIC_POOL_VM_SIZE=Standard_D2_v2 \
  ADMIN_USER_NAME=azureuser \
  SSH_PUBLIC_KEY="yoursshkeyhere" \
  SERVICE_PRINCIPLE_NAME=yourserviceprincipal \
  SERVICE_PRINCIPLE_PWD=yourserviceprincipalpassword \
  TENANT_ID=yourtennantid \
  RESOURCE_GROUP=clusterresourcegroup \
  LOCATION="West US"

az container show --name aciname --resource-group aciresourcegroupname --query state
az container logs --name aciname --resource-group aciresourcegroupname 
```

### Current state
* Currently, the code in this POC only deploys a DC/OS cluster with a fixed set of choices
* The requirement to create and supply the following is still challenging:
  * Service Principal information
  * Subscription and tenant information
* The POC currently does not provide details back to the caller such as FQDN
* Currently, when deployed via ACI,the deployment succeeds, but ACI detects an error and retrys.  This is under investigation.

### Possibilities
* It is possible that acs-engine-template-generator logic could be directly part of acs-engine
