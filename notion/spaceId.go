package notion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	tokenV2         string = os.Getenv("NOTION_TOKEN_V2")
	spaceId         string
	spaceIdFilePath = os.Getenv("HOME") + "/.notion_space_id" // space id file

	client      = &http.Client{Timeout: 60 * time.Second}
	tokenCookie = http.Cookie{Name: "token_v2", Value: tokenV2}
)

func CheckToken() error {
	if tokenV2 == "" {
		return fmt.Errorf("Please set NOTION_TOKEN environment variable")
	}
	sid, err := getSpaceId()
	if err != nil {
		return err
	}
	spaceId = *sid
	return nil
}

func getSpaceId() (*string, error) {
	// check
	if _, err := os.Stat(spaceIdFilePath); err == nil {
		f, err := os.Open(spaceIdFilePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		sid := string(b)
		return &sid, nil
	}

	httpReq, err := http.NewRequest("POST", GetSpaceIdUrl, nil)
	if err != nil {
		return nil, err
	}
	httpReq.AddCookie(&tokenCookie)
	httpReq.Header.Add("Content-Type", "application/json")
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := new(map[string]interface{})
	err = json.Unmarshal(b, response)
	if err != nil {
		return nil, err
	}

	spaceId = findSpaceId(*response)
	if spaceId == "" {
		return nil, fmt.Errorf("spaceId not found")
	}
	return &spaceId, saveSpaceId(spaceId)
}

func saveSpaceId(spaceId string) error {
	return os.WriteFile(spaceIdFilePath, []byte(spaceId), 644)
}

func findSpaceId(data map[string]interface{}) string {
	for key, value := range data {
		switch value := value.(type) {
		case string:
			if key == "spaceId" {
				return value
			}
		case map[string]interface{}:
			if spaceId := findSpaceId(value); spaceId != "" {
				return spaceId
			}
		case []interface{}:
			for _, item := range value {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if spaceId := findSpaceId(itemMap); spaceId != "" {
						return spaceId
					}
				}
			}
		}
	}

	return ""
}
