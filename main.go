package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	userField     = "username"
	passwordField = "password"
	hostField     = "host"
	schemeField   = "protocol"
	pathField     = "path"
)

func main() {
	log.SetFlags(0)
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Home directory:", err)
	}
	creds, err := loadCredentials(filepath.Join(home, ".git-credentials"))
	if err != nil {
		log.Fatal("Load credentials:", err)
	}

	vals, err := readMap(os.Stdin)
	if err != nil {
		log.Fatal("Read values:", err)
	}

	cred := lookupCredential(vals, creds)
	if cred == nil {
		os.Exit(1)
	}

	vals[schemeField] = cred.Scheme
	vals[hostField] = cred.Host
	if v := cred.User.Username(); v != "" {
		vals[userField] = v
	}
	if v, ok := cred.User.Password(); ok && v != "" {
		vals[passwordField] = v
	}
	if v := cred.Path; v != "" {
		vals[pathField] = v
	}

	printMap(os.Stdout, vals)
}

func lookupCredential(m map[string]string, creds []*url.URL) *url.URL {
	for _, cred := range creds {
		if cred.Scheme != m[schemeField] {
			continue
		}
		if cred.Host != m[hostField] {
			continue
		}

		wantUser := m[userField]
		haveUser := cred.User.Username()
		if wantUser != "" && haveUser != "" && wantUser != haveUser {
			continue
		}

		wantPath := m[pathField]
		if wantPath != "" && cred.Path != "" && !strings.HasPrefix(wantPath, cred.Path) {
			continue
		}

		return cred
	}
	return nil
}

func loadCredentials(path string) ([]*url.URL, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	br := bufio.NewScanner(fd)
	var urls []*url.URL
	for br.Scan() {
		uri, err := url.Parse(br.Text())
		if err != nil {
			return nil, err
		}
		uri.Path = strings.TrimLeft(uri.Path, "/")
		urls = append(urls, uri)
	}
	return urls, nil
}

func readMap(r io.Reader) (map[string]string, error) {
	br := bufio.NewScanner(r)
	m := make(map[string]string)
	for br.Scan() {
		if br.Text() == "" {
			break
		}
		fields := strings.SplitN(br.Text(), "=", 2)
		if len(fields) != 2 {
			continue
		}
		m[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
	}
	return m, br.Err()
}

func printMap(w io.Writer, m map[string]string) {
	for key, val := range m {
		fmt.Fprintf(w, "%s=%s\n", key, val)
	}
}
