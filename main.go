package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

type handler struct {
}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)

	msg.Authoritative = true
	msg.RecursionAvailable = false
	msg.Rcode = dns.RcodeNameError

	switch r.Question[0].Qtype {
	case dns.TypeAAAA:
		switch r.Question[0].Name {
		case "time.botcheckup.io":
			now := time.Now().Unix()
			binary.LittleEndian.PutUint64(a, now)

			ipv6 := net.IP{0x26, 0x06, 0xbd, 0xc0, 0xff, 0xff, 0x00, 0x00, a[7], a[6], a[5], a[4], a[3], a[2], a[1], a[0]}
			msg.Rcode = dns.RcodeSuccess
			msg.Answer = append(msg.Answer, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: domain, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 1},
				AAAA: ipv6,
			})
		}
	}
	w.WriteMsg(&msg)
}

func main() {
	address := flag.String("address", "127.0.0.1:9999", "address to listen on")
	flag.Parse()

	srv := &dns.Server{
		Addr: *address,
		Net:  "udp",
	}
	srv.Handler = &handler{}

	fmt.Println("listening on", *address)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
