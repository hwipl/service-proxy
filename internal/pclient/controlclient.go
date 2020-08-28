package pclient

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hwipl/service-proxy/internal/network"
)

// controlClient stores control client information
type controlClient struct {
	serverAddr *net.TCPAddr
	tlsConfig  *tls.Config
	specs      []*ServiceSpec
	conn       net.Conn
}

// runClient runs the control client
func (c *controlClient) runClient() {
	// connect to server
	conn, err := net.DialTCP("tcp", nil, c.serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
	if c.tlsConfig != nil {
		c.conn = tls.Client(conn, c.tlsConfig)
	}
	defer c.conn.Close()
	log.Println("Connected to server", c.serverAddr)

	// send service specs to server
	active := 0
	for _, spec := range c.specs {
		log.Printf("Sending service registration %s to server", spec)
		m := spec.ToMessage()
		network.WriteToConn(c.conn, m.Serialize())

		// read reply messages from server
		var msg network.Message
		buf := network.ReadFromConn(c.conn)
		if buf == nil {
			log.Println("Closing connection to server")
			return
		}
		msg.Parse(buf)

		// handle message types
		var spec ServiceSpec
		replyFmt := "Server reply: service registration %s %s\n"
		switch msg.Op {
		case network.MessageOK:
			spec.FromMessage(&msg)
			log.Printf(replyFmt, &spec, "OK")
			active++
		case network.MessageErr:
			spec.FromMessage(&msg)
			log.Printf(replyFmt, &spec, "ERROR")
		default:
			// unknown message, stop here
			log.Println("Unknown reply from server, " +
				"closing connection")
			return
		}
	}

	// are any services active on the server?
	if active == 0 {
		log.Println("Could not register any service on the server, " +
			"closing connection")
		return
	}
	log.Printf("Registered %d service(s) on the server, "+
		"keeping connection open", active)

	// keep connection open
	go func() {
		for {
			// send a keep-alive/NOP message every 15 seconds
			time.Sleep(15 * time.Second)
			keepAlive := network.Message{Op: network.MessageNop}
			if !network.WriteToConn(c.conn,
				keepAlive.Serialize()) {
				return
			}
		}
	}()
	for {
		if network.ReadFromConn(c.conn) == nil {
			log.Println("Closing connection to server")
			return
		}
	}
}

// RunControlClient runs the control client
func RunControlClient(cntrlAddr *net.TCPAddr, tlsConfig *tls.Config,
	// parse service specifications (format "<protocol>:<port>:<destPort>")
	serviceSpecs string) {
	services := strings.Split(serviceSpecs, ",")
	var specs []*ServiceSpec
	for _, s := range services {
		specs = append(specs, ParseServiceSpec(s))
	}

	// print info and run control client
	ip := ""
	if cntrlAddr.IP != nil {
		ip = fmt.Sprintf("%s", cntrlAddr.IP)
	}
	tlsInfo := ""
	if tlsConfig != nil {
		tlsInfo = "in mTLS mode "
	}
	log.Printf("Starting client %sand connecting to server %s:%d\n",
		tlsInfo, ip, cntrlAddr.Port)

	// create and run control client
	c := controlClient{
		serverAddr: cntrlAddr,
		tlsConfig:  tlsConfig,
		specs:      specs,
	}
	c.runClient()
}
