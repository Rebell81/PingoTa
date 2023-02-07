package ping

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func makeBody() ([]byte, error) {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: make([]byte, 64),
		},
	}

	return msg.Marshal(nil)
}

func Ping(dst net.Addr) (int64, error) {
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		return 0, fmt.Errorf("error on ListenPacket: %w", err)
	}

	defer func(conn *icmp.PacketConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(conn)

	if err := conn.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
		return 0, fmt.Errorf("can't set connection deadline %w", err)
	}

	body, err := makeBody()
	if err != nil {
		return 0, fmt.Errorf("issue with make body for request: %w", err)
	}

	start := time.Now()

	if _, err := conn.WriteTo(body, dst); err != nil {
		return 0, fmt.Errorf("write to connection error: %w", err)
	}

	replyBytes := make([]byte, 1500)
	replySize, _, err := conn.ReadFrom(replyBytes)
	if err != nil {
		return 0, fmt.Errorf("read from connection err: %w", err)
	}

	duration := time.Since(start)

	reply, err := icmp.ParseMessage(1, replyBytes[:replySize])
	if err != nil {
		return 0, fmt.Errorf("issue with parse raw reply: %w", err)
	}

	switch reply.Code {
	case 0:
		return duration.Milliseconds(), nil
	case 11:
		return 0, TimeoutErr
	default:
		return 0, HostUnreachableErr
	}
}
