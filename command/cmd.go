package command

import (
	"fmt"
	"net"
	"time"

	"pt/ping"
)

type store interface {
	Write(host string, time int64) error
}

type cmd struct {
	store store
}

func New(store store) *cmd {
	return &cmd{store: store}
}

func (c *cmd) Run(hosts []string) error {
	for _, host := range hosts {
		var total int64

		ip, err := net.ResolveIPAddr("ip4", host)
		if err != nil {
			return err
		}

		for i := 0; i < 10; i++ {
			t, err := ping.Ping(&net.UDPAddr{IP: net.ParseIP(ip.String())})
			if err != nil {
				fmt.Println(err)
				break
			}
			total += t
		}

		fmt.Printf("%s - ping %s for %dms\n", time.Now().Format(time.RFC3339), host, total/10)

		if err := c.store.Write(host, total/10); err != nil {
			return err
		}
	}

	return nil
}
