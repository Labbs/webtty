package main

import (
	"github.com/labbs/webtty/backend/localcommand"
	"github.com/labbs/webtty/utils"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			EnvVars:     []string{"CONFIG"},
			Usage:       "Config file path",
			Value:       "config.json",
			Destination: &appOptions.ConfigFile,
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "address",
			Usage:       "Address to listen on",
			EnvVars:     []string{"ADDRESS"},
			Value:       "0.0.0.0",
			Destination: &appOptions.Address,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "port",
			Usage:       "Port to listen on",
			EnvVars:     []string{"PORT"},
			Value:       "8080",
			Destination: &appOptions.Port,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "path",
			Usage:       "URL path to access the web terminal",
			Value:       "/",
			Destination: &appOptions.Path,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "permit-write",
			Usage:       "Permit clients to write to the TTY (BE CAREFUL)",
			EnvVars:     []string{"PERMIT_WRITE"},
			Value:       true,
			Destination: &appOptions.PermitWrite,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-basic-auth",
			Usage:       "Enable basic authentication",
			EnvVars:     []string{"ENABLE_BASIC_AUTH"},
			Value:       false,
			Destination: &appOptions.EnableBasicAuth,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "credential",
			Usage:       "Credential for basic authentication",
			EnvVars:     []string{"CREDENTIAL"},
			Value:       "",
			Destination: &appOptions.Credential,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-random-url",
			Usage:       "Enable random URL generation",
			EnvVars:     []string{"ENABLE_RANDOM_URL"},
			Value:       false,
			Destination: &appOptions.EnableRandomUrl,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "random-url-length",
			Usage:       "Length of the random URL",
			EnvVars:     []string{"RANDOM_URL_LENGTH"},
			Value:       8,
			Destination: &appOptions.RandomUrlLength,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-tls",
			Usage:       "Enable TLS",
			EnvVars:     []string{"ENABLE_TLS"},
			Value:       false,
			Destination: &appOptions.EnableTLS,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "tls-crt-file",
			Usage:       "Path to the TLS certificate file",
			EnvVars:     []string{"TLS_CRT_FILE"},
			Value:       "~/.gotty.crt",
			Destination: &appOptions.TLSCrtFile,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "tls-key-file",
			Usage:       "Path to the TLS key file",
			EnvVars:     []string{"TLS_KEY_FILE"},
			Value:       "~/.gotty.key",
			Destination: &appOptions.TLSKeyFile,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-tls-client-auth",
			Usage:       "Enable TLS client authentication",
			EnvVars:     []string{"ENABLE_TLS_CLIENT_AUTH"},
			Value:       false,
			Destination: &appOptions.EnableTLSClientAuth,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "tls-ca-crt-file",
			Usage:       "Path to the TLS CA certificate file",
			EnvVars:     []string{"TLS_CA_CRT_FILE"},
			Value:       "~/.gotty.ca.crt",
			Destination: &appOptions.TLSCACrtFile,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "index-file",
			Usage:       "Path to the index file",
			EnvVars:     []string{"INDEX_FILE"},
			Value:       "",
			Destination: &appOptions.IndexFile,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "title-format",
			Usage:       "Format string for the page title",
			EnvVars:     []string{"TITLE_FORMAT"},
			Value:       "{{ .command }}@{{ .hostname }}",
			Destination: &appOptions.TitleFormat,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "enable-reconnect",
			Usage:       "Enable automatic reconnection",
			EnvVars:     []string{"ENABLE_RECONNECT"},
			Value:       false,
			Destination: &appOptions.EnableReconnect,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "reconnect-time",
			Usage:       "Time in seconds between reconnection attempts",
			EnvVars:     []string{"RECONNECT_TIME"},
			Value:       10,
			Destination: &appOptions.ReconnectTime,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "max-connection",
			Usage:       "Maximum number of connections",
			EnvVars:     []string{"MAX_CONNECTION"},
			Value:       0,
			Destination: &appOptions.MaxConnection,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "once",
			Usage:       "Accept only one client and exit on disconnection",
			EnvVars:     []string{"ONCE"},
			Value:       false,
			Destination: &appOptions.Once,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "timeout",
			Usage:       "Timeout for connection",
			EnvVars:     []string{"TIMEOUT"},
			Value:       0,
			Destination: &appOptions.Timeout,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "permit-arguments",
			Usage:       "Permit clients to send command line arguments in URL",
			EnvVars:     []string{"PERMIT_ARGUMENTS"},
			Value:       true,
			Destination: &appOptions.PermitArguments,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "width",
			Usage:       "Width of the screen",
			EnvVars:     []string{"WIDTH"},
			Value:       0,
			Destination: &appOptions.Width,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "height",
			Usage:       "Height of the screen",
			EnvVars:     []string{"HEIGHT"},
			Value:       0,
			Destination: &appOptions.Width,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "ws-origin",
			Usage:       "Origin of WebSocket connection",
			EnvVars:     []string{"WS_ORIGIN"},
			Value:       "",
			Destination: &appOptions.WSOrigin,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "term",
			Usage:       "Terminal type [xterm]",
			EnvVars:     []string{"TERM"},
			Value:       "xterm",
			Destination: &appOptions.Term,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "close-signal",
			Usage:       "Signal sent to the command process when gotty close it",
			EnvVars:     []string{"CLOSE_SIGNAL"},
			Value:       1,
			Destination: &backendOptions.CloseSignal,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "close-timeout",
			Usage:       "Time in seconds to force kill process after client is disconnected",
			EnvVars:     []string{"CLOSE_TIMEOUT"},
			Value:       -1,
			Destination: &backendOptions.CloseTimeout,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "recording-url",
			Usage:       "URL to send recording data",
			EnvVars:     []string{"RECORDING_URL"},
			Value:       "",
			Destination: &utils.RecordingUrl,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "recording-enabled",
			Usage:       "Enable recording",
			EnvVars:     []string{"RECORDING_ENABLED"},
			Value:       true,
			Destination: &utils.RecordingEnabled,
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:        "blacklist",
			Usage:       "Blacklist of word to be censored",
			EnvVars:     []string{"BLACKLIST"},
			Destination: &localcommand.BlackList,
		}),
	}
}
