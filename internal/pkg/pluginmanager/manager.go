package pluginmanager

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

func DownloadPlugins(conf *configloader.Config) error {
	// create plugins dir if not exist
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf("plugins directory should not be \"\"")
	}
	log.Printf("Using dir <%s> to store plugins.", pluginDir)

	// download all plugins that don't exist locally
	dc := NewPbDownloadClient()

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); errors.Is(err, os.ErrNotExist) {
			// plugin does not exist
			err := dc.download(pluginDir, pluginFileName, tool.Version)
			if err != nil {
				return err
			}
			continue
		}
		// check md5
		dup, err := checkFileMD5(filepath.Join(pluginDir, pluginFileName), dc, pluginFileName, tool.Version)
		if err != nil {
			return err
		}

		if dup {
			log.Printf("Plugin: %s already exists, no need to download.", pluginFileName)
			continue
		}

		log.Printf("Plugin: %s changed and will be overrided.", pluginFileName)
		if err = os.Remove(filepath.Join(pluginDir, pluginFileName)); err != nil {
			return err
		}
		if err = dc.download(pluginDir, pluginFileName, tool.Version); err != nil {
			return err
		}
	}

	return nil
}

func checkFileMD5(file string, dc *PbDownloadClient, pluginFileName string, version string) (bool, error) {
	localmd5, err := LocalContentMD5(file)
	if err != nil {
		return false, err
	}
	remotemd5, err := dc.fetchContentMD5(pluginFileName, version)
	if err != nil {
		return false, err
	}

	if localmd5 == remotemd5 {
		return true, nil
	}
	return false, nil
}
