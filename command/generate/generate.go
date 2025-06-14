package generate

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/bigkevmcd/go-configparser"
	"github.com/mitchellh/cli"

	"github.com/blairham/aws-config/command/flags"
	"github.com/blairham/aws-config/provider/aws"
)

type cmd struct {
	UI    cli.Ui
	flags *flag.FlagSet
	help  string

	diff bool
}

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.Init()
	return c
}

func (c *cmd) Init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.BoolVar(&c.diff, "diff", false, "Enable diff output.")

	c.help = flags.Usage(help, c.flags)
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	configFile := aws.ConfigFile()

	cfg := aws.LoadDefaultConfig()
	token := aws.GetToken(cfg)

	// create sso client
	ssoClient := sso.NewFromConfig(cfg)

	if generateAwsConfigFile(ssoClient, token, configFile, c.diff) != nil {
		return 1
	}

	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return c.help
}

func showFileDiff(file1, file2 string) {
	// First check to make sure these are legit filenames so we don't confuse the "diff" command
	for _, file := range []string{file1, file2} {
		_, err := os.Stat(file)
		if err != nil {
			fmt.Printf("File %s does not exist\n", file)
			return
		}
	}
	cmd := exec.Command("diff", file1, file2)
	cmd.Stdout = os.Stdout
	cmd.Run() // Ignore error as diff returns non-zero when files differ
}

func generateAwsConfigFile(ssoClient *sso.Client, token *string, configFile string, diff bool) error {
	configFileNew := configFile + ".new"

	awsConfig, err := configparser.NewConfigParserFromFile(configFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Fetching list of all accounts for user")
	accountPaginator := sso.NewListAccountsPaginator(ssoClient, &sso.ListAccountsInput{
		AccessToken: token,
	})
	for accountPaginator.HasMorePages() {
		x, err := accountPaginator.NextPage(context.TODO())
		if err != nil {
			fmt.Println(err)
		}

		for _, y := range x.AccountList {
			// only add accounts that we care about
			accountName := aws.ToString(y.AccountName)
			if !strings.HasPrefix(accountName, "tripadvisor-") && !strings.HasPrefix(accountName, "trip-") && !strings.HasPrefix(accountName, "core-") {
				continue
			}

			trimmedAccountName := strings.TrimPrefix(accountName, "trip-")
			trimmedAccountName = strings.TrimPrefix(trimmedAccountName, "tripadvisor-")
			section := "profile " + trimmedAccountName

			// check if profile already exists and update it
			if !awsConfig.HasSection(section) {
				fmt.Printf("Adding profile %v\n", trimmedAccountName)
				awsConfig.AddSection(section)
			}

			awsConfig.Set(section, "sso_account_id", aws.ToString(y.AccountId))
			awsConfig.Set(section, "sso_role_name", "AdministratorAccess")
			awsConfig.Set(section, "sso_region", "us-east-1")
			awsConfig.Set(section, "sso_start_url", "https://tamg.awsapps.com/start")
			awsConfig.Set(section, "region", "us-east-1")
		}
	}
	err = awsConfig.SaveWithDelimiter(configFileNew, "=")
	if err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}
	if diff {
		showFileDiff(configFile, configFileNew)
	}
	err = os.Rename(configFileNew, configFile)
	if err != nil {
		return fmt.Errorf("failed to rename config file: %w", err)
	}

	return nil
}
