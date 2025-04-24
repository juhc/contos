package control

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/juhc/contos/config"
	"github.com/urfave/cli"
)

const (
	socket = "/run/containerd/containerd.sock"
)

func entrypointAction(c *cli.Context) error {
	if _, err := os.Stat("/host/dev"); err == nil {
		cmd := exec.Command("mount", "--rbind", "/host/dev", "/dev")
		if err := cmd.Run(); err != nil {
			//log.Errorf("Failed to mount /dev: %v", err)
		}
	}

	shouldWriteFiles := false

	if shouldWriteFiles {
		//writeFiles(cfg)
	}

	setupCommandSymlinks()

	if len(os.Args) < 3 {
		return nil
	}

	binary, err := exec.LookPath(os.Args[2])
	if err != nil {
		return err
	}

	return syscall.Exec(binary, os.Args[2:], os.Environ())
}

func setupCommandSymlinks() {
	for _, link := range []symlink{
		{config.ContosCtlBin, "/usr/bin/autologin"},
		{config.ContosCtlBin, "/usr/bin/recovery"},
		{config.ContosCtlBin, "/usr/bin/respawn"},
		{config.ContosCtlBin, "/usr/sbin/netconf"},
		{config.ContosCtlBin, "/usr/sbin/wait-for-containerd"},
		{config.ContosCtlBin, "/usr/sbin/poweroff"},
		{config.ContosCtlBin, "/usr/sbin/reboot"},
		{config.ContosCtlBin, "/usr/sbin/halt"},
		{config.ContosCtlBin, "/usr/sbin/shutdown"},
		{config.ContosCtlBin, "/sbin/poweroff"},
		{config.ContosCtlBin, "/sbin/reboot"},
		{config.ContosCtlBin, "/sbin/halt"},
		{config.ContosCtlBin, "/sbin/shutdown"},
	} {
		os.Remove(link.newname)
		if err := os.Symlink(link.oldname, link.newname); err != nil {
			//log.Error(err)
		}
	}
}
