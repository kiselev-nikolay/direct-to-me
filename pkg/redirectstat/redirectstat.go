package redirectstat

import "sync"

type Click struct {
	RedirectKey string

	Direct uint64
	Social uint64
}

type Error struct {
	RedirectKey string

	NotFound                   uint64
	DatabaseUnreachable        uint64
	TemplateProcessFailed      uint64
	ClientContentProcessFailed uint64
}

type StatChannel struct {
	ClicksChannel chan Click
	ErrorsChannel chan Error
}

var StatChannelInstance *StatChannel
var lock = &sync.Mutex{}

func GetStatChannel() *StatChannel {
	if StatChannelInstance == nil {
		lock.Lock()
		if StatChannelInstance == nil {
			StatChannelInstance = &StatChannel{}
		}
		lock.Unlock()
	}
	return StatChannelInstance
}
