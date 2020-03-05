package rr

import (
	"context"
	"os"

	"github.com/lixiangzhong/log"
)

func NewEngine(cfg EngineConfig) (*Engine, error) {
	w, err := NewWatcher(cfg.WatcherConfig)
	if err != nil {
		return nil, err
	}
	go w.Start()
	e := &Engine{
		Watcher:    w,
		cfg:        cfg,
		cancelFunc: nil,
	}
	w.OnChangeFunc = e.BuildRun
	return e, nil
}

type Engine struct {
	Watcher    *Watcher
	cfg        EngineConfig
	cancelFunc context.CancelFunc
}

type EngineConfig struct {
	BuildCmd string `yaml:"build_cmd"`
	RunCmd   string `yaml:"run_cmd"`
	WatcherConfig
}

func (e *Engine) gobuild(ctx context.Context) error {
	log.Println("build...")
	cmd := NewCommand(ctx, e.cfg.BuildCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (e *Engine) gorun(ctx context.Context) error {
	log.Println("run...")
	cmd := NewCommand(ctx, e.cfg.RunCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func (e *Engine) BuildRun() {
	ctx, cancel := context.WithCancel(context.Background())
	err := e.gobuild(ctx)
	if err != nil {
		log.Error(err)
		cancel()
		return
	}
	e.Cancel()
	e.cancelFunc = cancel
	err = e.gorun(ctx)
	if err != nil {
		cancel()
		log.Error(err)
	}
}

func (e *Engine) Cancel() {
	if e.cancelFunc != nil {
		e.cancelFunc()
	}
}

func (e *Engine) Watch(root string) error {
	return e.Watcher.Watch(root)
}

func (e *Engine) Stop() {
	e.Watcher.Close()
	e.Cancel()
}

func (e *Engine) Start() {
	e.BuildRun()
}
