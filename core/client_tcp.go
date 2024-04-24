package core

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/hati-sh/hati/common"
)

type ClientTcp struct {
	tlsCertificate tls.Certificate
	host           string
	port           string
	tlsEnabled     bool
}

func NewClientTcp(host string, port string, tlsEnabled bool) (ClientTcp, error) {
	cert, err := common.GenX509KeyPair()
	if err != nil {
		return ClientTcp{}, err
	}

	return ClientTcp{
		tlsCertificate: cert,
		host:           host,
		port:           port,
		tlsEnabled:     tlsEnabled,
	}, nil
}

func (s ClientTcp) Connect() error {
	// config := tls.Config{InsecureSkipVerify: true, Certificates: []tls.Certificate{s.tlsCertificate}}
	config := tls.Config{}
	config.Rand = rand.Reader

	conn, err := net.Dial("tcp", s.host+":"+s.port)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	// state := conn.ConnectionState()
	// for _, v := range state.PeerCertificates {
	// 	fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
	// 	fmt.Println(v.Subject)
	// }
	// log.Println("client: handshake: ", state.HandshakeComplete)
	// log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	// msg := NewMessage()
	// msg.SetPayload([]byte("dziala!!!"))
	// msg.SetExtraSpace([4]byte{'D', 'U', 'P', 'A'})

	// msgBytes := msg.Bytes()

	// fmt.Println(msgBytes)
	// fmt.Println(string(msgBytes))
	// msgBytes

	writer := bufio.NewWriter(conn)

	rc := 0

	for i := 0; i < 1000; i++ {
		_, err := writer.Write([]byte("dziala!!!\n"))
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		writer.Flush()

		// log.Printf("client: wrote %q (%d bytes)", string(msgBytes), n)

		reply := make([]byte, 256)
		n, err := conn.Read(reply)
		if err != nil {
			log.Fatal(err)
		}

		if n > 0 {
			rc++
		}
		// log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
		// log.Print("client: exiting")
	}

	fmt.Println(rc)

	return nil
}

// func handleClient(conn net.Conn) {
// 	defer conn.Close()
// 	buf := make([]byte, 512)
// 	for {
// 		log.Print("server: conn: waiting")
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			if err != nil {
// 				log.Printf("server: conn: read: %s", err)
// 			}
// 			break
// 		}
// 		log.Printf("server: conn: echo %q\n", string(buf[:n]))
// 		n, err = conn.Write(buf[:n])

// 		n, err = conn.Write(buf[:n])
// 		log.Printf("server: conn: wrote %d bytes", n)

// 		if err != nil {
// 			log.Printf("server: write: %s", err)
// 			break
// 		}
// 	}
// 	log.Println("server: conn: closed")
// }
