package aws

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/browser"
)

type SSOCacheEntry struct {
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

func ToString(p *string) string {
	return aws.ToString(p)
}

func ConfigFile() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return home + "/.aws/config"
}

func LoadDefaultConfig() aws.Config {
	// load default aws config
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal("error: ", err)
	}

	return cfg
}

func generateToken(cfg aws.Config) *string {
	// create sso oidc client to trigger login flow
	ssooidcClient := ssooidc.NewFromConfig(cfg)

	// register your client which is triggering the login flow
	register, err := ssooidcClient.RegisterClient(context.TODO(), &ssooidc.RegisterClientInput{
		ClientName: aws.String("sample-client-name"),
		ClientType: aws.String("public"),
		Scopes:     []string{"sso-portal:*"},
	})
	if err != nil {
		fmt.Println(err)
	}

	// authorize your device using the client registration response
	deviceAuth, err := ssooidcClient.StartDeviceAuthorization(context.TODO(), &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		StartUrl:     aws.String("https://tamg.awsapps.com/start"),
	})
	if err != nil {
		fmt.Println(err)
	}

	// trigger OIDC login. open browser to login. close tab once login is done. press enter to continue
	url := aws.ToString(deviceAuth.VerificationUriComplete)
	fmt.Printf("If browser is not opened automatically, please open link:\n%v\n", url)
	err = browser.OpenURL(url)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Press ENTER key once login is done")
	_ = bufio.NewScanner(os.Stdin).Scan()

	// generate sso token
	token, err := ssooidcClient.CreateToken(context.TODO(), &ssooidc.CreateTokenInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		DeviceCode:   deviceAuth.DeviceCode,
		GrantType:    aws.String("urn:ietf:params:oauth:grant-type:device_code"),
	})
	if err != nil {
		fmt.Println(err)
	}

	return token.AccessToken
}
func getCurrentToken() *string {
	// Best effort attempt to get token from sso cache.
	// If you can't for whatever reason, return nil, and the code will walk the user through generating a token

	usr, err := user.Current()
	if err != nil {
		return nil
	}
	dir := filepath.Join(usr.HomeDir, ".aws/sso/cache")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		filename := filepath.Join(dir, entry.Name())
		byteValue, err := os.ReadFile(filename)
		if err != nil {
			continue
		}
		var cacheEntry SSOCacheEntry
		err = json.Unmarshal(byteValue, &cacheEntry)
		// If the file ends in .json it should probably have valid json, but meh
		if err != nil {
			continue
		}
		// If it didn't fill in the fields, it is probably not a cache entry, some other random json
		if cacheEntry.AccessToken == "" || cacheEntry.ExpiresAt.IsZero() {
			continue
		}
		// Make sure it's not already expired
		if time.Now().After(cacheEntry.ExpiresAt) {
			continue
		}

		return &cacheEntry.AccessToken
	}

	return nil
}

func GetToken(cfg aws.Config) *string {
	token := getCurrentToken()
	if token != nil {
		return token
	}

	return generateToken(cfg)
}
