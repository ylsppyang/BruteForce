package main

import (
	bruteforce "BruteForce"
	"fmt"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Loginfailed(proto uint, ip *net.IP, time time.Time) {
	defer wg.Done()
	result, err := bruteforce.BruteForceCheck(proto, ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	if result {
		/* We just print a log here. You can add this ip in black ip list to not handle request from it */
		fmt.Printf("Yes! It's a brute force, protoid:%v\n", proto)
	}
}

func initConfig() {
	bruteforce.InitBruteForce()
	for i := bruteforce.BF_PROTO_HTTP; i < bruteforce.BF_PROTO_END; i++ {
		bruteforce.Config_bf_setting(uint(i), uint(i*3))
	}
	bruteforce.Config_bf_setting(0, 0)
}

func testAger() {
	fmt.Println("===================== start testAger ========================")
	bruteforce.Config_bf_setting(0, 3)
	wg.Add(3)
	ip := net.ParseIP("178.21.2.13")
	Loginfailed(bruteforce.BF_PROTO_HTTP, &ip, time.Now())
	Loginfailed(bruteforce.BF_PROTO_HTTP, &ip, time.Now())
	time.Sleep((bruteforce.BRUTE_TIME + 2) * time.Second)
	/* There is no log, because item is out of time,
	   but if you sleep only BRUTE_TIME-1 seconds, there should be a brute force log */
	Loginfailed(bruteforce.BF_PROTO_HTTP, &ip, time.Now())
	fmt.Println("===================== end of testAger ========================")
}

func main() {
	initConfig()
	fmt.Println(bruteforce.Bf_setting)
	ip := net.ParseIP("178.23.2.12")
	for i := 0; i < 10; i++ {
		for j := bruteforce.BF_PROTO_HTTP; j < bruteforce.BF_PROTO_END; j++ {
			wg.Add(1)
			go Loginfailed(uint(j), &ip, time.Now())
		}
		time.Sleep(1 * time.Second)
	}
	testAger()
	wg.Wait()
}
