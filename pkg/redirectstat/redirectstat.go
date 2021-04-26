package redirectstat

import "sync"

type Click struct {
	RedirectKey string `json:"-"`

	Direct uint64 `json:"direct"`
	Social uint64 `json:"social"`
}

type Fail struct {
	RedirectKey string `json:"-"`

	NotFound                   uint64 `json:"notFound"`
	DatabaseUnreachable        uint64 `json:"databaseUnreachable"`
	TemplateProcessFailed      uint64 `json:"templateProcessFailed"`
	ClientContentProcessFailed uint64 `json:"clientContentProcessFailed"`
}

type StatChannels struct {
	ClicksChannel chan *Click
	FailsChannel  chan *Fail
}

var StatChannelsInstance *StatChannels
var lock = &sync.Mutex{}

func GetStatChannels() *StatChannels {
	if StatChannelsInstance == nil {
		lock.Lock()
		if StatChannelsInstance == nil {
			StatChannelsInstance = &StatChannels{
				ClicksChannel: make(chan *Click, 1024),
				FailsChannel:  make(chan *Fail, 1024),
			}
		}
		lock.Unlock()
	}
	return StatChannelsInstance
}
