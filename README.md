# Advent of Code Input Downloader

This is a simple command-line tool written in Go that helps you download the input for a specific day of the Advent of Code event.

## Usage

To use this tool, you need to provide the year and day of the Advent of Code event as command-line arguments. For example, to download the input for day 1 of the 2021 event, you would run the following command:

```bash
go run main.go -y 2021 -d 1
```
The tool will create a directory for the specified day in the current working directory and download the input file into it.

## Session Cookie

In order to download the input, you need to provide your Advent of Code session cookie. This can be done by setting the `ADVENT_SESSION_COOKIE` environment variable to the value of your session cookie.

## Building

To build the tool, simply run the following command:

```bash
go build
```
