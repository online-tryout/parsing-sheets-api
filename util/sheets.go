package util

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsRowReader struct {
	Number string `json:"number"`
	Question string `json:"question"`
	Answer string `json:"answer"`
	Option []string `json:"option"`
}

func (sheet *SheetsRowReader) IsEmpty() bool {
	return len(sheet.Number) == 0 && len(sheet.Question) == 0 && len(sheet.Answer) == 0 && len(sheet.Question) == 0
}

func NumberToColumnLetter(n int64) string {
    var result string

    for n > 0 {
        remainder := (n - 1) % 26
        result = string(rune('A'+remainder)) + result
        n = (n - 1) / 26
    }

    return result
}

func GetSheetsClient(credentialsFile string) (*sheets.Service, error) {
    creds, err := os.ReadFile(credentialsFile)
    if err != nil {
        return nil, err
    }

    config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
    if err != nil {
        return nil, err
    }

    client := config.Client(context.Background())
    sheetsService, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))

    return sheetsService, err
}

func GetSpreadsheetInfo(srv *sheets.Service, spreadsheetID string) (*sheets.Spreadsheet, error) {
    spreadsheet, err := srv.Spreadsheets.Get(spreadsheetID).Do()
    if err != nil {
        return nil, fmt.Errorf("unable to retrieve spreadsheet: %v", err)
    }
    return spreadsheet, nil
}

func FetchData(srv *sheets.Service, spreadsheetID, sheetName, readRange string) ([][]interface{}, error) {
    resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, sheetName+"!"+readRange).Do()
    if err != nil {
        return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
    }

	if len(resp.Values) == 0 {
        return nil, fmt.Errorf("no data found in sheet")
    }

    return resp.Values, nil
}

func GetSheetID(url string) (string, error) {
	re := regexp.MustCompile(`\/spreadsheets\/d\/([a-zA-Z0-9-_]+)\/edit`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("URL is not a valid Google Sheets URL")
	}
	return matches[1], nil
}