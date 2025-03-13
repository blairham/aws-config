package run

import (
	"log"
	"os"
	"os/exec"

	"github.com/blairham/aws-config/command/flags"
	"github.com/blairham/aws-config/provider/aws"
	"github.com/mitchellh/cli"
)

type cmd struct {
	UI   cli.Ui
	help string
}

var Logger *log.Logger

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
	awsProfile, err := aws.GetProfile()
	if err != nil {
		Logger.Fatalf("Could not determine AWS account: %s", err)
	}
	if awsProfile != "" {
		os.Setenv(AwsProfile, awsProfile)
	}

	cmd := exec.Command("aws2-wrap", args...)
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
