package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/blend/go-sdk/env"
	logger "github.com/blend/go-sdk/logger"
	web "github.com/blend/go-sdk/web"
	"github.com/mat285/advent-bot/pkg/advent"
	"github.com/mat285/advent-bot/pkg/slack"
)

const errMessage = "Oops! Something's not quite right"

func main() {
	log := logger.All()

	wc, err := web.NewConfigFromEnv()
	if err != nil {
		log.SyncFatalExit(err)
	}

	app := web.NewFromConfig(wc).WithLogger(log)

	app.POST("/", handle)

	quit := make(chan os.Signal, 1)
	// trap ^C
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		<-quit
		log.SyncError(app.Shutdown())
	}()

	if err := web.StartWithGracefulShutdown(app); err != nil {
		log.SyncFatalExit(err)
	}
}

func handle(r *web.Ctx) web.Result {
	body, err := r.PostBody()
	if err != nil {
		return r.JSON().InternalError(err)
	}
	r.Request().Body = ioutil.NopCloser(bytes.NewReader(body))
	// user := web.StringValue(r.Param(slack.ParamUserIDKey))
	// text := web.StringValue(r.Param(slack.ParamTextKey))

	// err = verify(r) // verify the request came from slack
	// if err != nil {
	// 	r.Logger().Error(err)
	// 	return r.JSON().NotAuthorized()
	// }

	resp, err := advent.GetLeaderBoard()
	if err != nil {
		return r.JSON().InternalError(err)
	}

	return r.JSON().Result(getSlackMessage(resp))
}

func verify(r *web.Ctx) error {
	timestamp, err := r.HeaderValue(slack.TimestampHeaderParam)
	if err != nil {
		return err
	}
	body, err := r.PostBody()
	if err != nil {
		return err
	}
	sig, err := r.HeaderValue(slack.SignatureHeaderParam)
	if err != nil {
		return err
	}
	return slack.VerifyRequest(timestamp, string(body), string(sig), env.Env().String(slack.EnvVarSignatureSecret))
}

func getSlackMessage(board *advent.Board) *slack.Message {
	return &slack.Message{
		ResponseType: slack.ResponseTypeInChannel,
		Text:         getMessageText(board),
	}
}

func getMessageText(board *advent.Board) string {
	str := "```"
	for i, m := range board.Members {
		str += fmt.Sprintf("%d. %s - %d\n", i+1, m.Name, m.LocalScore)
	}
	return strings.TrimSpace(str) + "```"
}
