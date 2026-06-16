package metadata

import (
	"net/netip"
	"os"
	"testing"
)

func TestLookup_WellKnownIPs(t *testing.T) {
	// Set up paths to your local downloaded test files
	cityDB := "../../data/dbip-city.mmdb"
	asnDB := "../../data/asn.mmdb"

	// Skip if files don't exist yet locally to prevent breaking CI/CD pipelines
	if _, err := os.Stat(cityDB); os.IsNotExist(err) {
		t.Skip("Skipping test; local MMDB files not found in ../../data/")
	}

	err := Init(cityDB, asnDB)
	if err != nil {
		t.Fatalf("Failed to initialize metadata DBs: %v", err)
	}
	defer Close()

	tests := []struct {
		name        string
		ipStr       string
		expectedASN uint
		expectedISP string
		expectGeo   bool
	}{
		{
			name:        "Cloudflare DNS",
			ipStr:       "1.1.1.1",
			expectedASN: 13335,
			expectedISP: "Cloudflare",
			expectGeo:   true,
		},
		{
			name:        "Google DNS",
			ipStr:       "8.8.8.8",
			expectedASN: 15169,
			expectedISP: "Google",
			expectGeo:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := netip.MustParseAddr(tt.ipStr)
			info := Lookup(addr)

			if info.IP != tt.ipStr {
				t.Errorf("Expected IP %s, got %s", tt.ipStr, info.IP)
			}

			// Validate ASN if the database is loaded
			if asnReader != nil && info.ASN != tt.expectedASN {
				t.Errorf("Expected ASN %d, got %d (ISP: %s)", tt.expectedASN, info.ASN, info.ISP)
			}

			// Soft check on GeoData (Country codes are typically 2 letters)
			if cityReader != nil && tt.expectGeo && len(info.Country) != 2 {
				t.Errorf("Expected valid ISO country code, got '%s'", info.Country)
			}
		})
	}
}
