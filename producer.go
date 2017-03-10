package supermq

import (
	"fmt"
	nsqlib "github.com/nsqio/go-nsq"
)

type MProducer struct {
	poller chan *nsqlib.Producer
	pdlist []*nsqlib.Producer
}

func NewMProducer() *MProducer {
	return &MProducer{
		poller: make(chan *nsqlib.Producer, 10),
		pdlist: make([]*nsqlib.Producer, 0, 10),
	}
}

func (mp *MProducer) SetLogger(l logger, ll int) {
	for _, p := range mp.pdlist {
		p.SetLogger(l, nsqlib.LogLevel(ll))
	}
}

func (mp *MProducer) Stop() {
	for _, p := range mp.pdlist {
		p.Stop()
	}
}

func (m *MProducer) AddRx(addr string) error {
	config := nsqlib.NewConfig()
	w, err := nsqlib.NewProducer(addr, config)
	if err != nil {
		return err
	}
	if err := w.Ping(); err != nil {
		return err
	}
	m.pdlist = append(m.pdlist, w)
	m.poller <- w
	return nil
}

func (m *MProducer) MultiPublish(topic string, body [][]byte) (err error) {
	for i := 0; i < len(m.poller); i++ {
		onePder := <-m.poller
		err = onePder.MultiPublish(topic, body)
		m.poller <- onePder
		if err == nil {
			return
		} else {
			fmt.Printf("publish to %s failed, try next\n", onePder.String())
		}
	}
	if err != nil {
		return
	}
	return fmt.Errorf("no available producer")
}
