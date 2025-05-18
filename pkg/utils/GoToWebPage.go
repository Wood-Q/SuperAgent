package utils

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
)

func GoToWebPage(ctx context.Context, url string) (string, error) {
	but, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		log.Fatal(err)
	}

	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		return "", err
	}
	return result.Output, nil
}
