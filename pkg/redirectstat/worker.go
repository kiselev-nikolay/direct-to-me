package redirectstat

import (
	"context"
	"sync"
	"sync/atomic"
)

type RedirectAggregation struct {
	mutex  sync.Mutex
	Clicks map[string]*Click
	Fails  map[string]*Fail
}

func (rs *RedirectAggregation) Worker(ctx context.Context) {
	statChannels := GetStatChannels()
	select {
	default:
		rs.mutex.Lock()
		if rs.Clicks == nil {
			rs.Clicks = make(map[string]*Click)
		}
		if rs.Fails == nil {
			rs.Fails = make(map[string]*Fail)
		}
		rs.mutex.Unlock()
		go rs.clickWaiter(ctx, statChannels.ClicksChannel)
		go rs.failWaiter(ctx, statChannels.FailsChannel)
	case <-ctx.Done():
		return
	}
}

func (rs *RedirectAggregation) clickWaiter(ctx context.Context, ch chan *Click) {
	select {
	default:
		for {
			click := <-ch
			rs.mutex.Lock()
			v, ok := rs.Clicks[click.RedirectKey]
			if !ok {
				rs.Clicks[click.RedirectKey] = click
				rs.mutex.Unlock()
				continue
			}
			rs.mutex.Unlock()
			atomic.AddUint64(&v.Direct, click.Direct)
			atomic.AddUint64(&v.Social, click.Social)
		}
	case <-ctx.Done():
		return
	}
}

func (rs *RedirectAggregation) failWaiter(ctx context.Context, ch chan *Fail) {
	select {
	default:
		for {
			fail := <-ch
			rs.mutex.Lock()
			v, ok := rs.Fails[fail.RedirectKey]
			if !ok {
				rs.Fails[fail.RedirectKey] = fail
				rs.mutex.Unlock()
				continue
			}
			rs.mutex.Unlock()
			atomic.AddUint64(&v.NotFound, fail.NotFound)
			atomic.AddUint64(&v.DatabaseUnreachable, fail.DatabaseUnreachable)
			atomic.AddUint64(&v.TemplateProcessFailed, fail.TemplateProcessFailed)
			atomic.AddUint64(&v.ClientContentProcessFailed, fail.ClientContentProcessFailed)
		}
	case <-ctx.Done():
		return
	}
}
