package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"log"
	"os"
	"strings"
	"syscall"
)

func main() {
	// get user flags
	userConfig.Get().CheckValidity()

	// setup aws config
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(*userConfig.region.value),
		config.WithSharedConfigProfile(*userConfig.profile.value))
	if err != nil {
		log.Fatalln("could not create aws config: ", err.Error())
	}
	// get session token
	crds, err := sts.NewFromConfig(cfg).GetSessionToken(ctx, &sts.GetSessionTokenInput{
		SerialNumber:    userConfig.serialNumber.value,
		TokenCode:       userConfig.mfaToken.value,
		DurationSeconds: userConfig.SessionInt32(),
	})
	if err != nil {
		log.Fatalln("could not get session token:", err.Error())
	}

	switch *userConfig.mode.value {
	case "env":
		if err := os.Setenv("AWS_ACCESS_KEY_ID", *crds.Credentials.AccessKeyId); err != nil {
			log.Fatalln("could not set AWS_ACCESS_KEY_ID env var:", err.Error())
		}
		if err := os.Setenv("AWS_SECRET_ACCESS_KEY", *crds.Credentials.SecretAccessKey); err != nil {
			log.Fatalln("could not set AWS_ACCESS_KEY_ID env var:", err.Error())
		}
		if err := os.Setenv("AWS_SESSION_TOKEN", *crds.Credentials.SessionToken); err != nil {
			log.Fatalln("could not set AWS_ACCESS_KEY_ID env var:", err.Error())
		}
		if err := os.Setenv("AWS_REGION", *userConfig.region.value); err != nil {
			log.Fatalln("could not set AWS_ACCESS_KEY_ID env var:", err.Error())
		}

		// set env vars to shell
		if err := syscall.Exec(os.Getenv("SHELL"),
			[]string{os.Getenv("SHELL")}, syscall.Environ()); err != nil {
			log.Fatalln("could not set environment vars:", err.Error())
		}

		log.Println("credentials successfully stored into environment variables")
	case "conf":
		if err := storeToConfigFile(crds, userConfig); err != nil {
			log.Fatalln("could not store credentials to local files:", err.Error())
		}

		log.Println("credentials successfully stored to aws credentials/config files")
	default:
		log.Fatalln("Mode not supported. Supported modes are: [env, mode]")
	}
}

func storeToConfigFile(crds *sts.GetSessionTokenOutput, userConf availableUserFlags) error {
	// config file content to write
	configContent := fmt.Sprintf(`
[profile %s]
region=%s
`, *userConf.profile.value, *userConf.region.value)
	// credentials file content to write
	credContent := fmt.Sprintf(`
[%s]
aws_access_key_id=%s
aws_secret_access_key=%s
aws_session_token=%s
`, *userConf.profile.value, *crds.Credentials.AccessKeyId,
		*crds.Credentials.SecretAccessKey, *crds.Credentials.SessionToken)

	// read config file
	currentConf, err := os.ReadFile(*userConf.confFile.value)
	if err != nil {
		return fmt.Errorf("could not read config file err=%w", err)
	}
	// read credentials file
	currentCred, err := os.ReadFile(*userConf.credFile.value)
	if err != nil {
		return fmt.Errorf("coudl not read credentials file err=%w", err)
	}

	// remove old information from config
	tempConf := strings.Split(string(currentConf), "\n")
	for i, line := range tempConf {
		// delete old profile and data if it contains region
		if strings.TrimSpace(line) == fmt.Sprintf("[profile %s]", *userConf.profile.value) {
			tempConf[i] = ""
			if strings.Contains(tempConf[i+1], "region") {
				tempConf[i+1] = ""
			}
		}
	}
	// remove old information from credentials
	tempCred := strings.Split(string(currentCred), "\n")
	for i, line := range tempCred {
		if strings.TrimSpace(line) == fmt.Sprintf("[%s]", *userConf.profile.value) {
			tempCred[i] = ""
			if strings.Contains(tempCred[i+1], "aws_access_key_id") {
				tempCred[i+1] = ""
			}
			if strings.Contains(tempCred[i+2], "aws_secret_access_key") {
				tempCred[i+2] = ""
			}
			if strings.Contains(tempCred[i+3], "aws_session_token") {
				tempCred[i+3] = ""
			}
		}
	}

	// append existing data
	newCred := strings.Join(purgeEmpty(tempCred), "\n") + credContent
	newConf := strings.Join(purgeEmpty(tempConf), "\n") + configContent

	// write config file
	if err = os.WriteFile(*userConf.credFile.value, []byte(newCred), 0600); err != nil {
		return fmt.Errorf("could not write config file err=%w", err)
	}
	// append credentials file
	if err = os.WriteFile(*userConf.confFile.value, []byte(newConf), 0600); err != nil {
		return fmt.Errorf("could not write credentials file err=%w", err)
	}

	return nil
}

func purgeEmpty(s []string) []string {
	var r []string

	for i, v := range s {
		if v != "" {
			r = append(r, s[i])
		}
	}

	return r
}
