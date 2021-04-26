package redirectstat_test

import (
	"context"
	"testing"
	"time"

	"github.com/kiselev-nikolay/direct-to-me/pkg/redirectstat"
	"github.com/stretchr/testify/require"
)

func TestWorker(t *testing.T) {
	require := require.New(t)
	ra := redirectstat.RedirectAggregation{}
	ch := redirectstat.GetStatChannels()

	ch.ClicksChannel <- &redirectstat.Click{RedirectKey: "test", Direct: 1}
	ch.ClicksChannel <- &redirectstat.Click{RedirectKey: "test", Direct: 1}

	ch.FailsChannel <- &redirectstat.Fail{RedirectKey: "test", NotFound: 1}
	for i := 0; i < 9; i++ {
		ch.FailsChannel <- &redirectstat.Fail{RedirectKey: "test", DatabaseUnreachable: 1}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ra.Worker(ctx)
	<-ctx.Done()

	testClicks, ok := ra.Clicks["test"]
	require.True(ok)
	require.EqualValues(redirectstat.Click{RedirectKey: "test", Direct: 2}, *testClicks)

	testFails, ok := ra.Fails["test"]
	require.True(ok)
	require.EqualValues(redirectstat.Fail{RedirectKey: "test", NotFound: 1, DatabaseUnreachable: 9}, *testFails)
}
