package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/infrakit/spi/instance"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type plugin struct {
	Dir string
}

type schema struct {
	Memory          int
	Core            int
	DiskSize        int
	SourceArchiveID int64
	Zone            string
	Password        string
}

func NewSakuraCloudInstancePlugin(dir string) instance.Plugin {
	log.Println("sakuracloud instance plugin. dir=", dir)
	return &plugin{
		Dir: dir,
	}
}

func (v *plugin) Validate(req json.RawMessage) error {
	log.Printf("validate : %s", string(req))

	var token, secret string
	var ok bool

	if token, ok = os.LookupEnv("SAKURACLOUD_ACCESS_TOKEN"); !ok {
		return fmt.Errorf("Invalid environment values setting: %s", "Please set SAKURACLOUD_ACCESS_TOKEN environment value")
	}

	if secret, ok = os.LookupEnv("SAKURACLOUD_ACCESS_TOKEN_SECRET"); !ok {
		return fmt.Errorf("Invalid environment values setting: %s", "Please set SAKURACLOUD_ACCESS_TOKEN_SECRET environment value")
	}

	prop := &schema{
		Core:   1,
		Memory: 1,
		Zone:   "tk1a",
	}
	if err := json.Unmarshal([]byte(req), &prop); err != nil {
		return fmt.Errorf("Invalid instance properties: %s", err)
	}

	if prop.SourceArchiveID <= 0 {
		return fmt.Errorf("Invalid instance properties: %s", "SourceArchiveID is required")
	}
	if prop.Password == "" {
		return fmt.Errorf("Invalid instance properties: %s", "Password is required")
	}

	if prop.Zone != "" {
		if prop.Zone != "tk1a" && prop.Zone != "is1b" {
			return fmt.Errorf("Invalid instance properties: %s", "Zone is need in [is1b , tk1a]")
		}
	}

	client := api.NewClient(token, secret, prop.Zone)

	plan, err := client.Product.Server.GetBySpec(prop.Core, prop.Memory)
	if err != nil {
		return fmt.Errorf("Invalid instance properties: Core or Memory is invalid : %s", err)
	}

	if plan == nil {
		return fmt.Errorf("Invalid instance properties: %s", "Plan is not found. Please change Core or Memory")
	}

	return nil

}

func (v *plugin) Provision(spec instance.Spec) (*instance.ID, error) {

	properties := &schema{
		Memory:   1,
		Core:     1,
		DiskSize: 20,
		Zone:     "tk1a",
	}
	if spec.Properties != nil {
		if err := json.Unmarshal(*spec.Properties, properties); err != nil {
			return nil, fmt.Errorf("Invalid instance properties: %s", err)
		}
	}

	server, err := v.createServer(spec, properties)
	if err != nil {
		return nil, err
	}

	machineDir, err := ioutil.TempDir(v.Dir, "infrakit-")
	if err != nil {
		return nil, err
	}

	id := instance.ID(path.Base(machineDir))

	tagData, err := json.Marshal(spec.Tags)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(path.Join(machineDir, "tags"), tagData, 0666); err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(path.Join(machineDir, "ip"), []byte(server.Interfaces[0].IPAddress), 0666); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path.Join(machineDir, "id"), []byte(fmt.Sprintf("%d", server.ID)), 0666); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path.Join(machineDir, "zone"), []byte(properties.Zone), 0666); err != nil {
		return nil, err
	}

	return &id, nil

}

// Destroy terminates an existing instance.
func (v *plugin) Destroy(id instance.ID) error {
	fmt.Println("Destroying ", id)

	machineDir := path.Join(v.Dir, string(id))
	_, err := os.Stat(machineDir)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Instance does not exist")
		}
	}

	idData, err := ioutil.ReadFile(path.Join(machineDir, "id"))
	if err != nil {
		return fmt.Errorf("Failed on destroy server : %s", err)
	}

	zoneData, err := ioutil.ReadFile(path.Join(machineDir, "zone"))
	if err != nil {
		return fmt.Errorf("Failed on destroy server : %s", err)
	}

	intID, err := strconv.ParseInt(string(idData), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed on destroy server : %s", err)
	}

	token := os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	secret := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")

	client := api.NewClient(token, secret, string(zoneData))

	server, err := client.Server.Read(intID)
	if err != nil {
		return fmt.Errorf("Failed on destroy server : %s", err)
	}
	if server != nil {

		_, err := client.Server.DeleteWithDisk(server.ID, []int64{server.Disks[0].ID})
		if err != nil {
			return fmt.Errorf("Failed on destroy server : %s", err)
		}

	}

	return nil
}

