package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func main() {
	lines := ReadStdIn()
	if len(lines) != 4 {
		log.Fatalln("expected 4 lines, recieved ", len(lines))
	}
	profileLine := lines[0]
	if !IsProfileName(profileLine) {
		log.Fatal("First line is not a profile tag")
	}
	name := ParseProfileName(profileLine)
	profile = GetAlias(name)
	file, err := ReadCredentialsFile()
	if err != nil {
		log.Fatalln(err)
	}
	credentials := Unmarshal(file)

	contents := credentials.Marshal()
	err = ioutil.WriteFile("hello", contents, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAlias(name string) string {
	aliases := ReadConfigFile()
	var alias string
	if aliases == nil {
		alias = "default"
	} else {
		alias = aliases[name]
	}
	return alias
}

func ReadStdIn() []string {
	reader := bufio.NewReader(os.Stdin)
	var lines []string
	for {
		line, readErr := reader.ReadString('\n')
		lines = append(lines, line)
		if readErr == io.EOF {
			return lines
		}
	}
}

func ReadConfigFile() map[string]string {
	const configFile = "~/.aws/pcreds.yaml"

	if fileExists(configFile) {
		data, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalln("could not read config file: ", err)
		}
		type Config struct {
			Profiles map[string]string `yaml:"profiles"`
		}
		var config Config
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return nil
		}
		return config.Profiles
	}
	return nil
}

func ReadCredentialsFile() ([]byte, error) {
	credentialsFile := "~/.aws/credentials"
	if !fileExists(credentialsFile) {
		return nil, fmt.Errorf("file not found: %s", credentialsFile)
	}
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}
	return true
}
