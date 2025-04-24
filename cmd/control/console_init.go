package control

const (
	consoleDone = "/run/console-done"
	gettyCmd    = "/sbin/agetty"
	contosHome = "/home/contos"
	runLockDir  = "/run/lock"
	sshdFile    = "/etc/ssh/sshd_config"
)

type symlink struct {
	oldname, newname string
}