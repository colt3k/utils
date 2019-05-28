package tcp

import (
	"net"
	"strconv"

	log "github.com/colt3k/nglog/ng"
)

//Server start tcp server
func Server(host, port, contype string) {

	if len(host) > 0 && len(port) > 0 {
		if port, _ := strconv.ParseInt(port, 10, 0); port <= 1024 {
			log.Println("If using a port at or below 1024, it may require root or superuser access.")
		}
		// Listen for incoming connections.
		l, err := net.Listen(contype, host+":"+port)
		if err != nil {
			log.Logf(log.FATAL, "error listening\n%+v", err)
		}
		// Close the listener when the application closes.
		defer l.Close()
		log.Logln(log.INFO, "TCP Server Listening on "+host+":"+port)
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				log.Logf(log.FATAL, "issue accepting\n%+v", err)
			}
			// Handle connections in a new goroutine.
			go handleServerRequestTCP(conn)
		}
	}
}

/*
Handle any incoming TCP requests.
*/
func handleServerRequestTCP(conn net.Conn) {
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
