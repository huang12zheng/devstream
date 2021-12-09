package plugin

import (
	"fmt"
	"os"
	"plugin"

	"github.com/merico-dev/stream/internal/pkg/config"
	"github.com/merico-dev/stream/internal/pkg/download"
)

// DevStreamPlugin is a struct, on which install/reinstall/uninstall interfaces are defined.
type DevStreamPlugin interface {
	Install(*map[string]interface{})
	Reinstall(*map[string]interface{})
	Uninstall(*map[string]interface{})
}

func loadPlugin(tool *config.Tool) DevStreamPlugin {
	mod := fmt.Sprintf("plugins/%s_%s.so", tool.Name, tool.Version)
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var devStreamPlugin DevStreamPlugin
	symDevStreamPlugin, err := plug.Lookup("DevStreamPlugin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	devStreamPlugin, ok := symDevStreamPlugin.(DevStreamPlugin)
	if !ok {
		fmt.Println(err)
		os.Exit(1)
	}

	return devStreamPlugin
}

// Install loads the plugin and calls the Install method of that plugin.
func Install(tool *config.Tool) {

	path := fmt.Sprintf("plugins/%s_%s.so", tool.Name, tool.Version)
	appname := fmt.Sprintf("%s_%s.so", tool.Name, tool.Version)
	if !FileExist(path) {
		loader := download.NewDownloadClient()
		loader.AssetName = appname
		loader.Version = tool.Version
		loader.Filepath = path
		err := loader.GetAssetswithretry()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	p := loadPlugin(tool)
	p.Install(&tool.Options)
}

// Reinstall loads the plugin and calls the Reinstall method of that plugin.
func Reinstall(tool *config.Tool) {
	p := loadPlugin(tool)
	p.Reinstall(&tool.Options)
}

// Uninstall loads the plugin and calls the Uninstall method of that plugin.
func Uninstall(tool *config.Tool) {
	p := loadPlugin(tool)
	p.Uninstall(&tool.Options)
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