// DescribeInstances returns descriptions of all instances matching all of the provided tags.
func (v *plugin) DescribeInstances(tags map[string]string) ([]instance.Description, error) {
	files, err := ioutil.ReadDir(v.Dir)
	if err != nil {
		return nil, err
	}

	descriptions := []instance.Description{}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		machineDir := path.Join(v.Dir, file.Name())

		tagData, err := ioutil.ReadFile(path.Join(machineDir, "tags"))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}

		machineTags := map[string]string{}
		if err := json.Unmarshal(tagData, &machineTags); err != nil {
			return nil, err
		}

		allMatched := true
		for k, v := range tags {
			value, exists := machineTags[k]
			if !exists || v != value {
				allMatched = false
				break
			}
		}

		if allMatched {
			var logicalID *instance.LogicalID
			ipData, err := ioutil.ReadFile(path.Join(machineDir, "ip"))
			if err == nil {
				id := instance.LogicalID(ipData)
				logicalID = &id
			} else {
				if !os.IsNotExist(err) {
					return nil, err
				}
			}

			descriptions = append(descriptions, instance.Description{
				ID:        instance.ID(file.Name()),
				LogicalID: logicalID,
				Tags:      machineTags,
			})
		}
	}

	return descriptions, nil
}

func (v *plugin) createServer(spec instance.Spec, prop *schema) (*sacloud.Server, error) {

	token := os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	secret := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")

	client := api.NewClient(token, secret, prop.Zone)

	// get plan
	plan, err := client.Product.Server.GetBySpec(prop.Core, prop.Memory)
	if err != nil {
		return nil, fmt.Errorf("Failed on creating server: %s", err)
	}

	// create init script
	var startupScriptID int64
	if spec.Init != "" {
		noteReq := client.Note.New()
		noteReq.Name = "infrakit-plugin-sakuracloud"
		noteReq.Content = spec.Init

		note, err := client.Note.Create(noteReq)
		if err != nil {
			return nil, fmt.Errorf("Failed on creating startup script: %s", err)
		}
		startupScriptID = note.ID

		defer func() {
			client.Note.Delete(startupScriptID)
		}()
	}

	// create disk
	diskReq := client.Disk.New()
	diskReq.Name = "infrakit-plugin-sakuracloud"
	diskReq.SizeMB = prop.DiskSize * 1024
	diskReq.SetSourceArchive(prop.SourceArchiveID)

	disk, err := client.Disk.Create(diskReq)
	if err != nil {
		return nil, fmt.Errorf("Failed on creating disk: %s", err)
	}

	err = client.Disk.SleepWhileCopying(disk.ID, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("Timeout on creating disk: %s", err)
	}

	// edit disk
	editReq := client.Disk.NewCondig()
	editReq.SetPassword(prop.Password)
	if startupScriptID > 0 {
		editReq.AddNote(fmt.Sprintf("%d", startupScriptID))
	}
	_, err = client.Disk.Config(disk.ID, editReq)
	if err != nil {
		return nil, fmt.Errorf("Failed on editing disk: %s", err)
	}

	// create server
	serverReq := client.Server.New()
	serverReq.SetServerPlanByID(fmt.Sprintf("%d", plan.ID))
	serverReq.Name = fmt.Sprintf("infrakit-plugin-sakuracloud%s", spec.LogicalID)

	serverReq.AppendTag("@virtio-net-pci")
	serverReq.AppendTag("infrakit-plugin-sakuracloud")
	serverReq.AddPublicNWConnectedParam()

	server, err := client.Server.Create(serverReq)
	if err != nil {
		return nil, fmt.Errorf("Failed on creating server: %s", err)
	}

	// connect disk
	_, err = client.Disk.ConnectToServer(disk.ID, server.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed on connect disk to server: %s", err)
	}

	_, err = client.Server.Boot(server.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed on booting server: %s", err)
	}

	err = client.Server.SleepUntilUp(server.ID, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("Failed on booting server: %s", err)
	}

	return server, nil
}
