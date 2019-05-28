package udp

import (
	"net"

	log "github.com/colt3k/nglog/ng"
)

//Server create UDP Server
func Server(host, port, contype string) {

	if len(host) > 0 && len(port) > 0 {
		udpAddr, err := net.ResolveUDPAddr("udp4", host+":"+port)

		if err != nil {
			log.Logf(log.FATAL, "issue resolving UDP address\n%+v", err)
		}

		// setup listener for incoming UDP connection
		ln, err := net.ListenUDP(contype, udpAddr)
		if err != nil {
			log.Logf(log.FATAL, "error listening:\n%+v", err.Error())
		}

		log.Logln(log.INFO, "UDP Server Listening on "+host+":"+port)

		defer ln.Close()

		for {
			// Handle connections in a new goroutine.
			go handleServerRequestUDP(ln)
		}

	}

}

/*
Handle any incoming requests.
*/
func handleServerRequestUDP(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.Logf(log.ERROR, "issue reading\n%+v", err)
	}
	//log.Println("BufSize:",len(strings.TrimSpace(string(buf[:reqLen]))))
	log.Println("Data Received:", string(buf[:reqLen]))
	// Send a response back to person contacting us.
	_, err = conn.Write([]byte("Message received.\n"))
	if err != nil {
		log.Logf(log.ERROR, "issue writing %+v", err)
	}
	// Close the connection when you're done with it.
	err = conn.Close()
	if err != nil {
		log.Logf(log.ERROR, "issue closing %+v", err)
	}
}
