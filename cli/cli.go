package cli

import (
	"os"

	"github.com/jandre/passward/commands"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

const VERSION = "0.1.0"
const NAME = "passward"

var (
	app   = kingpin.New("passward", "Securely store and share passwords.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()
	setup = app.Command("setup", "Setup passward environment.")

	// vault new
	vault         = app.Command("vault", "Create and manage vaults.")
	vaultNew      = vault.Command("new", "Create a new vault.")
	vaultNewName  = vaultNew.Arg("name", "Name of the vault to create").Required().String()
	vaultShow     = vault.Command("show", "Show vault.")
	vaultShowName = vaultShow.Arg("name", "Name of the vault to show").Required().String()
	vaultList     = vault.Command("list", "List all vaults.")

	vaultPull = vault.Command("clone", "Clone a remote vault.")
	vaultSync = vault.Command("sync", "Sync local vault with a remote vault.")
)

func Run() {
	app.Version(VERSION)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case setup.FullCommand():
		commands.Setup()

	case vault.FullCommand():
		println("Subcommand for `vault` is required.")
		app.CommandUsage(os.Stderr, vault.FullCommand())

	case vaultNew.FullCommand():
		commands.VaultNew(*vaultNewName)

	case vaultShow.FullCommand():
		commands.VaultShow(*vaultShowName)

	case vaultList.FullCommand():
		commands.VaultList()

	default:
		app.Usage(os.Stderr)
	}
}
