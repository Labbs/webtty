package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/labbs/webtty/backend/localcommand"
	"github.com/labbs/webtty/server"
)

var appOptions *server.Options = &server.Options{}
var backendOptions *localcommand.Options = &localcommand.Options{}
var Version string = "unknown_version"
var CommitID string = "unknown_commit"

func main() {
	app := cli.NewApp()
	app.Name = "gotty"
	app.Version = Version + "+" + CommitID
	app.Usage = "Share your terminal as a web application"
	app.HideHelp = true
	cli.AppHelpTemplate = helpTemplate
	app.Flags = flags()
	app.Before = altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewJSONSourceFromFlagFunc("config"))
	app.Action = action
	app.Run(os.Args)
}

func waitSignals(errs chan error, cancel context.CancelFunc, gracefullCancel context.CancelFunc) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-errs:
		return err

	case s := <-sigChan:
		switch s {
		case syscall.SIGINT:
			gracefullCancel()
			fmt.Println("C-C to force close")
			select {
			case err := <-errs:
				return err
			case <-sigChan:
				fmt.Println("Force closing...")
				cancel()
				return <-errs
			}
		default:
			cancel()
			return <-errs
		}
	}
}

func action(c *cli.Context) error {
	if c.Args().Len() == 0 {
		msg := "Error: No command given."
		cli.ShowAppHelp(c)
		return fmt.Errorf(msg)
	}

	args := c.Args()
	factory, err := localcommand.NewFactory(args.First(), args.Slice()[1:], backendOptions)
	if err != nil {
		return err
	}

	hostname, _ := os.Hostname()
	appOptions.TitleVariables = map[string]interface{}{
		"command":  args.First(),
		"argv":     args.Slice()[1:],
		"hostname": hostname,
	}

	srv, err := server.New(factory, appOptions)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	gCtx, gCancel := context.WithCancel(context.Background())

	log.Printf("GoTTY is starting with command: %s", strings.Join(args.Slice(), " "))

	errs := make(chan error, 1)
	go func() {
		errs <- srv.Run(ctx, server.WithGracefullContext(gCtx))
	}()
	err = waitSignals(errs, cancel, gCancel)

	if err != nil && err != context.Canceled {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	return nil
}
