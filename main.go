package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
	
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/yaml.v3"
)



var configFilePath = fmt.Sprintf("%s/.aws/config", os.Getenv("HOME"))



type roleConfig struct {
	Role string `yaml:"role"`
	MFA  string `yaml:"string"`
}

type config map[string]roleConfig


func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <role> ", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}


func readToken() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "Enter MFA Code:")
	mfatoken, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(mfatoken), nil
}


func readConfig() (config, error) {
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	roleConfig := make(config)
	return roleConfig, yaml.Unmarshal(file, &roleConfig)
}

func cherror(err error) {
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "error:%v\n", err)
		os.Exit(1)
	}

}

func assumeProfile(profile string) (*credentials.Value, error) {
	apsession := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profile,
		SharedConfigState: session.SharedConfigEnable,
		AssumeRoleTokenProvider: readToken,
	}))

	creds, err := apsession.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	return &creds, nil
}



func printCredentialExports(role string, creds *credentials.Value) {
	fmt.Printf("export AWS_ACCESS_KEY_ID=\"%s\"\n", creds.AccessKeyID)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=\"%s\"\n", creds.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("export AWS_SECURITY_TOKEN=\"%s\"\n", creds.SessionToken)
	fmt.Printf("export ASSUMED_ROLE=\"%s\"\n", role)
	fmt.Printf("# Run this to configure your shell:\n")
	fmt.Printf("# eval $(%s)\n", strings.Join(os.Args, " "))
}



func main() {
	duration := flag.Duration("duration", time.Hour, "total duration credentials will be valid for.")
	flag.Parse()
	argv := flag.Args()
	if len(argv) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	stscreds.DefaultDuration = *duration

	role := argv[0]

	var creds *credentials.Value
	var err error
	creds, err = assumeProfile(role)

	
	cherror(err)
	printCredentialExports(role, creds)
}