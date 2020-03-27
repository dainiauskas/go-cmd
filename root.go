package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dainiauskas/flog"
	"github.com/dainiauskas/go-service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	banner = `
▓█████▄  ██▀███   ▄▄▄       ▄████▄   ▒█████  ▓█████▄ ▓█████
▒██▀ ██▌▓██ ▒ ██▒▒████▄    ▒██▀ ▀█  ▒██▒  ██▒▒██▀ ██▌▓█   ▀
░██   █▌▓██ ░▄█ ▒▒██  ▀█▄  ▒▓█    ▄ ▒██░  ██▒░██   █▌▒███
░▓█▄   ▌▒██▀▀█▄  ░██▄▄▄▄██ ▒▓▓▄ ▄██▒▒██   ██░░▓█▄   ▌▒▓█  ▄
░▒████▓ ░██▓ ▒██▒ ▓█   ▓██▒▒ ▓███▀ ░░ ████▓▒░░▒████▓ ░▒████▒
▒▒▓  ▒ ░ ▒▓ ░▒▓░ ▒▒   ▓▒█░░ ░▒ ▒  ░░ ▒░▒░▒░  ▒▒▓  ▒ ░░ ▒░ ░
░ ▒  ▒   ░▒ ░ ▒░  ▒   ▒▒ ░  ░  ▒     ░ ▒ ▒░  ░ ▒  ▒  ░ ░  ░
░ ░  ░   ░░   ░   ░   ▒   ░        ░ ░ ░ ▒   ░ ░  ░    ░
  ░       ░           ░  ░░ ░          ░ ░     ░       ░  ░
░                         ░                  ░
%s - %s [%s]
https://dracode.xyz

`
)

// CMD structore to work with cmd configuration and commands
type CMD struct {
	info    *Info
	root    *cobra.Command
	verbose bool
	config  string
	log     bool

	Configuration Configuration
}

// New create new root command
func New(name, desc, ver, build string) *CMD {
	name = strings.ToLower(name)

	rootCmd := &cobra.Command{
		Use:   name,
		Short: name,
		Long:  fmt.Sprintf(banner, strings.ToUpper(name), desc, ver),
	}

	var verbose bool
	config := fmt.Sprintf("%s.yaml", name)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&config, "config", config, "configuration file")

	return &CMD{
		root: rootCmd,
		info: &Info{
			Name:        name,
			Description: desc,
			Version:     ver,
			Build:       build,
		},
		verbose: verbose,
		config:  config,
	}
}

// LoadConfig used to read and load configuration file
func (c *CMD) LoadConfig(cfg Configuration) error {
	if err := c.configExists(); err != nil {
		return err
	}

	if err := c.loadConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	c.Configuration = cfg

	return nil
}

// Execute used to start using cmd executable
func (c *CMD) Execute() error {
	if c.Configuration == nil {
		return errors.New("Not found configuration")
	}

	c.initLog()

	return c.root.Execute()
}

// AddVersion added to show version command
func (c *CMD) AddVersion() {
	c.root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(c.info.Name, c.info.Version)
		},
	})
}

// AddService used to install and control services
func (c *CMD) AddService(cb func()) {
	service := &cobra.Command{
		Use:   "service",
		Short: "Aplication Service Control",
		Args:  cobra.RangeArgs(1, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			srv := service.New(c.info.Name, c.info.Name, c.info.Description)
			srv.SetCb(cb)

			return srv.Controller(args[0])
		},
	}

	c.root.AddCommand(service)
}

func (c *CMD) configExists() error {
	if _, err := os.Stat(c.config); os.IsNotExist(err) {
		app, err := os.Executable()
		if err != nil {
			return err
		}

		if err := os.Chdir(filepath.Dir(app)); err != nil {
			return err
		}
	}

	return nil
}

func (c *CMD) loadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigFile(c.config)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func (c *CMD) initLog() {
	flog.SetLogTrace(c.verbose)
	flog.SetLogPath(fmt.Sprintf("./log_%s/", c.info.Name))
}
