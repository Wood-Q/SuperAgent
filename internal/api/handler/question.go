package handler

import (
	"MoonAgent/cmd/di"
	"MoonAgent/internal/pipeline"
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sse"
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

func (h *ChatHandler) StreamChatWithModel(ctx context.Context, c *app.RequestContext) {
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

	// 设置SSE响应头
	c.SetStatusCode(http.StatusOK)
	stream := sse.NewStream(c)

	// 创建带有用户输入的上下文
	ctx = context.WithValue(context.Background(), "user_input", req.UserInput)

	// 构建助手
	runnable, err := pipeline.BuildAssitant(ctx, h.app)
	if err != nil {
		// 发送错误事件
		errorEvent := &sse.Event{
			Event: "error",
			Data:  []byte(err.Error()),
		}
		stream.Publish(errorEvent)
		return
	}

	// 调用流式处理
	streamReader, err := runnable.Stream(ctx, req.UserInput)
	if err != nil {
		// 发送错误事件
		errorEvent := &sse.Event{
			Event: "error",
			Data:  []byte(err.Error()),
		}
		stream.Publish(errorEvent)
		return
	}

	// 从流中读取数据并发送给客户端
	for {
		chunk, err := streamReader.Recv()
		if err != nil {
			if err.Error() == "EOF" || err.Error() == "stream is finished" {
				// 正常结束
				break
			}
			// 发送错误事件
			errorEvent := &sse.Event{
				Event: "error",
				Data:  []byte(err.Error()),
			}
			stream.Publish(errorEvent)
			return
		}

		// 发送消息事件
		event := &sse.Event{
			Event: "message",
			Data:  []byte(chunk.Content),
		}

		if err := stream.Publish(event); err != nil {
			// 客户端断开连接
			break
		}
	}

	// 发送完成事件
	doneEvent := &sse.Event{
		Event: "done",
		Data:  []byte("Stream completed"),
	}
	stream.Publish(doneEvent)
}
