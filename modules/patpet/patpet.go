package patpet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"galched-bot/modules/settings"
)

type (
	Pet struct {
		sync.RWMutex
		path    string
		counter int
	}
)

func New(s *settings.Settings) (*Pet, error) {
	var (
		err     error
		counter int
	)

	petData, err := ioutil.ReadFile(s.PetDataPath)
	if err != nil {
		log.Print("pet: cannot read data file", err)
		log.Print("pet: creating new counter")
	} else {
		err = json.Unmarshal(petData, &counter)
		if err != nil {
			counter = 0
			log.Print("pet: cannot unmarshal data file", err)
			log.Print("pet: creating new counter")
		} else {
			log.Print("pet: using previously saved counter")
		}
	}

	return &Pet{
		RWMutex: sync.RWMutex{},
		path:    s.PetDataPath,
		counter: counter,
	}, nil
}

func (p *Pet) Pet() int {
	p.Lock()
	defer p.Unlock()

	p.counter++
	return p.counter
}

func (p *Pet) Counter() int {
	p.RUnlock()
	defer p.RUnlock()

	return p.counter
}

func (p *Pet) Dump() {
	p.RLock()
	defer p.RUnlock()

	data, err := json.Marshal(p.counter)
	if err != nil {
		log.Print("pet: cannot marshal counter", err)
		return
	}
	file, err := os.Create(p.path)
	if err != nil {
		log.Print("pet: cannot open counter file", err)
		return
	}
	_, err = fmt.Fprintf(file, string(data))
	if err != nil {
		log.Print("pet: cannot write to counter file")
	}
	err = file.Close()
	if err != nil {
		log.Print("pet: cannot close counter file")
	}
	log.Print("pet: counter dumped to file:", p.counter)
}
