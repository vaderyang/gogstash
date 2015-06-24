package outputstdout

import (
	log "github.com/Sirupsen/logrus"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/gogstash/config"
	"os"
)

const (
	ModuleName = "file"
)

type OutputConfig struct {
	config.CommonConfig
	Path      string `json:"path"`
	Is_Append bool   `json:"is_append"`
	fp        *os.File
}

func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		CommonConfig: config.CommonConfig{
			Type: ModuleName,
		},
	}
}

func create_or_append_file(fpath string, isappend bool) (fp *os.File, err error) {

	flag := os.O_CREATE | os.O_RDWR
	if isappend {
		log.Infof("Open file to write %q\n", fpath)
		flag = flag | os.O_APPEND
	}

	fp, err = os.OpenFile(fpath, flag, 0660)
	if err != nil {
		log.Errorf("Open file to write failed: %q\n%v", fpath, err)
		return
	}
	return fp, err

}

func init() {
	config.RegistOutputHandler(ModuleName, func(mapraw map[string]interface{}) (retconf config.TypeOutputConfig, err error) {
		conf := DefaultOutputConfig()

		if err = config.ReflectConfig(mapraw, &conf); err != nil {
			return
		}

		path := conf.Path
		isappend := conf.Is_Append
		conf.fp, err = create_or_append_file(path, isappend)

		if err != nil {
			err = errutil.New("Create file error in output", err)
			return
		}

		retconf = &conf
		return
	})

}

func (t *OutputConfig) Event(event config.LogEvent) (err error) {

	raw, err := event.MarshalIndent()
	if err != nil {
		return
	}

	_, err = t.fp.WriteString(string(raw) + "\n")

	if err != nil {

		return
	}
	t.fp.Sync()

	return
}
