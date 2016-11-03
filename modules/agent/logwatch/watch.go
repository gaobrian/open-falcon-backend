package logwatch

import (
	cmap "github.com/streamrail/concurrent-map"
	"github.com/fsnotify/fsnotify"
	log "github.com/Sirupsen/logrus"
	"os"
	"github.com/hpcloud/tail"
	"path/filepath"
	"strings"
	"path"
	"time"
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"runtime"
)

var (
  workers chan bool
  keywords cmap.ConcurrentMap
)

func Start()  {
	workers = make(chan bool, runtime.NumCPU()*2)
	keywords = cmap.New()

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(int64(logCfg.Timer)))
		for range ticker.C {
			fillData()
			postData()
		}
	}()

	go func() {
		setLogFile()
		for i := 0; i < len(logCfg.WatchFiles); i++ {
			readFileAndSetTail(&(logCfg.WatchFiles[i]))
			go logFileWatcher(&(logCfg.WatchFiles[i]))
		}

	}()
}

func logFileWatcher(file *WatchFile) {
	logTail := file.ResultFile.LogTail
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
				if file.PathIsFile && event.Op == fsnotify.Create && event.Name == file.Path {
					log.Info("continue to watch file:", event.Name)
					if file.ResultFile.LogTail != nil {
						logTail.Stop()
						file.ResultFile.LogTail = nil
					}
					readFileAndSetTail(file)
				} else {

					if file.ResultFile.FileName == event.Name && (event.Op == fsnotify.Remove || event.Op == fsnotify.Rename) {
						log.Warn(event, "stop to tail")
					} else if event.Op == fsnotify.Create {
						//log.Infof("created file %v, basePath:%v", event.Name, path.Base(event.Name))
						if strings.HasSuffix(event.Name, file.Suffix) && strings.HasPrefix(path.Base(event.Name), file.Prefix) {
							if logTail != nil {
								logTail.Stop()
								file.ResultFile.LogTail = nil
							}
							file.ResultFile.FileName = event.Name
							readFileAndSetTail(file)
						}
					}
				}

			case err := <-watcher.Errors:
				log.Error(err)
			}
		}
	}()

	watchPath := file.Path
	if file.PathIsFile {
		watchPath = filepath.Dir(file.Path)
	}
	err = watcher.Add(watchPath)
	if err != nil {
		log.Fatal(err)

	}
	<-done
}

func readFileAndSetTail(file *WatchFile) {
	if file.ResultFile.FileName == "" {
		return
	}

	_, err := os.Stat(file.ResultFile.FileName)
	if err != nil {
		log.Fatal(file.ResultFile.FileName, err)
		return
	}

	log.Println("read file", file.ResultFile.FileName)
	tail, err := tail.TailFile(file.ResultFile.FileName, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	file.ResultFile.LogTail = tail

	go func() {
		for line := range tail.Lines {
		    log.Debugln(line.Text)
		    handleKeywords(*file, line.Text)
		}
	}()
}

func setLogFile() {
	for i, v := range logCfg.WatchFiles {
		log.Println(v.Path,v.Prefix,v.Suffix)
		if v.PathIsFile {
			logCfg.WatchFiles[i].ResultFile.FileName = v.Path
			continue
		}

		filepath.Walk(v.Path, func(path string, info os.FileInfo, err error) error {
			cfgPath := v.Path
			if strings.HasSuffix(cfgPath, "/") {
				cfgPath = string([]rune(cfgPath)[:len(cfgPath)-1])
			}

			if (filepath.Dir(path) != cfgPath) || (info.IsDir() == true) {
				return nil
			}
			if strings.HasPrefix(filepath.Base(path), v.Prefix) && strings.HasSuffix(path, v.Suffix) && !info.IsDir() {
				if logCfg.WatchFiles[i].ResultFile.FileName == "" || info.ModTime().After(logCfg.WatchFiles[i].ResultFile.ModTime) {
					logCfg.WatchFiles[i].ResultFile.FileName = path
					logCfg.WatchFiles[i].ResultFile.ModTime = info.ModTime()
				}
				return nil
			}

			return nil
		})

	}
}

// 查找关键词
func handleKeywords(file WatchFile, line string) {
	host,_ := g.Hostname()
	for _, p := range file.Keywords {
		value := 0.0
		if p.Regex.MatchString(line) {
			value = 1.0
		}

		key := file.ResultFile.FileName + p.Exp

		var data model.MetricValue
		if v, ok := keywords.Get(key); ok {
			d := v.(*model.MetricValue)
			tmpValue := value + d.Value.(float64)
			d.Value = tmpValue
			data = *d
		} else {
			data = model.MetricValue{Metric: logCfg.Metric,
				Endpoint:    host,
				Timestamp:   time.Now().Unix(),
				Value:       value,
				Step:        logCfg.Timer,
				Type:        "GAUGE",
				Tags:        "prefix=" + file.Prefix + ",suffix=" + file.Suffix + "," + p.Tag + "=" + p.FixedExp,
			}
		}
		keywords.Set(key, &data)

	}
}

func postData() {

	workers <- true
	go func() {
		if len(keywords.Items()) != 0 {
			metrics := make([]*model.MetricValue, 0, 20)
			for k, v := range keywords.Items() {
				metrics = append(metrics, v.(*model.MetricValue))
				log.Println(v)
				keywords.Remove(k)
			}
			g.SendToTransfer(metrics)
		}

		<-workers
	}()

}

func fillData() {
	host,_ := g.Hostname()
	for _, v := range logCfg.WatchFiles {
		for _, p := range v.Keywords {

			key := v.ResultFile.FileName + p.Tag

			if _, ok := keywords.Get(key); ok {
				continue
			}

			data := model.MetricValue{Metric: logCfg.Metric,
				Endpoint:    host,
				Timestamp:   time.Now().Unix(),
				Value:       0.0,
				Step:        logCfg.Timer,
				Type:        "GAUGE",
				Tags:        "prefix=" + v.Prefix + ",suffix=" + v.Suffix + "," + p.Tag + "=" + p.FixedExp,
			}
			keywords.Set(key, &data)
		}
	}

}








