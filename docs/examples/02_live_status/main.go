// 02 — Realtime: WatchStatus terlihat realtime (connected/warning/error).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/quiqxiq/roscatalog/rosrun"
)

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	addr := flag.String("addr", envOr("ROS_ADDR", "192.168.88.1:8728"), "alamat API RouterOS")
	user := flag.String("user", envOr("ROS_USER", "admin"), "username API")
	pass := flag.String("pass", envOr("ROS_PASS", ""), "password API")
	target := flag.String("ping-target", "8.8.8.8", "target ping status (kosong = API health saja)")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sess, err := rosrun.Open(ctx, rosrun.Config{
		Address:  *addr,
		Username: *user,
		Password: *pass,
		Timeout:  10 * time.Second,
		Status:   rosrun.StatusConfig{Target: *target},
	})
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	ch := sess.WatchStatus(ctx)
	for i := 0; i < 5; i++ {
		select {
		case st := <-ch:
			fmt.Printf("state=%s connected=%v api=%v warnings=%v ping(%d/%d) ms=%.0f\n",
				st.State, st.Connected, st.APIConnected, st.Warnings,
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
