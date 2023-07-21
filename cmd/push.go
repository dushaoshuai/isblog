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
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	localBlog     string
	localBlogName = "file"
)

func pushRunE(cmd *cobra.Command, args []string) error {
	var issu issue
	err := issu.fromFile(localBlog)
	if err != nil {
		return err
	}
	body, err := json.Marshal(issu)
	if err != nil {
		return err
	}

	req := httpReq{
		method:     http.MethodPatch,
		pathParams: []string{strconv.Itoa(issu.Number)},
		body:       bytes.NewReader(body),
	}
	_, err = req.do(cmd.Context())
	return err
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push -f file",
	Short: "Push the local blog to the remote Github issue",
	Long: `Push the local blog to the remote Github issue.
Currently the -f flag is required.`,
	RunE: pushRunE,
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pushCmd.Flags().StringVarP(&localBlog, localBlogName, "f", "", "The local blog file")
	pushCmd.MarkFlagRequired(localBlogName)
}
