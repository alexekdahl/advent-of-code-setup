package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Cred struct {
	SessionCookie string `json:"SESSION_COOKIE"`
}

type Config struct {
	Year          string
	Day           string
	SessionCookie string
	Path          string
}

var DIR_NAME = map[string]string{
	"1":  "one",
	"2":  "two",
	"3":  "three",
	"4":  "four",
	"5":  "five",
	"6":  "six",
	"7":  "seven",
	"8":  "eight",
	"9":  "nine",
	"10": "ten",
	"11": "eleven",
	"12": "twelve",
	"13": "thirteen",
	"14": "fourteen",
	"15": "fifteen",
	"16": "sixteen",
	"17": "seventeen",
	"18": "eighteen",
	"19": "nineteen",
	"20": "twenty",
	"21": "twenty_one",
	"22": "twenty_two",
	"23": "twenty_three",
	"24": "twenty_four",
}

//go:embed creds.json
var content embed.FS

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Need arguments")
		os.Exit(1)
	}

	year, day := os.Args[1], os.Args[2]

	if year == "" || day == "" {
		fmt.Println("Both year and day needs to be provided")
		os.Exit(1)
	}

	err := assertArgs(year, day)
	if err != nil {
		log.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dirName := year + "_" + DIR_NAME[day]
	path := filepath.Join(pwd, dirName)

	err = assertAndCreateFolder(path)
	if err != nil {
		log.Fatal(err)
	}

	sessionCookie, err := getSessionCookie()
	if err != nil {
		log.Fatal(err)
	}
	config := Config{
		Year:          year,
		Day:           day,
		Path:          filepath.Join(path, "input.txt"),
		SessionCookie: sessionCookie,
	}
	err = downloadInput(config)
	if err != nil {
		log.Fatal(err)
	}
}

func assertArgs(year string, day string) error {
	if len(year) != 4 {
		return errors.New("Not a valid year")
	}

	if len(day) > 2 {
		return errors.New("Not a valid day")
	}

	_, err := strconv.Atoi(day)
	if err != nil {
		return errors.New("Not a valid day")
	}

	return nil
}

func assertAndCreateFolder(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func getSessionCookie() (string, error) {
	bytes, err := content.ReadFile("creds.json")
	if err != nil {
		return "", err
	}

	var cred Cred
	json.Unmarshal(bytes, &cred)
	if cred.SessionCookie == "" {
		return "", errors.New("No session")
	}

	return cred.SessionCookie, nil
}

func downloadInput(config Config) error {
	client := new(http.Client)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%s/day/%s/input", config.Year, config.Day), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "github.com/AlexEkdahl")

	cookie := new(http.Cookie)
	cookie.Name, cookie.Value = "session", config.SessionCookie
	req.AddCookie(cookie)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	file, err := os.OpenFile(config.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
