package cmd

import (
	"os"

	log "github.com/inconshreveable/log15"
	"github.com/ncsa/ssh-auditor/sshauditor"
	"github.com/spf13/cobra"
)

var store *sshauditor.SQLiteStore
var dbPath string
var debug bool

func initStore() error {
	//This should really return err, but it doesn't look as nice as when I fail immediately
	//cobra gives the help for the current command, which is irrelevant
	s, err := sshauditor.NewSQLiteStore(dbPath)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	err = s.Init()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	_, err = s.Begin()
	store = s
	return err
}

var RootCmd = &cobra.Command{
	Use:   "ssh-auditor",
	Short: "ssh-auditor tests ssh server password security",
	Long:  `Complete documentation is available at https://github.com/ncsa/ssh-auditor`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.Root().SetHandler(log.LvlFilterHandler(
				log.LvlDebug,
				log.StderrHandler))
		} else {
			log.Root().SetHandler(log.LvlFilterHandler(
				log.LvlInfo,
				log.StderrHandler))
		}
		return initStore()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return store.Commit()
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&dbPath, "db", "ssh_db.sqlite", "Path to database file")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug")
}
