package util

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
)

var ClientID string

func persistMachineID(confDir string, machineID []byte) {
	log.Printf("persisting machine id %v", machineID)
	err := ioutil.WriteFile(confDir+"/.machine_id", machineID, 0600)
	if err != nil {
		panic(fmt.Sprintf("Can't persist machine id %v", err))
	}
}

func generateUUID() []byte {
	uuid, _ := uuid.NewUUID()
	return []byte(uuid.String())
}

func generateMachineID(confDir string) []byte {
	machineID := []byte{}
	if _, err := os.Stat("/etc/machine-id"); os.IsNotExist(err) {
		machineID = generateUUID()
	} else {
		machineID, err = ioutil.ReadFile("/etc/machine-id")
		if err != nil {
			panic(fmt.Sprintf("Failed to get machine id %v", err.Error()))
		}
	}
	persistMachineID(confDir, machineID)
	return machineID
}

// Read unique machine ID from `cfg.ConfigPath/.machine_id'
// if not exist, use /etc/machine-id if possible
// otherwise generate a new UUID
// and write into `cfg.ConfigPath/.machine_id'
func GenerateClientID(confDir string) {
	machineID, err := ioutil.ReadFile(confDir + "/.machine_id")
	if err != nil {
		fmt.Println("error reading machine id")
		machineID = generateMachineID(confDir)
	}
	ClientID = string(machineID[:len(machineID)-1])
	fmt.Println("generated ClientID", ClientID)
}

func GetClientID() string {
	return ClientID
}
