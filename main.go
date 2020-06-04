package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s <get|store|erase>\n", os.Args[0])
		os.Exit(2)
	}

	switch flag.Arg(0) {
	case "store", "erase":
		os.Exit(0)
	case "get":
		handleGet()
	default:
		fmt.Printf("Usage: %s <get|store|erase>\n", os.Args[0])
		os.Exit(2)
	}
}

func handleGet() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Home directory:", err)
		os.Exit(1)
	}

	creds, err := loadCredentials(filepath.Join(home, ".git-credentials"))
	if err != nil {
		fmt.Println("Load credentials:", err)
		os.Exit(1)
	}

	vals, err := readMap(os.Stdin)
	if err != nil {
		fmt.Println("Read values:", err)
		os.Exit(1)
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
