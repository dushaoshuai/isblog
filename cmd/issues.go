package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	bufW bufio.Writer
)

type issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (i *issue) fileName() string {
	mapping := func(r rune) rune {
		switch r {
		case ' ':
			return '_'
		case '/': // https://stackoverflow.com/a/9847306
			return '\u2044'
		default:
			return r
		}
	}
	title := strings.Map(mapping, i.Title)

	return fmt.Sprintf("%04d_%s.md", i.Number, title)
}

func (i *issue) toFile() error {
	f, err := os.Create(i.fileName())
	if err != nil {
		return err
	}
	defer f.Close()
	bufW.Reset(f)

	bufW.WriteString("===================================\n")
	bufW.WriteString("Do not edit the issue number line.\n")
	bufW.WriteString("Do not delete the issue title line.\n")
	bufW.WriteString(strconv.Itoa(i.Number))
	bufW.WriteByte('\n')
	bufW.WriteString(i.Title)
	bufW.WriteByte('\n')
	bufW.WriteString("===================================\n\n")
	bufW.WriteString(i.Body)
	return bufW.Flush()
}

func (i *issue) fromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	var (
		fileScanner = bufio.NewScanner(f)
		bodyBuilder strings.Builder // to build issue body

		l = 0 // line number
	)
	for fileScanner.Scan() {
		l++
		if l <= 3 || l == 6 || l == 7 {
			continue
		}
		if l == 4 { // issue number
			n, err := strconv.Atoi(fileScanner.Text())
			if err != nil {
				return err
			}
			i.Number = n
			continue
		}
		if l == 5 { // issue title
			i.Title = fileScanner.Text()
			continue
		}
		bodyBuilder.Write(fileScanner.Bytes())
		bodyBuilder.WriteByte('\n')
	}
	if fileScanner.Err() != nil {
		return err
	}
	i.Body = bodyBuilder.String()

	return nil
}

type httpReq struct {
	method      string
	needAuth    bool
	pathParams  []string
	body        io.Reader
	queryParams url.Values
}

func (r *httpReq) do(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, r.method, "https://api.github.com/repos/dushaoshuai/dushaoshuai.github.io/issues", r.body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	if r.needAuth {
		req.Header.Add("Authorization", "Bearer <TOKEN>")
	}
	if len(r.pathParams) != 0 {
		req.URL = req.URL.JoinPath(r.pathParams...)
	}

	if r.queryParams != nil {
		req.URL.RawQuery = r.queryParams.Encode()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Query Github REST API failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
		return nil, err
	}

	return body, nil
}
