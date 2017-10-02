# Acs Engine Web Architecture

## UI

[POC](https://acsengweb.azureedge.net/index.html)

This UI will send mandatory variables to API such as :
* apiVersion
* Orchestrator
* Number of master
* Vm Size for the master
* DNS prefix
* Azure credentials such as ClientId, TenantID, Subscription ID, etc...
* ...

> Note : The Json on the right will be removed

## API 

The api will construct the JSON payload mandatory for ACS-Engine to build the ARM template

When the payload is constructed it will send it to the runner.

## The Runner

This layer will :
* Generate the ARM template thanks to acs-engine
* Publish the crendentials informations (Certificates, ARM template, KubeConfig, SSHKey, etc...) to a storage account in the client subscription
* Deploy the cluster in the client subscription thanks to the AZ CLI
* When the deployment is done it will send back to the UI more information such as the FQDN, location of the storage account with the crendentials infos...