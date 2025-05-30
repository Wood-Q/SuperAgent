package handler

import (
	"MoonAgent/cmd/di"
	"MoonAgent/internal/pipeline"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type ChatHandler struct {
	app *di.Application
}

func NewChatHandler(app *di.Application) *ChatHandler {
	return &ChatHandler{app: app}
}

type Req struct {
	UserInput string `json:"userInput"`
}

func (h *ChatHandler) ChatWithModel(ctx context.Context, c *app.RequestContext) {
	var req Req

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if req.UserInput == "" {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "userInput cannot be empty",
		})
		return
	}

	ctx = context.WithValue(context.Background(), "user_input", req.UserInput)

	runnable, err := pipeline.BuildAssitant(ctx, h.app)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	out, err := runnable.Invoke(ctx, req.UserInput)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]string{
		"message": out.Content,
	})
}
