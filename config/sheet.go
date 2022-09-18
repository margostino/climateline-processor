package config

import (
	"context"
	"github.com/margostino/climateline-processor/common"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
)

func GetUrls() []string {
	urls := make([]string, 0)

	ctx := context.Background()
	sheets, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("GSHEET_API_KEY")))

	if !common.IsError(err, "when creating new Google API Service") {
		spreadsheetId := os.Getenv("SPREADSHEET_ID")
		readRange := os.Getenv("SPREADSHEET_RANGE")
		resp, err := sheets.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()

		if !common.IsError(err, "unable to retrieve data from sheet") && len(resp.Values) > 0 {
			for _, row := range resp.Values {
				urls = append(urls, row[0].(string))
			}
		}
	}

	return urls
}
