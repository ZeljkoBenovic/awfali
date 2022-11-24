package main

import (
	"flag"
	"log"
	"os"
)

var userConfig = availableUserFlags{
	region: userFlagStr{
		value:        new(string),
		name:         "region",
		defaultValue: "eu-central-1",
		description:  "AWS region",
	},
	profile: userFlagStr{
		value:        new(string),
		name:         "profile",
		defaultValue: "mfa_user",
		description:  "AWS profile name",
	},
	serialNumber: userFlagStr{
		value:        new(string),
		name:         "serial",
		defaultValue: "virtual",
		description:  "The identification number of the MFA device",
	},
	mfaToken: userFlagStr{
		value:        new(string),
		name:         "mfa",
		defaultValue: "",
		description:  "MFA token for your AWS account",
	},
	mode: userFlagStr{
		value:        new(string),
		name:         "mode",
		defaultValue: "conf",
		description: "Set mode to store credentials as environment variables " +
			"or write them to the aws credentials file(env or conf). " +
			"Env mode will start a new shell with env vars loaded in.",
	},
	sessionDuration: userFlagInt{
		value:        new(int),
		name:         "duration",
		defaultValue: 43200,
		description:  "Session token duration",
	},
	confFile: userFlagStr{
		value:        new(string),
		name:         "conf",
		defaultValue: getHomeDir() + "/.aws/config",
		description:  "AWS config file location",
	},
	credFile: userFlagStr{
		value:        new(string),
		name:         "cred",
		defaultValue: getHomeDir() + "/.aws/credentials",
		description:  "AWS credentials file location",
	},
	accessKeyId: userFlagStr{
		value:        new(string),
		name:         "access-key-id",
		defaultValue: "",
		description:  "IAM access key id to authenticate the request",
	},
	secretAccessKey: userFlagStr{
		value:        new(string),
		name:         "secret-access-key",
		defaultValue: "",
		description:  "IAM secret access key to authenticate the request",
	},
}

type availableUserFlags struct {
	region          userFlagStr
	profile         userFlagStr
	serialNumber    userFlagStr
	mfaToken        userFlagStr
	mode            userFlagStr
	sessionDuration userFlagInt
	confFile        userFlagStr
	credFile        userFlagStr

	accessKeyId     userFlagStr
	secretAccessKey userFlagStr
}

type userFlagStr struct {
	value        *string
	name         string
	defaultValue string
	description  string
}

type userFlagInt struct {
	value        *int
	name         string
	defaultValue int
	description  string
}

func (f *availableUserFlags) Get() *availableUserFlags {
	flag.StringVar(f.region.value, f.region.name, f.region.defaultValue, f.region.description)
	flag.StringVar(f.profile.value, f.profile.name, f.profile.defaultValue, f.profile.description)
	flag.StringVar(f.serialNumber.value, f.serialNumber.name, f.serialNumber.defaultValue, f.serialNumber.description)
	flag.StringVar(f.mfaToken.value, f.mfaToken.name, f.mfaToken.defaultValue, f.mfaToken.description)
	flag.StringVar(f.mode.value, f.mode.name, f.mode.defaultValue, f.mode.description)
	flag.StringVar(f.confFile.value, f.confFile.name, f.confFile.defaultValue, f.confFile.description)
	flag.StringVar(f.credFile.value, f.credFile.name, f.credFile.defaultValue, f.credFile.description)
	flag.StringVar(f.accessKeyId.value, f.accessKeyId.name, f.accessKeyId.defaultValue, f.accessKeyId.description)
	flag.StringVar(f.secretAccessKey.value, f.secretAccessKey.name, f.secretAccessKey.defaultValue, f.secretAccessKey.description)
	flag.IntVar(f.sessionDuration.value, f.sessionDuration.name, f.sessionDuration.defaultValue, f.sessionDuration.description)

	flag.Parse()

	return f
}

func (f *availableUserFlags) CheckValidity() {
	// check for mandatory flags
	if *f.serialNumber.value == "" || *f.mfaToken.value == "" {
		log.Println("Serial number and mfa must be defined")
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func (f *availableUserFlags) AreIAMCredentialsSet() bool {
	return *f.accessKeyId.value != "" && *f.secretAccessKey.value != ""
}

func (f *availableUserFlags) SessionInt32() *int32 {
	i := int32(*f.sessionDuration.value)

	return &i
}

func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("could not read home dir path")
	}

	return home
}
