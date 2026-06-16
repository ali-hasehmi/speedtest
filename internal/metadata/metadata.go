package metadata

import (
	"errors"
	"net/netip"

	"github.com/ali-hasehmi/speedtest/logger"
	"github.com/oschwald/maxminddb-golang/v2"
)

var (
	cityReader *maxminddb.Reader
	asnReader  *maxminddb.Reader
)

// Info represents the unified metadata payload
type Info struct {
	IP      string `json:"ip"`
	Country string `json:"country,omitempty"`
	ASN     uint   `json:"asn,omitempty"`
	ISP     string `json:"isp,omitempty"`
}

type CountryRecord struct {
	CountryCode string `maxminddb:"country_code"`
}

// ASNRecord matches the schema for DB-IP ASN / GeoLite2 ASN
type ASNRecord struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
}

// Init memory-maps the v2 databases. Empty paths are safely ignored.
func Init(cityDBPath, asnDBPath string) error {
	var err error

	if cityDBPath != "" {
		cityReader, err = maxminddb.Open(cityDBPath)
		if err != nil {
			return errors.New("failed to open city database: " + err.Error())
		}
		logger.Infof("City v2 database loaded from %s", cityDBPath)
	}

	if asnDBPath != "" {
		asnReader, err = maxminddb.Open(asnDBPath)
		if err != nil {
			return errors.New("failed to open ASN database: " + err.Error())
		}
		logger.Infof("ASN v2 database loaded from %s", asnDBPath)
	}

	return nil
}

// Lookup resolves a v2 netip.Addr into our Info struct.
func Lookup(ip netip.Addr) Info {
	info := Info{IP: ip.String()}

	// 1. Lookup City & Country Data
	if cityReader != nil {
		var countryRecord CountryRecord
		err := cityReader.Lookup(ip).Decode(&countryRecord)
		if err != nil {
			logger.Errorf("City lookup failed for IP %s: %v", info.IP, err)
		} else {
			info.Country = countryRecord.CountryCode
		}
	}

	// 2. Lookup ASN/ISP Data
	if asnReader != nil {
		var asnRecord ASNRecord
		err := asnReader.Lookup(ip).Decode(&asnRecord)
		if err != nil {
			logger.Errorf("ASN lookup failed for IP %s: %v", info.IP, err)
		} else {
			info.ASN = asnRecord.AutonomousSystemNumber
			info.ISP = asnRecord.AutonomousSystemOrganization
		}
	}

	return info
}

// Close unmaps the databases from memory.
func Close() {
	if cityReader != nil {
		cityReader.Close()
	}
	if asnReader != nil {
		asnReader.Close()
	}
}
