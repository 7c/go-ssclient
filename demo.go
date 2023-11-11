package main

import (
	"log"
	"os"
	"time"

	"github.com/7c/go-ssclient/shared"

	"github.com/fatih/color"
)

func main() {
	ssc, err := shared.NewSSClient("ss://AEAD_CHACHA20_POLY1305:HwhYX94emfSVhMD@247.246.79.138:903", true, time.Second*15)
	if err != nil {
		log.Fatalln(err)
	}

	// launch socks5 as frontend listener and connect via tcp to shadowsocks client
	ssc.LaunchWithSocks5(":1080", false)

	// go func() {
	// 	time.Sleep(time.Second * 5)
	// 	ssc.Disconnect()
	// }()

	// manage the ssc channels
	for {
		select {
		case <-ssc.Ctx.Done():
			color.Green("Disconnected")
			os.Exit(1)
		case err3 := <-ssc.Channels.ChanError:
			color.Red("Errored", err3)
		case <-ssc.Channels.ChanTCPReady:
			color.Blue("TCP Listener is ready")
			_, err3 := ssc.TestSocks5(":1080", "tcp")
			if err3 != nil {
				color.Yellow("TCP Test failed:%s", err3)
				ssc.Disconnect()
			}
			color.Yellow("TCP Test succeeed")
		case <-ssc.Channels.ChanUDPReady:
			color.Yellow("UDP Listener is ready")
			// ssc.TestSocks5(":1080", "udp")
		}
	}

}
