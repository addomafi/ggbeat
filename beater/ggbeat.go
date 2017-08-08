package beater

import (
	"fmt"
	"time"
	"bytes"
	"bufio"
	"regexp"
	"strings"
	"strconv"
	"os/exec"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/addomafi/ggbeat/config"
)

type Ggbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Ggbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func convertToMinutes(str string) int {
	time := strings.Split(str, ":")
	hr, er := strconv.Atoi(time[0])
	if (er != nil) {
		return 0
	}
	min, er := strconv.Atoi(time[1])
	if (er != nil) {
		return 0
	}
	return hr*60 + min
}

func (bt *Ggbeat) Run(b *beat.Beat) error {
	logp.Info("ggbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		out, err := exec.Command("bash","-c","echo | cat teste.txt").Output()

		r := regexp.MustCompile(`^(\w+)\s+(\w+)\s+(\w+)\s+([\w0-9:]+)\s+([\w0-9:]+)$`)

		if (err == nil) {
			scanner := bufio.NewScanner(bytes.NewReader(out))
		  for scanner.Scan() {
				logp.Info("Find for regex: " + scanner.Text());
				matches := r.FindStringSubmatch(scanner.Text())
				if len(matches) == 6 {
					logp.Info("Adding info: " + strings.Join(matches,","))
					bt.client.PublishEvent(common.MapStr{
						"@timestamp": common.Time(time.Now()),
						"type":       matches[1],
						"status":     matches[2],
						"name":       matches[3],
						"value":       convertToMinutes(matches[4]) + convertToMinutes(matches[5]),
					})
				}
		  }
		}
	}
}

func (bt *Ggbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
