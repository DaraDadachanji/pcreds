package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

type Profiles map[string]Profile

type Profile map[string]string

func (p *Profiles) Marshal() []byte {
	var lines []string
	for name, profile := range *p {
		line := fmt.Sprintf("[%s]", name)
		lines = append(lines, line)
		for key, value := range profile {
			line = fmt.Sprintf("%s = %s", key, value)
			lines = append(lines, line)
		}
		lines = append(lines, "") //separate profiles with blank line
	}
	contents := []byte(strings.Join(lines, "\n"))
	return contents

}

func Unmarshal(data []byte) Profiles {
	reader := bufio.NewReader(bytes.NewReader(data))
	profiles := Profiles{}
	var profileName string
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if IsBlank(line) {
			continue
		}
		if IsProfileName(line) {
			profileName = ParseProfileName(line)
			profiles[profileName] = Profile{}
			continue
		}
		key, value := ParseKeyValue(line)
		profiles[profileName][key] = value
	}
	return profiles
}

func IsBlank(line string) bool {
	match, _ := regexp.MatchString(`^\s*$`, line)
	return match
}

func IsProfileName(line string) bool {
	match, _ := regexp.MatchString(`^\[[ A-Za-z0-9\-_]+]`, line)
	return match
}

func ParseProfileName(line string) string {
	r := regexp.MustCompile(`\[([^]]*)]`)
	parts := r.FindStringSubmatch(line)
	if len(parts) < 2 {
		log.Panic("could not parse credentials file")
	}
	return parts[1]
}

func ParseKeyValue(line string) (key string, value string) {
	r := regexp.MustCompile(`([^= ]*)\s*=\s*(\S*)`)
	parts := r.FindStringSubmatch(line)
	if len(parts) < 3 {
		log.Println(line)
		log.Panic("could not parse credentials file")
	}
	return parts[1], parts[2]
}

func CredentialsFilepath() string {
	return filepath.Join(HomeDirectory(), ".aws", "credentials")
}
