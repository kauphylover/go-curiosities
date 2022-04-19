package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type NexusDmProps struct {
	Type      string `yaml:"type"`
	Url       string `yaml:"url"`
	IsDefault bool   `yaml:"isDefault"`
}

const nexusDmsFile = "nexus-dms.yaml"

func WriteToNexusDms(dms map[string]NexusDmProps) error {
	_, err := os.Stat(nexusDmsFile)
	if err != nil {
		fmt.Printf("Couldn't find file %s, creating\n", nexusDmsFile)
		_, err = os.Create(nexusDmsFile)
		if err != nil {
			return fmt.Errorf("Couldn't create file %s\n", nexusDmsFile)
		}
	}
	data, err := ioutil.ReadFile(nexusDmsFile)
	if err != nil {
		return fmt.Errorf("Could not read %s\n", nexusDmsFile)
	}

	var nexusDmMap map[string]NexusDmProps
	err = yaml.Unmarshal(data, &nexusDmMap)
	if err != nil {
		return fmt.Errorf("could not unmarshal %s\n", nexusDmsFile)
	}

	if nexusDmMap == nil {
		nexusDmMap = make(map[string]NexusDmProps)
	}
	for k, v := range dms {
		nexusDmMap[k] = v
	}

	data, err = yaml.Marshal(&nexusDmMap)
	if err != nil {
		return fmt.Errorf("Error while Marshaling nexus-dms. %v", err)
	}

	err = ioutil.WriteFile(nexusDmsFile, data, 0644)
	if err != nil {
		return fmt.Errorf("Could not write to %s: %v\n", nexusDmsFile, err)
	}

	return nil
}

func main() {
	dms := make(map[string]NexusDmProps)

	dms["dm2"] = NexusDmProps{"remote", "gitlab.eng.vmware.com/amallela/x", true}

	err := WriteToNexusDms(dms)

	if err != nil {
		fmt.Errorf("ERROR")
	}
}
