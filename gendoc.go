package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

func main() {
	srcFlag := cli.StringFlag{
		Name:  "src",
		Usage: "yaml files directory",
	}
	dstFlag := cli.StringFlag{
		Name:  "dst",
		Usage: "json directory",
	}
	templateFlag := cli.StringFlag{
		Name:  "template",
		Usage: "template directory",
		Value: getGoPath() + "/src/github.com/hiroosak/gendoc/template",
	}
	metaFlag := cli.StringFlag{
		Name:  "meta",
		Usage: "meta file",
	}
	overviewFlag := cli.StringFlag{
		Name:  "overview",
		Usage: "overview file",
	}

	app := cli.NewApp()
	app.Name = "gendoc"
	app.Usage = "make an document"
	app.Version = "0.1"

	app.Commands = []cli.Command{
		cli.Command{
			Name:   "init",
			Usage:  "Generate initialized json schema YAML",
			Action: scaffoldAction,
		},
		cli.Command{
			Name:   "doc",
			Usage:  "Generate html from json schema",
			Action: docAction,
			Flags:  []cli.Flag{srcFlag, templateFlag, metaFlag, overviewFlag},
		},
		cli.Command{
			Name:   "valid",
			Usage:  "Validation YAML or JSON file",
			Action: validAction,
			Flags:  []cli.Flag{srcFlag},
		},
		cli.Command{
			Name:   "gen",
			Usage:  "Generate JSON from YAML",
			Action: genAction,
			Flags:  []cli.Flag{srcFlag, dstFlag},
		},
	}
	app.Run(os.Args)
}

func scaffoldAction(c *cli.Context) {
	scaffold(c.Args().First())
}

func genAction(c *cli.Context) {
	src := c.String("src")
	dst := c.String("dst")
	if err := generateJSON(src, dst); err != nil {
		fmt.Println(err)
		fmt.Println("")
		cli.ShowAppHelp(c)
	} else {
		fmt.Println("ok.")
	}
}

func validAction(c *cli.Context) {
	src := c.String("src")
	if err := validSchemaTree(src); err != nil {
		fmt.Println(err)
		fmt.Println("")
	} else {
		fmt.Println("ok.")
	}
}

func docAction(c *cli.Context) {
	src := c.String("src")
	meta := c.String("meta")
	template := c.String("template")
	overview := c.String("overview")

	if err := generateHTML(src, meta, overview, template); err != nil {
		fmt.Println(err)
		fmt.Println("")
		cli.ShowAppHelp(c)
	}
}

func getGoPath() string {
	paths := strings.Split(os.Getenv("GOPATH"), ":")
	return paths[len(paths)-1]
}
