package redirectstat

import "sync"

type Click struct {
	RedirectKey string

	Direct uint64
	Social uint64
}

type Fail struct {
	RedirectKey string

	NotFound                   uint64
	DatabaseUnreachable        uint64
	TemplateProcessFailed      uint64
	ClientContentProcessFailed uint64
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
