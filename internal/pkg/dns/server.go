package dns

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/mirzakhany/pure_dns/internal/pkg/config"
)

type handler struct{}

var domains map[string]string

func mapToIP(domain string) (string, bool) {

	ipAddr, ok := domains[strings.TrimSuffix(domain, ".")]
	if !ok {
		log.Println("not ok")
		return "", ok
	}
	return ipAddr, ok
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := mapToIP(domain)
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		}
	}
	w.WriteMsg(&msg)
}

// StartDNSServer start dns server
func StartDNSServer(conf config.ConfYaml) error {
	domains = conf.DomainMap
	srv := &dns.Server{Addr: conf.Server.Address + ":" + strconv.Itoa(conf.Server.Port), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
	return nil
}
