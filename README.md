# AWFALI (aws-mfa-cli)
A handy tool for generating AWS Session Token on the fly.   
In the situations where an AWS IAM user is required to have MFA authentication even from the cli.    
A user must generate a Session Token using MFA Serial Number or ARN and a MFA generated token.

More information on [get-session-token](https://docs.aws.amazon.com/cli/latest/reference/sts/get-session-token.html)

## Usage
Download the release binary for your OS, and run the binary with parameters.

### Modes
#### CONF
By default (`-mode conf`), the program will grab session token from AWS and place it to the local AWS files (`.aws/.config` and `.aws/.credentials`).    
#### ENV
When `-mode` is set to `env`, the credentials will be stored as environment variables.    
When running in this mode, it is advisable to run the binary with `exec` to prevent opening a shell in shell, 
as the program needs to set shell environment variables.

### Parameters
*  `-conf string` AWS config file location (default "~/.aws/config")
*  `-cred string` AWS credentials file location (default "~/.aws/credentials")
*  `-duration int` Session token duration (default 43200)
*  `-mfa string` MFA token for your AWS account
*  `-mode string` Set mode to store credentials as environment variables or write them to the aws credentials file(env or conf). Env mode will start a new shell with env vars loaded in. (default "conf")
*  `-profile string` AWS profile name (default "mfa_user")
*  `-region string` AWS region (default "eu-central-1")
*  `-serial string` The identification number of the MFA device, hardware serial number or user ARN
