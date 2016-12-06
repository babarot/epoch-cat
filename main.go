package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hhkbp2/go-strftime"
)

const (
	unixTimeBegin = 946684800 // 946684800 2000-01-01T00:00:00+0000
)

var (
	unixTimeRe = regexp.MustCompile(`([0-9]{9,}(\.[0-9]+)?)`) // e.g. 1476707866.8115
	// dateTimeRe TODO
)

var (
	format = flag.String("f", "", "Specify the time format for convert")
	quote  = flag.Bool("q", false, "Quote the date to be output")
	prefix = flag.String("p", "", "Add the prefix for searching to prevent mismatch when looking for UNIX time")
)

func init() {
	var location string

	if os.Getenv("TZ") == "" {
		location = "UTC"
	} else {
		location = os.Getenv("TZ")
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func doEchoLine(r io.Reader) error {
	var formattedTime string
	var buf bytes.Buffer

	unixTimeNow := time.Now().Unix()
	br := bufio.NewReader(r)

	for {
		line, isPrefix, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if _, err = buf.Write(line); err != nil {
			return err
		}
		if !isPrefix {
			var matched [][]byte
			if *prefix != "" {
				unixTimeRe = regexp.MustCompile(*prefix + `([0-9]{9,}(\.[0-9]+)?)`)
			}
			matched = unixTimeRe.FindSubmatch(line)
			if len(matched) > 0 {
				matchedUnixTime, err := strconv.ParseFloat(strings.TrimPrefix(string(matched[0]), *prefix), 64)
				if err != nil {
					panic(err)
					fmt.Println(line)
					continue
				}
				unixTime := int64(matchedUnixTime)
				if *format == "" {
					formattedTime = time.Unix(unixTime, 0).Format(time.RFC3339)
				} else {
					formattedTime = strftime.Format(*format, time.Unix(unixTime, 0))
				}
				if *quote {
					formattedTime = fmt.Sprintf(`"%s"`, formattedTime)
				}
				if *prefix != "" {
					formattedTime = *prefix + formattedTime
				}
				if unixTimeBegin <= unixTime && unixTime <= unixTimeNow {
					fmt.Println(unixTimeRe.ReplaceAllString(string(line), formattedTime))
					continue
				}
			}
			buf.Reset()
		}
	}
	return nil
}

func openFile(s string) (io.ReadWriteCloser, error) {
	fi, err := os.Stat(s)
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeSocket != 0 {
		return net.Dial("unix", s)
	}
	return os.Open(s)
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		doEchoLine(os.Stdin)
	}

	for _, fname := range flag.Args() {
		if fname == "-" {
			doEchoLine(os.Stdin)
		} else {
			f, err := openFile(fname)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer f.Close()
			doEchoLine(f)
		}
	}
}
