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
	api, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("GSHEET_API_KEY")))

	if !common.IsError(err, "when creating new Google API Service") {
		spreadsheetId := os.Getenv("SPREADSHEET_ID")
		readRange := os.Getenv("SPREADSHEET_RANGE")
		resp, err := api.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()

		if !common.IsError(err, "unable to retrieve data from sheet") && len(resp.Values) > 0 {
			for _, row := range resp.Values {
				urls = append(urls, row[0].(string))
			}
		}
	}

	return urls
}

func Mock() {
	ctx := context.Background()
	api, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("GSHEET_API_KEY")))

	if !common.IsError(err, "when creating new Google API Service") {
		var vr sheets.ValueRange
		spreadsheetId := os.Getenv("SPREADSHEET_ID")
		writeRange := os.Getenv("SPREADSHEET_RANGE")
		myval := []interface{}{"One", "Two", "Three"}
		vr.Values = append(vr.Values, myval)
		_, err = api.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
		common.Check(err, "unable to update data into sheet")
	}
}
