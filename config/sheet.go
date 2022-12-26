package config

import (
	"context"
	"github.com/margostino/climateline-processor/common"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"strconv"
)

type UrlConfig struct {
	Url            string
	Category       string
	Tags           string
	BotEnabled     bool
	TwitterEnabled bool
}

func GetUrls(inputCategory string) []*UrlConfig {
	urls := make([]*UrlConfig, 0)

	ctx := context.Background()
	api, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("GSHEET_API_KEY")))

	if !common.IsError(err, "when creating new Google API Service") {
		spreadsheetId := os.Getenv("SPREADSHEET_ID")
		readRange := os.Getenv("SPREADSHEET_RANGE")
		resp, err := api.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()

		if !common.IsError(err, "unable to retrieve data from sheet") && len(resp.Values) > 0 {
			for _, row := range resp.Values {
				var isBotEnabled, isTwitterEnabled bool
				var category string
				if len(row) == 5 {
					category = row[1].(string)
					isBotEnabled, err = strconv.ParseBool(row[2].(string))
					common.SilentCheck(err, "when fetching Bot enabled config from feed urls")
					isTwitterEnabled, err = strconv.ParseBool(row[3].(string))
					common.SilentCheck(err, "when fetching Twitter enabled config from feed urls")
				} else {
					log.Printf("Configuration sheet for Feed Urls is not valid. It must have 3 columns. It has %d\n", len(row))
				}

				matchCategory := (inputCategory != "*" && inputCategory == category) || inputCategory == "*"

				if matchCategory {
					//isBotEnabled = true
					//isTwitterEnabled = false
					urlConfig := &UrlConfig{
						Url:            row[0].(string),
						Tags:           row[4].(string),
						Category:       category,
						BotEnabled:     isBotEnabled,
						TwitterEnabled: isTwitterEnabled,
					}
					urls = append(urls, urlConfig)
				}

			}
		}
	}

	return urls
}
