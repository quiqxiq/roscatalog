package rosrun

import (
	"context"
	"errors"
	"time"

	"github.com/quiqxiq/goros/v4"

	ros "example.com/roscatalog"
)

// WatchAuto seperti Watch tetapi bila stream ditolak router sebelum ada data
// satu pun (mis. listen/follow tidak didukung di menu itu), otomatis beralih
// ke poll dengan jeda yang sama (Plan.Every, atau interval bawaan view, atau 1s).
// Poll dari view listen/follow memakai diff, jadi bentuk event tetap serupa.
func (r *Runner) WatchAuto(ctx context.Context, res *ros.Resource, view string, ver ros.Version, plan ros.Plan, p ros.Params) (<-chan Event, error) {
	if plan.Delivery == "" {
		plan.Delivery = ros.DeliveryStream
	}
	req, err := res.Build(view, ver, plan, p)
	if err != nil {
		return nil, err
	}
	first, err := r.Watch(ctx, req)
	if err != nil {
		return nil, err
	}
	if req.Delivery != ros.DeliveryStream {
		return first, nil
	}

	out := make(chan Event, 64)
	forward := func(ch <-chan Event) {
		for ev := range ch {
			select {
			case out <- ev:
			case <-ctx.Done():
				return
			}
		}
	}
	go func() {
		defer close(out)
		got := false
		for ev := range first {
			if ev.Kind == KindError && !got {
				var de *routeros.DeviceError
				if errors.As(ev.Err, &de) {
					every := plan.Every
					if v, ok := res.View(view, ver); ok && every < ros.MinPollEvery {
						every = v.Interval
					}
					if every < ros.MinPollEvery {
						every = time.Second
					}
					if preq, perr := res.Build(view, ver, ros.Poll(every), p); perr == nil {
						r.logf("rosrun: stream %s/%s ditolak router (%v); beralih ke poll %s", res.ID, view, ev.Err, every)
						if second, werr := r.Watch(ctx, preq); werr == nil {
							forward(second)
							return
						}
					}
				}
			}
			if ev.Kind != KindError {
				got = true
			}
			select {
			case out <- ev:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}
