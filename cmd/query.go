package cmd

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/satoqz/lyr/query"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"get", "find", "search"},
	Args:    cobra.MinimumNArgs(1),
	RunE:    queryExec,
}

func queryExec(_ *cobra.Command, args []string) error {

	q := query.New(strings.Join(args, " "))
	res, err := q.Search()
	if err != nil {
		return err
	}
	data := res.Collect()

	if len(data) == 0 {
		return errors.New("No songs found. Exiting.\n")
	}

	song := data[0]

	var lyrics string
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		lyrics, err = song.ScrapeLyrics()
		wg.Done()
	}()

	wg.Wait()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n\n%s\n", song.Name, lyrics)
	return nil
}
