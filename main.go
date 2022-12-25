package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type Cred struct {
	SessionCookie string `json:"sessionCookie"`
}

type Config struct {
	Year          string
	Day           int
	SessionCookie string
	Path          string
}

func newConfig(y string, d int, sc string, p string) Config {
	return Config{
		Year:          y,
		Day:           d,
		SessionCookie: sc,
		Path:          p,
	}
}

var (
	year     = flag.String("year", "", "year of Advent of Code")
	day      = flag.Int("day", 0, "day of Advent of Code")
	dirnames = []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
		"eleven",
		"twelve",
		"thirteen",
		"fourteen",
		"fifteen",
		"sixteen",
		"seventeen",
		"eighteen",
		"nineteen",
		"twenty",
		"twenty_one",
		"twenty_two",
		"twenty_three",
		"twenty_four",
	}
)

//go:embed creds.json
var content embed.FS

func main() {
	flag.Parse()

	err := validateArgs(*year, *day)
	if err != nil {
		log.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dirName := fmt.Sprintf("%s_%s", *year, dirnames[*day-1])
	path := filepath.Join(pwd, dirName)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	sessionCookie, err := getSessionCookie()
	if err != nil {
		log.Fatal(err)
	}

	path = filepath.Join(path, "input.txt")
	config := newConfig(*year, *day, sessionCookie, path)

	if err := downloadInput(config); err != nil {
		log.Fatal(err)
	}
}

func getSessionCookie() (string, error) {
	var cred Cred

	bytes, err := content.ReadFile("creds.json")
	if err != nil {
		return "", fmt.Errorf("failed to read creds.json: %w", err)
	}

	if err := json.Unmarshal(bytes, &cred); err != nil {
		return "", fmt.Errorf("failed to unmarshal creds.json: %w", err)
	}

	return cred.SessionCookie, nil
}

func validateArgs(y string, d int) error {
	if y == "" || d == 0 {
		return fmt.Errorf("both year and day needs to be provided")
	}

	if v := isValidYear(y); !v {
		return fmt.Errorf("not a valid year")
	}

	return nil
}

func isValidYear(year string) bool {
	match, _ := regexp.MatchString("^[0-9]{4}$", year)

	return match
}

func downloadInput(config Config) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%s/day/%d/input", config.Year, config.Day), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: config.SessionCookie,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download input: %s", resp.Status)
	}

	if _, err := os.Stat(config.Path); err == nil {
		// file exists, return without performing copy operation
		return nil
	}

	f, err := os.Create(config.Path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
