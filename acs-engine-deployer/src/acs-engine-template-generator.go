package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	vlabs "github.com/Azure/acs-engine/pkg/api/vlabs"
)

func main() {
	envvars := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = getval(splits)
		return
	})

	template := getTemplate(envvars)

	s, _ := json.MarshalIndent(template, "", "    ")
	writefile("/tmp/template.json", s)

	fmt.Println("Created acs-engine template:")
	fmt.Println("")
	fmt.Println(string(s))
	fmt.Println("")
}

func getTemplate(envvars map[string]string) Template {
	var template Template
	template.APIVersion = "vlabs"
	template.Properties = vlabs.Properties{}

	setOrchestratorprofile(&template.Properties, envvars)
	setMasterProfile(&template.Properties, envvars)
	setAgentPoolProfiles(&template.Properties, envvars)
	setLinuxProfiles(&template.Properties, envvars)

	return template
}

func setOrchestratorprofile(properties *vlabs.Properties, envvars map[string]string) {
	properties.OrchestratorProfile = &vlabs.OrchestratorProfile{}

	// var orchestratorprofile vlabs.OrchestratorProfile

	var orchestratortypevar = envvars["ORCHESTRATOR_TYPE"]
	switch orchestratortypevar {
	case "Kubernetes":
		properties.OrchestratorProfile.OrchestratorType = vlabs.Kubernetes
	case "DCOS":
		properties.OrchestratorProfile.OrchestratorType = vlabs.DCOS
	case "Swarm":
		properties.OrchestratorProfile.OrchestratorType = vlabs.Swarm
	case "SwarmMode":
		properties.OrchestratorProfile.OrchestratorType = vlabs.SwarmMode
	}

	properties.OrchestratorProfile.OrchestratorVersion = envvars["ORCHESTRATOR_VERSION"]

	return // orchestratorprofile
}

func setMasterProfile(properties *vlabs.Properties, envvars map[string]string) {
	properties.MasterProfile = &vlabs.MasterProfile{}

	i, _ := strconv.Atoi(envvars["MASTER_COUNT"])
	properties.MasterProfile.Count = i
	properties.MasterProfile.DNSPrefix = envvars["MASTER_DNS_PREFIX"]
	properties.MasterProfile.VMSize = envvars["MASTER_VM_SIZE"]
}

func setAgentPoolProfiles(properties *vlabs.Properties, envvars map[string]string) {
	properties.AgentPoolProfiles = []*vlabs.AgentPoolProfile{}

	publicAgentPoolProfile := getPublicAgentPoolProfile(envvars)
	properties.AgentPoolProfiles = append(properties.AgentPoolProfiles, &publicAgentPoolProfile)
}

func getPublicAgentPoolProfile(envvars map[string]string) vlabs.AgentPoolProfile {
	var publicAgentPoolProfile vlabs.AgentPoolProfile
	publicAgentPoolProfile.Name = "agentpublic"

	i, _ := strconv.Atoi(envvars["PUBLIC_POOL_COUNT"])
	publicAgentPoolProfile.Count = i
	publicAgentPoolProfile.DNSPrefix = envvars["PUBLIC_POOL_DNS_PREFIX"]
	publicAgentPoolProfile.VMSize = envvars["PUBLIC_POOL_VM_SIZE"]

	//TODO: ports
	return publicAgentPoolProfile
}

func setLinuxProfiles(properties *vlabs.Properties, envvars map[string]string) {
	properties.LinuxProfile = &vlabs.LinuxProfile{}

	properties.LinuxProfile.AdminUsername = envvars["ADMIN_USER_NAME"]

	properties.LinuxProfile.SSH.PublicKeys = []vlabs.PublicKey{}

	publicKey := removeQuotes(envvars["SSH_PUBLIC_KEY"])
	fmt.Println(publicKey)
	properties.LinuxProfile.SSH.PublicKeys = append(properties.LinuxProfile.SSH.PublicKeys, vlabs.PublicKey{
		KeyData: publicKey,
	})
}

func removeQuotes(input string) string {
	result := strings.TrimLeft(input, "\"")
	result = strings.TrimRight(result, "\"")
	return result
}

func writefile(filepath string, data []byte) {
	f, err := os.Create(filepath)
	check(err)
	defer f.Close()

	_, err = f.Write(data)
	check(err)

	f.Sync()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writebytes(data vlabs.OrchestratorProfile) {
	var out io.Writer
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
}

func getenvironment(data []string, getkeyval func(item string) (key, val string)) map[string]string {
	items := make(map[string]string)
	for _, item := range data {
		key, val := getkeyval(item)
		items[key] = val
	}
	return items
}

func getval(valarray []string) string {
	var val string
	for i := 1; i < len(valarray); i++ {
		if i > 1 {
			val += "="
		}

		val += valarray[i]
	}

	return val
}

//Containing struct for the templates being generated
type Template struct {
	APIVersion string           `json:"apiVersion,omitempty"`
	Properties vlabs.Properties `json:"properties"`
}
