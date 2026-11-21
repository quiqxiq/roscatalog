// 04 — Workflow Mikhmon: profile + on-login, monitor profile scheduler, dan voucher.
package main

import (
	"context"
	"fmt"
	"time"

	"example.com/roscatalog/mikhmon"
	"example.com/roscatalog/rosrun"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sess, err := rosrun.Open(ctx, rosrun.Config{
		Address:  "192.168.88.1:8728",
		Username: "admin",
		Password: "secret",
		Timeout:  10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	client := mikhmon.New(sess)

	// 1. Buat script on-login sesuai format Mikhmon v4
	onLogin, err := mikhmon.OnLoginScript(mikhmon.OnLoginOptions{
		Mode:         mikhmon.ExpiryRemoveRecord, // "remc": hapus user saat expired + catat penjualan
		Price:        "5000",                     // harga modal
		Validity:     "1d",                       // masa aktif 1 hari
		SellingPrice: "7000",                     // harga jual
		Profile:      "paket-1hari",              // nama profile
		LockUser:     true,                       // lock MAC user saat pertama login
	})
	if err != nil {
		panic(err)
	}

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
}
