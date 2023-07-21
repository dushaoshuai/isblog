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
	buf bufio.Writer
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
	buf.Reset(f)

	// todo: How to do error handling better?
	_, err = buf.WriteString(strconv.Itoa(i.Number))
	if err != nil {
		return err
	}
	err = buf.WriteByte('\n')
	if err != nil {
		return err
	}
	_, err = buf.WriteString(i.Title)
	if err != nil {
		return err
	}
	err = buf.WriteByte('\n')
	if err != nil {
		return err
	}
	_, err = buf.WriteString(i.Body)
	if err != nil {
		return err
	}
	return buf.Flush()
}

type httpReq struct {
	method      string
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
