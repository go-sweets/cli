package gen

import (
	"github.com/go-sweets/cli/internal/gen/gorm"
	"github.com/go-sweets/cli/internal/gen/migrate"
	"github.com/spf13/cobra"
)

var CmdGen = &cobra.Command{
	Use:   "gen",
	Short: "gen: Generate Directory. gen gorm ",
	Long:  "gen: Generate Directory. gen gorm ",
}

func init() {
	CmdGen.AddCommand(gorm.CmdGorm)
	CmdGen.AddCommand(migrate.CmdMigrate)
}
