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
	vaultShowName = vaultShow.Arg("name", "Name of the vault to show").String()
	vaultList     = vault.Command("list", "List all vaults.")
	vaultUse      = vault.Command("use", "Select active vault.")
	vaultUseName  = vaultUse.Arg("name", "Name of the vault to use").Required().String()

	vaultAdd              = vault.Command("add", "")
	vaultAddUser          = vaultAdd.Command("user", "Add a user to the vault")
	vaultAddUserEmail     = vaultAddUser.Arg("email", "Email address, e.g. bob@foo.com").Required().String()
	vaultAddUserVaultName = vaultAddUser.Flag("vault", "(optional) name of vault to use").String()

	vaultRemove              = vault.Command("remove", "")
	vaultRemoveUser          = vaultRemove.Command("user", "Remove a user from the vault")
	vaultRemoveUserEmail     = vaultRemoveUser.Arg("email", "Email address, e.g. bob@foo.com, of the user to remove").Required().String()
	vaultRemoveUserVaultName = vaultRemoveUser.Flag("vault", "(optional) name of vault to use").String()

	vaultSetRemote     = vault.Command("remote", "Set or view remote url (git).")
	vaultSetRemoteUrl  = vaultSetRemote.Arg("url", "Remote url").Required().String()
	vaultSetRemoteName = vaultSetRemote.Flag("vault", "Name of the vault to set remote for.").String()
	// vaultList     = vault.Command("list", "List all vaults.")

	addSecret            = app.Command("store", "Store a secret.")
	addSecretName        = addSecret.Flag("vault", "Name of the vault.").String()
	addSecretSite        = addSecret.Flag("site", "The site is the container for the secrets.").Required().String()
	addSecretUsername    = addSecret.Flag("username", "Username associated with the site.").Required().String()
	addSecretPassword    = addSecret.Flag("passphrase", "Passphrase to store with the site.").Required().String()
	addSecretDescription = addSecret.Flag("description", "Description to store with the site.").String()

	revealSecret          = app.Command("reveal", "Reveal a secret.")
	revealSecretSite      = revealSecret.Arg("site", "Name of site to reveal.").Required().String()
	revealSecretVaultName = revealSecret.Flag("vault", "Name of the vault.").String()

	vaultSync     = vault.Command("sync", "Sync local vault with a remote vault.")
	vaultSyncName = vaultSync.Flag("vault", "(optional) Name of the vault to sync.").String()

	vaultFetch     = vault.Command("fetch", "Fetch a remote vault.")
	vaultFetchUrl  = vaultFetch.Arg("url", "Remote url, e.g. git@github.com/passward/test.git").Required().String()
	vaultFetchName = vaultFetch.Flag("name", "(optional) Name of vault to use.").String()
)

func Run() {
	app.Version(VERSION)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case setup.FullCommand():
		commands.Setup()

	case vault.FullCommand():
		println("Subcommand for `vault` is required.")
		app.CommandUsage(os.Stderr, vault.FullCommand())

	case revealSecret.FullCommand():
		commands.VaultSecretReveal(*revealSecretVaultName, *revealSecretSite)

	case vaultAddUser.FullCommand():
		commands.VaultAddUser(*vaultAddUserVaultName, *vaultAddUserEmail)

	case vaultNew.FullCommand():
		commands.VaultNew(*vaultNewName)

	case vaultUse.FullCommand():
		commands.VaultUse(*vaultUseName)

	case vaultSetRemote.FullCommand():
		commands.VaultSetRemote(*vaultSetRemoteName, *vaultSetRemoteUrl)

	case vaultFetch.FullCommand():
		commands.VaultFetch(*vaultFetchUrl, *vaultFetchName)

	case vaultSync.FullCommand():
		commands.VaultSync(*vaultSyncName)

	case vaultShow.FullCommand():
		commands.VaultShow(*vaultShowName)

	case vaultList.FullCommand():
		commands.VaultList()

	case addSecret.FullCommand():
		commands.VaultSecretAdd(*addSecretName, *addSecretSite, *addSecretUsername, *addSecretPassword, *addSecretDescription)

	default:
		app.Usage(os.Stderr)
	}
}
