package rr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/lixiangzhong/log"
)

type WatcherConfig struct {
	ExcludeDir []string `yaml:exclude_dir`
	Ext        []string `yaml:ext`
}

type Watcher struct {
	watcher      *fsnotify.Watcher
	cfg          WatcherConfig
	OnChangeFunc func()
}

func NewWatcher(cfg WatcherConfig) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{watcher: watcher, cfg: cfg, OnChangeFunc: NoneFunc}
	return w, err
}

func (w *Watcher) Start() {
	w.start()
}

func (w *Watcher) start() {
	for {
		select {
		case ev := <-w.watcher.Events:
			var err error
			if ev.Op&fsnotify.Remove == fsnotify.Remove {
				err = w.UnWatch(ev.Name)
				w.OnChangeFunc()
				break
			}
			info, err := os.Lstat(ev.Name)
			if err != nil {
				break
			}
			ext := filepath.Ext(ev.Name)
			if ev.Op&fsnotify.Create == fsnotify.Create {
				if info.IsDir() {
					err = w.Watch(ev.Name)
				} else {
					if StringInSlice(ext, w.cfg.Ext) {
						w.OnChangeFunc()
					}
				}
				break
			}
			if ev.Op&fsnotify.Write == fsnotify.Write {
				if StringInSlice(ext, w.cfg.Ext) {
					w.OnChangeFunc()
				}
			}
			if err != nil {
				log.Error(ev, err)
			}
		case err := <-w.watcher.Errors:
			log.Error("ev:", err)
		}
	}
}

func (w *Watcher) Watch(root string) error {

	for _, v := range w.cfg.ExcludeDir {
		if strings.HasPrefix(root, v) {
			return nil
		}
	}
	if isHidden(root) {
		return nil
	}
	w.add(root)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if root == path {
			return nil
		}
		if info.IsDir() {
			return w.Watch(path)
		}
		return nil
	})
}

func (w *Watcher) UnWatch(root string) error {
	w.remove(root)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if root == path {
			return nil
		}
		if info.IsDir() {
			w.UnWatch(path)
			return filepath.SkipDir
		}
		return nil
	})
}

func (w *Watcher) add(name string) {
	if StringInSlice(name, w.cfg.ExcludeDir) {
		return
	}
	fmt.Println("watch:", name)
	w.watcher.Add(name)
}

func (w *Watcher) remove(name string) {
	w.watcher.Remove(name)
}

func (w *Watcher) Close() {
	w.watcher.Close()
}
