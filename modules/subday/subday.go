package subday

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"galched-bot/modules/settings"
)

type Subday struct {
	sync.RWMutex
	path     string
	database map[string][]string
}

func New(s *settings.Settings) (*Subday, error) {
	var (
		err  error
		data = make(map[string][]string)
	)

	subdayData, err := ioutil.ReadFile(s.SubdayDataPath)
	if err != nil {
		log.Print("subday: cannot read subday data file", err)
		log.Print("subday: creating new subday database")
	} else {
		err = json.Unmarshal(subdayData, &data)
		if err != nil {
			data = make(map[string][]string)
			log.Print("subday: cannot unmarshal subday data file", err)
			log.Print("subday: creating new subday database")
		} else {
			log.Print("subday: using previously saved subday database")
		}
	}

	subday := &Subday{
		RWMutex:  sync.RWMutex{},
		path:     s.SubdayDataPath,
		database: data,
	}
	return subday, nil
}

func (s *Subday) Database() map[string][]string {
	return s.database
}

func (s *Subday) DumpToFile() {
	data, err := json.Marshal(s.database)
	if err != nil {
		log.Print("subday: cannot marshal database file", err)
		return
	}
	file, err := os.Create(s.path)
	if err != nil {
		log.Print("subday: cannot open database file", err)
		return
	}
	_, err = fmt.Fprintf(file, string(data))
	if err != nil {
		log.Print("subday: cannot write to database file")
	}
	err = file.Close()
	if err != nil {
		log.Print("subday: cannot close database file")
	}
	log.Print("subday: database dumped to file")
}
