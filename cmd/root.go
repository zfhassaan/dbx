package cmd

import (
	"github.com/spf13/cobra"
)

/*
rootCmd defines the main Cobra command for the dbx CLI, providing database backup and restore functionality.
Execute runs the root command and handles any execution errors.

	Long:  "Cross-platform tool to backup and restore databases like MySQL, PostgreSQL, and SQLite",
*/
var rootCmd = &cobra.Command{
	Use:   "dbx",
	Short: "dbx is a simple database backup CLI utility",
	Long:  "Cross-platform tool to backup and restore databases like MySQL, PostgreSQL, MongoDB, and SQLite",
}

func PrintBanner() {
	banner := `

 ________  ________  ________  ___  __    ___  ___  ________        ___  ___  _________  ___  ___       ___  _________    ___    ___ 
|\   __  \|\   __  \|\   ____\|\  \|\  \ |\  \|\  \|\   __  \      |\  \|\  \|\___   ___\\  \|\  \     |\  \|\___   ___\ |\  \  /  /|
\ \  \|\ /\ \  \|\  \ \  \___|\ \  \/  /|\ \  \\\  \ \  \|\  \     \ \  \\\  \|___ \  \_\ \  \ \  \    \ \  \|___ \  \_| \ \  \/  / /
 \ \   __  \ \   __  \ \  \    \ \   ___  \ \  \\\  \ \   ____\     \ \  \\\  \   \ \  \ \ \  \ \  \    \ \  \   \ \  \   \ \    / / 
  \ \  \|\  \ \  \ \  \ \  \____\ \  \\ \  \ \  \\\  \ \  \___|      \ \  \\\  \   \ \  \ \ \  \ \  \____\ \  \   \ \  \   \/  /  /  
   \ \_______\ \__\ \__\ \_______\ \__\\ \__\ \_______\ \__\          \ \_______\   \ \__\ \ \__\ \_______\ \__\   \ \__\__/  / /    
    \|_______|\|__|\|__|\|_______|\|__| \|__|\|_______|\|__|           \|_______|    \|__|  \|__|\|_______|\|__|    \|__|\___/ /     
                                                                                                                        \|___|/

 Simple DB Backup Tool (dbx)
 ---------------------------
`
	println(banner)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
