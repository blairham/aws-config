package run

import (
	"log"
	"os"
	"os/exec"

	"github.com/mitchellh/cli"

	"github.com/blairham/aws-sso-config/command/flags"
	"github.com/blairham/aws-sso-config/provider/aws"
)

type cmd struct {
	UI   cli.Ui
	help string
}

var Logger *log.Logger

// Variables to facilitate testing by allowing us to mock external dependencies
var (
	execCommand    = exec.Command
	getProfileFunc = aws.GetProfile
)

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.Init()
	return c
}

func (c *cmd) Init() {
	Logger = log.New(os.Stderr, "", 0)
	c.help = help
}

func (c *cmd) Run(args []string) int {
	awsProfile, err := getProfileFunc()
	if err != nil {
		Logger.Fatalf("Could not determine AWS account: %s", err)
	}
	if awsProfile != "" {
		err = os.Setenv(AwsProfile, awsProfile)
		if err != nil {
			Logger.Fatalf("Could not set environment variable: %s", err)
		}
	}

	cmd := execCommand("aws2-wrap", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		Logger.Fatal(err)
	}
	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return flags.Usage(c.help, nil)
}
