package beater

import (
	"fmt"
	"sync"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/fupslot/fsstatbeat/config"
)

// fsstatbeat configuration.
type fsstatbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	b      *beat.Beat
}

// New creates an instance of fsstatbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &fsstatbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts fsstatbeat.
func (bt *fsstatbeat) Run(b *beat.Beat) error {
	logp.Info("fsstatbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	bt.b = b

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		var wg sync.WaitGroup
		bt.RunChecks(&wg)
		wg.Wait()
	}
}

func (bt *fsstatbeat) RunChecks(wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() { wg.Done() }()

	for _, r := range bt.config.Rules {
		bt.Check(&r)
	}
}

func (bt *fsstatbeat) Check(r *config.Rule) {
	for _, res := range r.Resources {
		if res.File.Path != "" {
			stat, err := Fsstat(res)
			if err != nil {
				logp.Err(err.Error())
				continue
			}

			bt.PublishEventFile(stat)
		} else if res.Proc.Name != "" {
			logp.Info("!!!%s", res.Proc.Name)
		}
	}
}

func (bt *fsstatbeat) PublishEventFile(st *FileState) {
	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type": bt.b.Info.Name,
			"file": common.MapStr{
				"name":  st.name,
				"path":  st.path,
				"umask": st.umask,
				"owner": st.owner,
				"perm":  st.perm,
				"group": st.group,
				"octal": st.octal,
			},
		},
	}
	bt.client.Publish(event)
	logp.Info("Event sent")
}

// Stop stops fsstatbeat.
func (bt *fsstatbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
