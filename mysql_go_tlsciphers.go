package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	ciphers := []uint16{0x0005, 0x000a, 0x002f, 0x0035, 0x003c, 0x009c, 0x009d,
		0xc007, 0xc009, 0xc00a, 0xc011, 0xc012, 0xc013, 0xc014, 0xc023,
		0xc027, 0xc02f, 0xc02b, 0xc030, 0xc02c, 0xcca8, 0xcca9}
	rounds := 5

	for n, cipher := range ciphers {
		mysql.RegisterTLSConfig("custom"+string(n), &tls.Config{
			InsecureSkipVerify:       true,
			PreferServerCipherSuites: true,
			CipherSuites:             []uint16{cipher},
		})
		dsn := "msandbox:msandbox@tcp(127.0.0.1:5717)/test?tls=custom" + string(n)
		var ciphername string
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("Connection with cipher=%x failed: %s\n", cipher, err)
			continue
		}
		rows, err := db.Query("SHOW SESSION STATUS LIKE 'Ssl_cipher'")
		if err != nil {
			fmt.Printf("Connection with cipher=%x failed: %s\n", cipher, err)
			continue
		}
		for rows.Next() {
			var settingname string
			rows.Scan(&settingname, &ciphername)
		}
		rows.Close()
		db.Close()
		fmt.Printf("Testing with cipher=%x (MySQL cipher name: %s)\n", cipher, ciphername)
	CIPHERLOOP:
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 2)
			start := time.Now()
			for j := 0; j < rounds; j++ {
				db, err := sql.Open("mysql", dsn)
				if err != nil {
					println(err)
				}
				_, err = db.Exec("DO 1")
				if err != nil {
					fmt.Printf("%s\n", err)
					break CIPHERLOOP
				}
				db.Close()
			}
			elapsed := time.Since(start)
			fmt.Printf("%d rounds in %s (%fms per loop)\n", rounds, elapsed, (float64(elapsed.Nanoseconds()/int64(rounds)))/1000000)
		}

		// Fetch a large resultset
		db, err = sql.Open("mysql", dsn)
		_, _ = db.Exec("DO 1")
		if err != nil {
			fmt.Printf("Connection with cipher=%x failed: %s\n", cipher, err)
			continue
		}
		for i := 0; i < 3; i++ {
			bigstart := time.Now()
			rows, err = db.Query("SELECT REPEAT('x', @@global.max_allowed_packet)")
			if err != nil {
				fmt.Printf("Connection with cipher=%x failed: %s\n", cipher, err)
				continue
			}
			for rows.Next() {
				var bigpacket string
				rows.Scan(&bigpacket)
			}
			bigelapsed := time.Since(bigstart)
			fmt.Printf("Fetching a big result: %s\n", bigelapsed)
			rows.Close()
		}
		db.Close()
	}
}
