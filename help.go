package main

var helpTemplate string = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options] <command> [<arguments...>]

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`
