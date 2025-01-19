package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/clover0/issue-agent/cli"
	"github.com/clover0/issue-agent/config"
)

const defaultConfigPath = "./issue_agent.yml"

// This value is set at release build time
// ldflags "-X github.com/clover0/issue-agent/main.containerImageTag=v0.0.1"
var containerImageTag = "dev"

// Use the docker command to start a container and execute the agent binary
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configPath, err := getConfigPathOrDefault()
	if err != nil {
		panic(err)
	}

	conf, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}

	promptPath, err := getPromptPath(conf)
	if err != nil {
		panic(err)
	}

	flags, err := parseArgs()
	if err != nil {
		panic(err)
	}

	var awsDockerEnvs []string
	if flags.Common.AWSProfile != "" {
		awsKeys := getAWSKeys(flags.Common.AWSProfile)
		fmt.Println("Using AWS credentials from profile:", flags.Common.AWSProfile)
		awsDockerEnvs = append(awsDockerEnvs, "-e", "AWS_REGION="+awsKeys.Region)
		awsDockerEnvs = append(awsDockerEnvs, "-e", "AWS_ACCESS_KEY_ID="+awsKeys.AccessKeyID)
		awsDockerEnvs = append(awsDockerEnvs, "-e", "AWS_SECRET_ACCESS_KEY="+awsKeys.SecretAccessKey)
		awsDockerEnvs = append(awsDockerEnvs, "-e", "AWS_SESSION_TOKEN="+awsKeys.SessionToken)
	}

	imageName := "ghcr.io/clover0/issue-agent"
	imageTag := containerImageTag
	dockerEnvs := passEnvs()
	containerName := "issue-agent"
	args := []string{
		"run",
		"--rm",
		"--name", containerName,
	}
	// Mount files to the container
	if len(configPath) > 0 {
		args = append(args, "-v", configPath+":"+config.ConfigFilePath)
	}
	if len(promptPath) > 0 {
		path, err := filepath.Abs(conf.Agent.PromptPath)
		if err != nil {
			panic(err)
		}
		args = append(args, "-v", path+":"+config.PromptFilePath)
	}
	args = append(args, dockerEnvs...)
	args = append(args, awsDockerEnvs...)
	args = append(args, imageName+":"+imageTag)
	args = append(args, os.Args[1:]...)
	for _, a := range os.Args[1:] {
		if strings.HasSuffix(a, "-config") {
			break
		}
		args = append(args, "-config", configPath)
	}

	cmd := exec.CommandContext(ctx, dockerCmd(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error running container:", err)
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func(containerName string) {
		sig := <-sigChan
		fmt.Println("Received signal")
		if err := cmd.Process.Signal(sig); err != nil {
			fmt.Println("Error sending signal to container:", err)
		}
		stopContainer(containerName)
		cancel()
	}(containerName)

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
			fmt.Printf("Process exited with error: %v\n", err)
			os.Exit(exitErr.ExitCode())
		}
	}
}

func parseArgs() (*cli.IssueInputs, error) {
	flags, mapper := cli.IssueFlags()

	start := 1
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "-") {
			start = i
			break
		}
	}

	buf := bytes.NewBuffer([]byte{})
	flags.SetOutput(buf)

	if err := flags.Parse(os.Args[start:]); err != nil {
		if strings.Contains(err.Error(), "flag provided but not defined") {
			// pass to the next starting container
			fmt.Printf("Parsed input: %v\n", mapper)
		}
		return mapper, fmt.Errorf("failed to parse input: %w", err)
	}

	return mapper, nil
}

func dockerCmd() string {
	com, ok := os.LookupEnv("_DOCKER_CMD")
	if ok {
		return com
	}

	return "docker"
}

func stopContainer(containerName string) {
	cmd := exec.Command(dockerCmd(), "kill", containerName)
	bytes, _ := cmd.CombinedOutput()
	fmt.Println(string(bytes))
}

func getConfigPathOrDefault() (string, error) {
	configStart := len(os.Args)
	foundConfig := false
	for i, arg := range os.Args {
		if strings.HasSuffix(arg, "-config") {
			configStart = i
			foundConfig = true
			break
		}
	}

	// when -config option not found
	// default config file or empty
	if !foundConfig {
		if _, err := os.Stat(defaultConfigPath); err != nil {
			return "", nil
		}
		return filepath.Abs(defaultConfigPath)
	}

	if len(os.Args) <= configStart+1 {
		return "", fmt.Errorf("-config option value is required")
	}

	path := os.Args[configStart+1]
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func getPromptPath(conf config.Config) (string, error) {
	if len(conf.Agent.PromptPath) == 0 {
		return "", nil
	}
	return filepath.Abs(conf.Agent.PromptPath)
}

// Pass only the environment variables that are required by the agent.
// This is to avoid passing sensitive information to the container.
func passEnvs() []string {
	var passEnvs []string
	for _, env := range os.Environ() {
		envName := strings.Split(env, "=")[0]
		if slices.Contains(cli.EnvNames(), envName) {
			passEnvs = append(passEnvs, env)
		}
	}

	var dockerEnvs []string
	for _, env := range passEnvs {
		varName := strings.Split(env, "=")
		if len(varName) == 2 {
			dockerEnvs = append(dockerEnvs, "-e", env)
		}
	}

	return dockerEnvs
}

type awsCredentials struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func getAWSKeys(profile string) awsCredentials {
	ctx := context.Background()

	var opts []func(*awsconfig.LoadOptions) error
	if profile != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(profile))
	}

	sdkConfig, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return awsCredentials{}
	}

	cred, err := sdkConfig.Credentials.Retrieve(ctx)
	if err != nil {
		fmt.Println("Couldn't retrieve credentials")
		fmt.Println(err)
		return awsCredentials{}
	}

	return awsCredentials{
		Region:          sdkConfig.Region,
		AccessKeyID:     cred.AccessKeyID,
		SecretAccessKey: cred.SecretAccessKey,
		SessionToken:    cred.SessionToken,
	}
}
