/*
Copyright © 2023 shaouai <shaoshuai@shuai.host>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	issueNum     int
	issueNumName = "issue-number"
)

func pullOne(cmd *cobra.Command, args []string) error {
	req := httpReq{
		method:     http.MethodGet,
		pathParams: []string{strconv.Itoa(issueNum)},
	}
	resp, err := req.do(cmd.Context())
	if err != nil {
		return err
	}

	var issu issue
	err = json.Unmarshal(resp, &issu)
	if err != nil {
		return err
	}
	return issu.toFile()
}

func pullList(cmd *cobra.Command, args []string) error {
	var (
		ctx = cmd.Context()

		perPage   = 100
		page      = 0
		issueChan = make(chan issue, perPage)

		es      []error
		errChan = make(chan error)
	)

	go func() {
		for err := range errChan {
			es = append(es, err)
		}
	}()

	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				break loop
			default:
				page++
			}

			q := url.Values{}
			q.Add("per_page", strconv.Itoa(perPage))
			q.Add("page", strconv.Itoa(page))
			req := httpReq{
				method:      http.MethodGet,
				queryParams: q,
			}

			resp, err := req.do(ctx)
			if err != nil {
				errChan <- err
				break // this should be a fatal error
			}

			var issueSlice []issue
			err = json.Unmarshal(resp, &issueSlice)
			if err != nil {
				errChan <- err
				break // this should be a fatal error
			}

			for _, issu := range issueSlice {
				issueChan <- issu
			}
			if len(issueSlice) < perPage {
				break
			}
		}

		close(issueChan)
	}()

	for issu := range issueChan {
		err := issu.toFile()
		if err != nil {
			errChan <- err
		}
	}
	close(errChan)

	return errors.Join(es...)
}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull the remote Github issue(s) to the local blog(s)",
	Long: `Pull the remote Github issue(s) to the local blog(s).
It overwrites any existing local blogs, so be careful not to lose local changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed(issueNumName) {
			return pullOne(cmd, args)
		}
		return pullList(cmd, args)
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pullCmd.Flags().IntVarP(&issueNum, issueNumName, "i", 0, "The number that identifies the issue")
}
