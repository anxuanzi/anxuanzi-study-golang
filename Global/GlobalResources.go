package Global

import (
	_ "az-ops/resources/statik"
	"github.com/pterm/pterm"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/afero"
	"net/http"
	"os"
)

var etcFolder afero.Fs
var tmpFolder afero.Fs
var statikFS http.FileSystem

func InitResources() {
	etcFolder = afero.NewBasePathFs(afero.NewOsFs(), "/etc")
	tmpFolder = afero.NewBasePathFs(afero.NewOsFs(), "/tmp")
	var err error
	statikFS, err = fs.New()
	if err != nil {
		pterm.Error.Println("An error has occurred while trying load the program.")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	initMirrors()
	initTimezones()
	initDNSList()
	initBasicSoftwareList()
}

func GetEtcFolder() afero.Fs {
	return etcFolder
}
func GetTmpFolder() afero.Fs {
	return tmpFolder
}

func GetStatikFS() http.FileSystem {
	return statikFS
}
