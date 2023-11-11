package shared

import (
	"os"
	"sync"

	nested "github.com/antonfisher/nested-logrus-formatter"
	logrus "github.com/sirupsen/logrus"
)

type SSClientManager struct {
	mu      sync.Mutex
	clients []*SSClient
	log     *logrus.Logger

	channelNewClient chan *SSClient
}

func NewSSClientManager() *SSClientManager {

	man := SSClientManager{
		channelNewClient: make(chan *SSClient),
	}

	man.log = logrus.New()
	// man.log = logrus.WithField("component", "SSClient")
	man.log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		CallerFirst: true,

		FieldsOrder: []string{"component", "category"},
	})
	// man.log.WithField("component", "SSClient").Info("test")

	go man.eventloop()
	return &man
}

func (manager *SSClientManager) AddClient(client *SSClient) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.clients = append(manager.clients, client)
	manager.channelNewClient <- client
	manager.log.Info("New client has been added")
}

func (manager *SSClientManager) eventloop() {
	manager.log.Info("SSClientManager event loop initiated")
	for {
		select {
		case client := <-manager.channelNewClient:
			manager.log.Infof("managing new ssclient %s", client.SSUrl.Host)
			// Handle the new client, e.g., start listening to its ChanError
			go func(c *SSClient) {
				for err := range c.Channels.ChanError {
					// Handle the error, log it, restart the client, etc.
					manager.log.Error(err)
				}
			}(client)
			// Add other select cases as needed...
		}
	}
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
}
