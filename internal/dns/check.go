package dns

import (
	"errors"
	"net"
)

var (
	ErrDomainIsBusy = errors.New("domain is busy")
)

func CheckDomain(domain string) error {
	addrs, err := net.LookupHost(domain)
	if err != nil {
		var dnsError *net.DNSError
		if errors.As(err, &dnsError) {
			return nil
		}

		return err
	}

	if len(addrs) > 0 {
		return ErrDomainIsBusy
	}

	return nil
}
