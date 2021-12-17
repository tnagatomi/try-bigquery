package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"os"
)

var projectID string
var dataset string

func init() {
	projectID = os.Getenv("GCP_PROJECT")
	dataset = os.Getenv("BQ_DATASET")
}

func main() {
	err := queryBasic(os.Stdout, projectID)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

// queryBasic demonstrates issuing a query and reading results.
func queryBasic(w io.Writer, projectID string) error {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	q := client.Query(
		"SELECT	name, count	FROM " + dataset + ".names_2014 " +
			"WHERE gender = 'M' " +
			"ORDER BY count DESC " +
			"LIMIT 5")
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"
	// Run the query and print results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	it, err := job.Read(ctx)
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(w, row)
	}
	return nil
}
