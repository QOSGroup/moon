package main

import (
	"github.com/QOSGroup/qmoon/cmd/commands"
)

func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		commands.InitFilesCmd,
		commands.MigrationCmd,
		commands.NodeCmd,
		commands.TxCmd,
		commands.ServerCmd,
		commands.VersionCmd)

	//cmd := cli.PrepareBaseCmd(rootCmd, "", os.ExpandEnv(filepath.Join("$HOME", cfg.DefaultQMoonDir)))
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
