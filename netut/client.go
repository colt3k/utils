package netut

//Host stores parts of a Server for network communication
type Host struct {
	IP       string
	Port     string
	URL      string
	Pass     bool
	Protocol string
}
type Client interface {
	Client(contype string, hostAr []string, proxy string) *[]Host
}
