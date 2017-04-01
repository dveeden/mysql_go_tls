package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	tlsvalues := []string{"false", "skip-verify", "cache"}
	rounds := 100

	scache := tls.NewLRUClientSessionCache(32)
	mysql.RegisterTLSConfig("cache", &tls.Config{
		InsecureSkipVerify: true,
		ClientSessionCache: scache,
	})

	for _, tlsvalue := range tlsvalues {
		dsn := "msandbox:msandbox@tcp(localhost:5717)/test?tls=" + tlsvalue
		fmt.Printf("Testing with tls=%s\n", tlsvalue)
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
					break
				}
				db.Close()
			}
			elapsed := time.Since(start)
			fmt.Printf("%d rounds in %s (%fms per loop)\n", rounds, elapsed,
				(float64(elapsed.Nanoseconds()/int64(rounds)))/1000000)
		}
	}
}
