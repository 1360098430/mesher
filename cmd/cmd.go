package cmd

import (
	"fmt"
	chassiscommon "github.com/ServiceComb/go-chassis/core/common"
	"github.com/go-chassis/mesher/common"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

//Local is a constant
const Local = "127.0.0.1"

//ConfigFromCmd store cmd params
type ConfigFromCmd struct {
	ConfigFile        string
	Mode              string
	LocalServicePorts string
	PortsMap          map[string]string
}

//Configs is a pointer of struct ConfigFromCmd
var Configs *ConfigFromCmd

// parseConfigFromCmd
func parseConfigFromCmd(args []string) (err error) {
	app := cli.NewApp()
	app.HideVersion = true
	app.Usage = "Service mesh."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Usage:       "mesher config file, example: --config=mesher.yaml",
			Destination: &Configs.ConfigFile,
		},
		cli.StringFlag{
			Name:        "mode",
			Value:       common.ModeSidecar,
			Usage:       fmt.Sprintf("mesher running mode [ %s|%s ]", common.ModePerHost, common.ModeSidecar),
			Destination: &Configs.Mode,
		},
		cli.StringFlag{
			Name:        "service-ports",
			EnvVar:      common.EnvServicePorts,
			Usage:       fmt.Sprintf("service protocol and port,examples: --service-ports=http:3000,grpc:8000"),
			Destination: &Configs.LocalServicePorts,
		},
	}
	app.Action = func(c *cli.Context) error {
		return nil
	}

	err = app.Run(args)
	return
}

//Init get config and parses those command
func Init() error {
	Configs = &ConfigFromCmd{}
	return parseConfigFromCmd(os.Args)
}

//GeneratePortsMap generates ports map
func (c *ConfigFromCmd) GeneratePortsMap() error {
	c.PortsMap = make(map[string]string)
	if c.LocalServicePorts != "" { //parse service ports
		s := strings.Split(c.LocalServicePorts, ",")
		for _, v := range s {
			p := strings.Split(v, ":")
			if len(p) != 2 {
				return fmt.Errorf("[%s] is invalid", p)
			}
			c.PortsMap[p[0]] = Local + ":" + p[1]
		}
		return nil
	}
	//support deprecated env
	addr := os.Getenv(common.EnvSpecificAddr)
	if addr != "" {
		addr = strings.TrimSpace(addr)
		log.Printf("%s is deprecated, plz use SERVICE_PORTS=http:8080,grpc:90000 instead", common.EnvSpecificAddr)
		s := strings.Split(addr, ":")
		if len(s) != 2 {
			return fmt.Errorf("[%s] is invalid", addr)
		}
		c.PortsMap[chassiscommon.ProtocolRest] = Local + ":" + s[1]
	}

	return nil
}
