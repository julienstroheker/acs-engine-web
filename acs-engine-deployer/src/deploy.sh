acs-engine-template-generator
acs-engine generate --output-directory template /tmp/template.json
az login --service-principal -u ${SERVICE_PRINCIPLE_NAME} -p ${SERVICE_PRINCIPLE_PWD} --tenant ${TENANT_ID}
az group create -l "${LOCATION}" -n ${RESOURCE_GROUP}
az group deployment create -g ${RESOURCE_GROUP} --template-file template/azuredeploy.json --parameters @template/azuredeploy.parameters.json