package control

import (
	"fmt"
	"os"
	"os/exec"
	"strings"


	"github.com/urfave/cli"
)

func AutologinMain() {
	app := cli.NewApp()

	app.Name = os.Args[0]
	app.Usage = "autologin console"
	app.Version = "1.0"
	app.Author = "Vsevolod Chalkov"
	app.Email = "chalkov@centos.com"
	app.EnableBashCompletion = true
	app.Action = autologinAction
	app.HideHelp = true
	app.Run(os.Args)
}

func autologinAction(c *cli.Context) error {
	cmd := exec.Command("/bin/stty", "sane")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	usertty := ""
	user := "root"
	if c.NArg() > 0 {
		usertty = c.Args().Get(0)
		s := strings.SplitN(usertty, ":", 2)
		user = s[0]
		if len(s) > 1 {
		//	tty = s[1]
		}
	}

	banner := `
█▀▀ █▀█ █▄░█ ▀█▀ █▀█ █▀
█▄▄ █▄█ █░▀█ ░█░ █▄█ ▄█
Autologin ContOS v1
`
	fmt.Println(banner)

	loginBin := ""
	args := []string{}
	loginBin = "login"
	args = append(args, "-f", user)
	
	loginBinPath, err := exec.LookPath(loginBin)
	if err != nil {
		fmt.Printf("error finding %s in path: %s", cmd.Args[0], err)
		return err
	}
	os.Setenv("TERM", "linux")

	cmd = exec.Command(loginBinPath, args...)
	cmd.Env = os.Environ()

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return nil
}
