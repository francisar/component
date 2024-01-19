package main

import (
	"context"
	"fmt"
	"github.com/francisar/component/network/pkg/tools/ping"
	"net"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

func main()  {
	pingCheck()
}


var a  = [...]string{
	"10.11.4.2",
}

func pingCheck()  {
	defer func() {
		if err := recover(); err != nil {
			// 在这里处理panic
			debug.PrintStack()
			panic(err)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(len(a))
	localAddr := net.IPAddr{
		IP:  net.ParseIP("10.31.0.18"),
		Zone: "",
	}
	for i, ip := range a {
		tempIp := ip
		index := i
		ctx := context.Background()

		go func() {
			defer func(ctx context.Context) {
				if err := recover(); err != nil {
					// 在这里处理panic
					println(err)
					debug.PrintStack()
				}
				wg.Done()
			}(ctx)
			ctxTimeout,_ := context.WithTimeout(ctx, 5 * time.Second)
			raddr := net.IPAddr{
				IP:  net.ParseIP(tempIp),
				Zone: "",
			}
			conn, err := ping.NewPingV4(&localAddr, &raddr, 3* time.Second, os.Getpid()+1 & 0xFFFF)
			defer func() {
				err := conn.Close()
				fmt.Println("close", tempIp, err)
			}()
			fmt.Println("start ping", tempIp, index)
			if err != nil {
				fmt.Println(tempIp, "new udpv4 failed", err)
				return
			}
			msg := fmt.Sprintf("hello %s", tempIp)
			sendErr := conn.SendIcmp([]byte(msg))
			if sendErr != nil {
				fmt.Println(tempIp, "send hello", sendErr)
				return
			}
			resp, recvErr := conn.RecvIcmp()
			for ;; {
				select {
				case <- ctxTimeout.Done():
					fmt.Println(string(resp.Body), tempIp, index, "timeout")
					return
				default:
					if recvErr != nil {
						fmt.Println(tempIp, "recv hello", recvErr)
						return
					}
					if resp.Dst.String() == tempIp {
						fmt.Println(string(resp.Body), resp.Src, resp.Dst, tempIp, index, resp.Timeout)
						return
					}
				}
			}
		}()
	}
    wg.Wait()
	fmt.Println("finished")


/*
   addr2 := net.IPAddr{
   		IP:   net.ParseIP("10.113.47.2"),
   		Zone: "",
   	}
	p2, err2 := ping.NewUdpV4(&addr2, 3* time.Second, os.Getpid())
	if err2 != nil {
		panic(err2)
	}

	sendErr2 := p2.SendIcmp([]byte("bb"))
	if sendErr2 != nil {
		panic(sendErr2)
	}

	recv2,recvErr2 := p2.RecvIcmp()

	if recvErr2 != nil {
		panic(recvErr2)
	}

	fmt.Println(recv2)


 */

}