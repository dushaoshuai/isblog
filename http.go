package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type issues []issue

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/dushaoshuai/dushaoshuai.github.io/issues", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	q := url.Values{}
	q.Add("per_page", strconv.Itoa(1))
	q.Add("page", strconv.Itoa(2))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		panic(err)
	}

	var issus issues
	err = json.Unmarshal(body, &issus)
	if err != nil {
		panic(err)
	}

	for _, issu := range issus {
		fname := fmt.Sprintf("%04d %s", issu.Number, issu.Title)
		f, err := os.Create(fname)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(issu.Body)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}
