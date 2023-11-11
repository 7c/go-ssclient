# ssclient for go using go-shadowsocks2

This package is basic wrapper around a version of https://godoc.org/github.com/shadowsocks/go-shadowsocks2/. go-shadowsocks2 does not support acting as a package. This is why i had to take a stable version and build it in this package statically.. Also added kind of error handling for testing purphoses

Connection Error Handling, Timeout Handling has been added with context and go-channels.

## tcp/udp SSClient with socks5 listener
```
ssc, err := shared.NewSSClient("ss://AEAD_CHACHA20_POLY1305:HwhYX94emfSVhMD@217.244.79.108:903", true, time.Second*15)
if err != nil {
    log.Fatalln(err)
}
// launch tcp-only and listen on ':1080" socks5
ssc.LaunchWithSocks5(":1080", false)

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
```
