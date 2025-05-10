package tools

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino-ext/components/tool/googlesearch"
	"go.uber.org/zap"
)

func GoogleSearch(query string) []*googlesearch.SimplifiedSearchItem {
	ctx := context.Background()

	googleAPIKey := "AIzaSyDrKLkc290NdtNfC8fkOVQTVCPq_yuXZpA"
	googleSearchEngineID := "60af7f26ff64c4d55"

	if googleAPIKey == "" || googleSearchEngineID == "" {
		zap.L().Error("GStool", zap.String("error", "GOOGLE_API_KEY and GOOGLE_SEARCH_ENGINE_ID must set"))
	}

	// create tool
	searchTool, err := googlesearch.NewTool(ctx, &googlesearch.Config{
		APIKey:         googleAPIKey,
		SearchEngineID: googleSearchEngineID,
		Lang:           "zh-CN",
		Num:            5,
	})
	if err != nil {
		zap.L().Error("GStool", zap.String("error", err.Error()))
	}

	// prepare params
	req := googlesearch.SearchRequest{
		Query: query,
		Num:   3,
		Lang:  "zh-CN",
	}

	args, err := json.Marshal(req)
	if err != nil {
		zap.L().Error("GStool", zap.String("error", err.Error()))
	}

	// do search
	resp, err := searchTool.InvokableRun(ctx, string(args))
	if err != nil {
		zap.L().Error("GStool", zap.String("error", err.Error()))
	}

	var searchResp googlesearch.SearchResult
	if err := json.Unmarshal([]byte(resp), &searchResp); err != nil {
		zap.L().Error("GStool", zap.String("error", err.Error()))
	}

	return searchResp.Items
}
