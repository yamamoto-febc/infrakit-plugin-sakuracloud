package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/docker/infrakit/plugin/util"
	instance_plugin "github.com/docker/infrakit/spi/http/instance"
	"github.com/spf13/cobra"
	"log"
)

var (
	// PluginName is the name of the plugin in the Docker Hub / registry
	PluginName = "SakuraCloudInstance"

	// PluginType is the type / interface it supports
	PluginType = "infrakit.InstancePlugin/1.0"

	// Version is the build release identifier.
	Version = "Unspecified"

	// Revision is the build source control revision.
	Revision = "Unspecified"
)

func main() {

	if _, ok := os.LookupEnv("SAKURACLOUD_ACCESS_TOKEN"); !ok {
		log.Printf("Invalid environment values setting: %s", "Please set SAKURACLOUD_ACCESS_TOKEN environment value")
		os.Exit(1)
	}

	if _, ok := os.LookupEnv("SAKURACLOUD_ACCESS_TOKEN_SECRET"); !ok {
		log.Printf("Invalid environment values setting: %s", "Please set SAKURACLOUD_ACCESS_TOKEN_SECRET environment value")
		os.Exit(1)
	}

	listen := "unix:///run/infrakit/plugins/instance-sakuracloud.sock"
	dir := os.TempDir()

	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "SakuraCloud instance plugin",
		RunE: func(c *cobra.Command, args []string) error {

			if c.Use == "version" {
				return nil
			}

			log.Println("Starting plugin")
			log.Println("Listening on:", listen)

			_, stopped, err := util.StartServer(listen, instance_plugin.PluginServer(
				NewSakuraCloudInstancePlugin(dir)))

			if err != nil {
				log.Print(err)
			}

			<-stopped // block until done

			log.Println("Server stopped")
			return nil
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print build version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			buff, err := json.MarshalIndent(map[string]interface{}{
				"name":     PluginName,
				"type":     PluginType,
				"version":  Version,
				"revision": Revision,
			}, "  ", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(buff))
			return nil
		},
	})

	cmd.Flags().StringVar(&listen, "listen", listen, "listen address (unix or tcp) for the control endpoint")
	cmd.Flags().StringVar(&dir, "dir", dir, "Dir for storing the files")

	err := cmd.Execute()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
