// 02 — Realtime: WatchStatus terlihat realtime (connected/warning/error).
package main

import (
	"context"
	"fmt"
	"time"

	"example.com/roscatalog/rosrun"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sess, err := rosrun.Open(ctx, rosrun.Config{
		Address:  "192.168.88.1:8728",
		Username: "admin",
		Password: "r00t",
		Timeout:  10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	ch := sess.WatchStatus(ctx)
	for i := 0; i < 5; i++ {
		select {
		case st := <-ch:
			fmt.Printf("state=%s connected=%v warnings=%v ping(%d/%d) ms=%.0f\n",
				st.State, st.Connected, st.Warnings,
				st.Ping.Sent, st.Ping.Received, float64(st.Ping.AvgRTT)/float64(time.Millisecond))
			if st.Connected && st.State == rosrun.StatusConnected {
				fmt.Println("terhubung; tekan Ctrl+C untuk berhenti")
			}
		case <-time.After(5 * time.Second):
			fmt.Println("timeout menunggu status")
			return
		}
	}
}
