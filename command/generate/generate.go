package generate

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/bigkevmcd/go-configparser"
	"github.com/blairham/aws-config/command/flags"
	"github.com/blairham/aws-config/provider/aws"
	"github.com/mitchellh/cli"
)

type cmd struct {
	UI    cli.Ui
	flags *flag.FlagSet
	help  string

	diff    bool
	cleanup bool
}

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.Init()
	return c
}

func (c *cmd) Init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.BoolVar(&c.diff, "diff", false, "Enable diff output.")
	c.flags.BoolVar(&c.cleanup, "cleanup", false, "Remove profiles for deleted accounts.")

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

	if generateAwsConfigFile(ssoClient, token, configFile, c.diff, c.cleanup) != nil {
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
	cmd.Run()
}

func getDeletedAccounts() []string {
	// TODO: Implement

	// grep '^account_id' aws/conf/deleted/*.yaml | awk '{print $NF}' | tr "'" '"' | tr '\n' ',' | sed 's/,/, /g'
	return []string{"682663373751", "214402454761", "888141291843", "645449977376", "238748435140", "838966304137", "470843223423", "488930975526", "873418971637", "871576525389", "315408011559", "542652408899", "500107894351", "372722424469", "482829957692", "913664604572", "308372660733", "575631896363", "587273070693", "646305171068", "916698063943", "094139108652", "880843368113", "609658169657", "652580916927", "305349583326", "055450967330", "041793920579", "768492706723", "593319960666", "986061530938", "662733315436"}
}

func generateAwsConfigFile(ssoClient *sso.Client, token *string, configFile string, diff, cleanup bool) error {
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
			if !(strings.HasPrefix(aws.ToString(y.AccountName), "tripadvisor-") || strings.HasPrefix(aws.ToString(y.AccountName), "trip-") || strings.HasPrefix(aws.ToString(y.AccountName), "core-")) {
				continue
			}

			trimmedAccountName := strings.TrimPrefix(aws.ToString(y.AccountName), "trip-")
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
	if cleanup {
		deletedAccounts := getDeletedAccounts()
		for _, section := range awsConfig.Sections() {
			accountID, err := awsConfig.Get(section, "sso_account_id")
			if err != nil {
				continue
			}
			if slices.Contains(deletedAccounts, accountID) {
				awsConfig.RemoveSection(section)
			}
		}
	}
	awsConfig.SaveWithDelimiter(configFileNew, "=")
	if diff {
		showFileDiff(configFile, configFileNew)
	}
	os.Rename(configFileNew, configFile)

	return nil
}
