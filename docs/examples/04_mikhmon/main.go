// 04 — Workflow Mikhmon: profile + on-login, monitor profile scheduler, dan voucher.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/quiqxiq/roscatalog/mikhmon"
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
	dryRun := flag.Bool("dry-run", true, "generate script on-login tanpa menulis ke router")
	flag.Parse()

	onLogin, err := mikhmon.OnLoginScript(mikhmon.OnLoginOptions{
		Mode:         mikhmon.ExpiryRemoveRecord,
		Price:        "5000",
		Validity:     "1d",
		SellingPrice: "7000",
		Profile:      "paket-1hari",
		LockUser:     true,
	})
	if err != nil {
		panic(err)
	}
	if *dryRun {
		fmt.Println("[dry-run] on-login script length:", len(onLogin))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sess, err := rosrun.Open(ctx, rosrun.Config{
		Address:  *addr,
		Username: *user,
		Password: *pass,
		Timeout:  10 * time.Second,
		Status:   rosrun.StatusConfig{Disable: true},
	})
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	client := mikhmon.New(sess)

	// 2. Tambah Profile Hotspot
	profileID, err := client.AddProfile(ctx, mikhmon.Profile{
		Name:              "paket-1hari",
		RateLimit:         "2M/2M",
		SharedUsers:       "1",
		StatusAutorefresh: "1m",
		OnLogin:           onLogin,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Profile dibuat: paket-1hari (id=%s)\n", profileID)

	// 3. Pasang scheduler Monitor Profile (agar lampu indikator di web Mikhmon v4 hijau)
	if err := client.InstallProfileMonitor(ctx, "paket-1hari", mikhmon.ExpiryRemoveRecord); err != nil {
		panic(err)
	}
	fmt.Println("Monitor Profile scheduler berhasil dipasang")

	// 4. Tambah voucher hotspot
	name := "vc101"
	comment := mikhmon.MikhmonComment(name, name, "batch-1") // format "vc-batch-1"
	userID, err := client.AddUser(ctx, mikhmon.User{
		Name:        name,
		Password:    name,
		Profile:     "paket-1hari",
		Server:      "all",
		Comment:     comment,
		LimitUptime: "1d",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Voucher dibuat: %s (id=%s, comment=%s)\n", name, userID, comment)

	// 5. Cek jumlah user dalam profile
	total, err := client.CountUsers(ctx, "paket-1hari")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total user berprofile paket-1hari: %d\n", total)

	// 6. Baca voucher secara typed menggunakan fields.MikhmonVoucherRow
	vouchers, err := client.VouchersTyped(ctx, comment)
	if err == nil {
		for _, v := range vouchers {
			fmt.Printf("- Voucher: %s, Profile: %s, LimitUptime: %s\n", v.Name, v.Profile, v.LimitUptime)
		}
	}
}
