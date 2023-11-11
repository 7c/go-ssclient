package shared

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/7c/go-ssclient/core"
	"github.com/7c/go-ssclient/global"
	"github.com/7c/go-ssclient/socks"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/proxy"
)

type SSClient struct {
	SSUrl    *url.URL
	Cipher   core.Cipher
	Password string

	Channels *global.SSClientChannels

	Ctx    context.Context
	cancel context.CancelFunc
}

func (ssc *SSClient) Disconnect() {
	log.Printf("SSClient disconnect signal received\n")
	ssc.cancel()
	// <-ssc.Ctx.Done()
	// log.Printf("SSClient disconnected")
}

func (ssc *SSClient) TestSocks5(socksListenAddr, proto string) (bool, error) {
	logf("TestSocks5 '%s','%s'", proto, socksListenAddr)
	directDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	// var dialContext proxy.DialContext
	dialer, err := proxy.SOCKS5(proto, socksListenAddr, nil, directDialer)
	if err != nil {
		return false, fmt.Errorf("could not parse proxy address")
	}
	transport := &http.Transport{
		DialContext: dialer.(proxy.ContextDialer).DialContext,
	}
	// xhr request to verify the address
	client := resty.New()
	client.SetTransport(transport)
	resp, err := client.R().Get("https://ip4.ip8.com")
	if err != nil {
		return false, fmt.Errorf("xhr error:%s", err)
	}

	logf("TestSocks5 response body: '%s'", color.MagentaString("%s", resp.Body()))
	if strings.TrimSpace(string(resp.Body())) == ssc.SSUrl.Hostname() {
		return true, nil
	}
	return false, fmt.Errorf("unexpected response: %s", resp)
}

func (ssc *SSClient) LaunchWithSocks5(socksListenAddr string, udpEnabled bool) *SSClient {
	serverAddr := ssc.SSUrl.Host

	go SocksLocal(ssc.Ctx, socksListenAddr, serverAddr, ssc.Cipher.StreamConn, ssc.Channels)
	if udpEnabled {
		socks.UDPEnabled = true
		go UdpSocksLocal(ssc.Ctx, socksListenAddr, serverAddr, ssc.Cipher.PacketConn, ssc.Channels)
	}

	logf("LaunchWithSocks5 done")
	return ssc
}

// valid ssURL format
// ss://[chiper]:your-password@[server_address]:[port]
func NewSSClient(ssURL string, verbose bool, timeout time.Duration) (*SSClient, error) {
	LoggerConfig.Verbose = verbose
	parsed, err := url.Parse(ssURL)
	if err != nil || parsed.Scheme != "ss" {
		return nil, fmt.Errorf("please have a valid ssURL DSN:%s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	ssc := &SSClient{
		SSUrl:  parsed,
		Ctx:    ctx,
		cancel: cancel,
		Channels: &global.SSClientChannels{
			ChanError:    make(chan error),
			ChanTCPReady: make(chan bool),
			ChanUDPReady: make(chan bool),
			Timeout:      timeout,
		},
	}

	cipher := parsed.User.Username()
	password, b1 := parsed.User.Password()
	if !b1 {
		return nil, fmt.Errorf("please provide valid password")
	}
	ssc.Password = password

	// key := make([]byte, 32)
	var key []byte
	// io.ReadFull(rand.Reader, key)
	chpr, err := core.PickCipher(cipher, key, password)
	if err != nil {
		return nil, fmt.Errorf("error by cipher selection: %s", err)
	}
	ssc.Cipher = chpr
	return ssc, nil
}
