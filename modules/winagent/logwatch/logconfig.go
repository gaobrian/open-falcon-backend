package logwatch
import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"errors"
	"github.com/hpcloud/tail"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/toolkits/file"
)

type LogConfig struct {
	Metric     string
	Timer      int64
	Agent      string
	WatchFiles []WatchFile `json:"files"`
	LogLevel   string
}

type resultFile struct {
	FileName string
	ModTime  time.Time
	LogTail  *tail.Tail
}

type WatchFile struct {
	Path       string
	Prefix     string
	Suffix     string
	Keywords   []keyWord
	PathIsFile bool
	ResultFile resultFile `json:"-"`
}

type keyWord struct {
	Exp      string
	Tag      string
	FixedExp string         `json:"-"`
	Regex    *regexp.Regexp `json:"-"`
}

const configFile = "log.json"

var (
	logCfg         *LogConfig
	fixExpRegex = regexp.MustCompile(`[\W]+`)
)

func init() {

	var err error
	logCfg, err = ReadConfig(configFile)
	if err != nil {
		log.Fatal("ERROR: ", err)
		return
	}
	if err = checkLogConfig(logCfg); err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		LogConfigFileWatcher()
	}()
}

func ReadConfig(configFile string) (*LogConfig, error) {
	if !file.IsExist(configFile) {
		log.Fatalln("config file:", configFile, "is not existent. please create it")
		return nil, errors.New("File not exists")
	}

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config *LogConfig
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	fmt.Println(config.LogLevel)

	// 检查配置项目
	if err := checkLogConfig(config); err != nil {
		return nil, err
	}

	log.Println("config init success, start to work ...")
	return config, nil
}

// 检查配置项目是否正确
func checkLogConfig(config *LogConfig) error {

	for i, v := range config.WatchFiles {
		//检查路径
		fInfo, err := os.Stat(v.Path)
		if err != nil {
			return err
		}

		if !fInfo.IsDir() {
			config.WatchFiles[i].PathIsFile = true
		}

		//检查后缀,如果没有,则默认为.log
		config.WatchFiles[i].Prefix = strings.TrimSpace(v.Prefix)
		config.WatchFiles[i].Suffix = strings.TrimSpace(v.Suffix)
		//if config.WatchFiles[i].Suffix == "" {
		//	log.Println("file pre ", config.WatchFiles[i].Path, "suffix is no set, will use .log")
		//	config.WatchFiles[i].Suffix = ".log"
		//}

		//agent不检查,可能后启动agent

		//检查keywords
		if len(v.Keywords) == 0 {
			return errors.New("ERROR: keyword list not set")
		}

		for _, keyword := range v.Keywords {
			if keyword.Exp == "" || keyword.Tag == "" {
				return errors.New("ERROR: keyword's exp and tag are requierd")
			}
		}

		// 设置正则表达式
		for j, keyword := range v.Keywords {

			if config.WatchFiles[i].Keywords[j].Regex, err = regexp.Compile(keyword.Exp); err != nil {
				return err
			}

			log.Println("INFO: tag:", keyword.Tag, "regex", config.WatchFiles[i].Keywords[j].Regex.String())

			config.WatchFiles[i].Keywords[j].FixedExp = string(fixExpRegex.ReplaceAll([]byte(keyword.Exp), []byte(".")))
		}
	}

	return nil
}

//配置文件监控,可以实现热更新
func LogConfigFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Name == configFile && (event.Op == fsnotify.Chmod || event.Op == fsnotify.Rename || event.Op == fsnotify.Write || event.Op == fsnotify.Create) {
					log.Println("modified config file", event.Name, "will reaload config")
					if cfg, err := ReadConfig(configFile); err != nil {
						log.Println("ERROR: config has error, will not use old config", err)
					} else if checkLogConfig(logCfg) != nil {
						log.Println("ERROR: config has error, will not use old config", err)
					} else {
						log.Println("config reload success")
						logCfg = cfg
					}

				}
			case err := <-watcher.Errors:
				log.Fatal(err)
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
