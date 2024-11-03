package ip2country_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/theshoutingparrot/ip2country"
)

func BenchmarkGetCountry(b *testing.B) {

	ips, err := ip2country.Load("./dbip-country.csv")
	if err != nil {
		b.Fatal(err)
	}

	for i := 1; i <= b.N; i++ {
		ips.GetCountry(createRandomIP())
	}

}

func createRandomIP() string {
	p1 := rand.Int31n(256)
	p2 := rand.Int31n(256)
	p3 := rand.Int31n(256)
	p4 := rand.Int31n(256)

	return fmt.Sprintf("%d.%d.%d.%d", p1, p2, p3, p4)
}

func TestIPtoCountryLookup(t *testing.T) {
	ips, err := ip2country.Load("./dbip-country.csv")
	if err != nil {
		t.Fatal(err)
	}

	country := ips.GetCountry("0.0.0.0")
	if country != "ZZ" {
		t.Errorf("Expected ZZ but found %s", country)
	}
	country = ips.GetCountry("0.0.0.1")
	if country != "ZZ" {
		t.Errorf("Expected ZZ but found %s", country)
	}
	country = ips.GetCountry("0.255.255.255")
	if country != "ZZ" {
		t.Errorf("Expected ZZ but found %s", country)
	}

	country = ips.GetCountry("50.97.196.208")
	if country != "US" {
		t.Errorf("Expected US but found %s", country)
	}

	country = ips.GetCountry("50.97.198.135")
	if country != "CN" {
		t.Errorf("Expected US but found %s", country)
	}

	country = ips.GetCountry("50.97.198.135")
	if country != "CN" {
		t.Errorf("Expected US but found %s", country)
	}

	country = ips.GetCountry("200.95.185.34")
	if country != "CL" {
		t.Errorf("Expected CL but found %s", country)
	}

	country = ips.GetCountry("223.255.255.255")
	if country != "AU" {
		t.Errorf("Expected AU but found %s", country)
	}
}

func TestGetCountryMulti(t *testing.T) {
	ips, err := ip2country.Load("./dbip-country.csv")
	if err != nil {
		t.Fatal(err)
	}

	countries := ips.GetCountryMulti("35.185.131.112", "35.185.133.191", "64.215.100.142", "94.46.48.46", "159.8.131.119")
	if countries[0] != "US" {
		t.Errorf("Expected US but found %s", countries[0])
	}
	if countries[1] != "TW" {
		t.Errorf("Expected TW but found %s", countries[1])
	}
	if countries[2] != "BR" {
		t.Errorf("Expected BR but found %s", countries[2])
	}
	if countries[3] != "GB" {
		t.Errorf("Expected GB but found %s", countries[3])
	}
	if countries[4] != "NL" {
		t.Errorf("Expected NL but found %s", countries[4])
	}

}

func TestIP4AddressToInt(t *testing.T) {
	ipNumb, err := ip2country.Ip4ToInt("0.0.0.255")
	if err != nil {
		t.Error(err)
	}
	if ipNumb != 255 {
		t.Errorf("Expected 255 but found %d", ipNumb)
	}
	//-----------------------------------------------
	ipNumb, err = ip2country.Ip4ToInt("0.0.1.255")
	if err != nil {
		t.Error(err)
	}
	if ipNumb != 511 {
		t.Errorf("Expected 511 but found %d", ipNumb)
	}
	//-----------------------------------------------
	ipNumb, err = ip2country.Ip4ToInt("255.0.0.0")
	if err != nil {
		t.Error(err)
	}
	if ipNumb != 4278190080 {
		t.Errorf("Expected 4278190080 but found %d", ipNumb)
	}
	//-----------------------------------------------
	ipNumb, err = ip2country.Ip4ToInt("255.255.255.255")
	if err != nil {
		t.Error(err)
	}
	if ipNumb != 4294967295 {
		t.Errorf("Expected 4294967295 but found %d", ipNumb)
	}
}

func TestIP6AddressToInt(t *testing.T) {
	b := big.NewInt(0)
	ip2country.Ip6ToInt("fec0::", b)
	t.Logf("result: %v", b)

	b = big.NewInt(0)
	ip2country.Ip6ToInt("::0000", b)
	t.Logf("result: %v", b)

	b = big.NewInt(0)
	ip2country.Ip6ToInt("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", b)
	t.Logf("result: %v", b)

	b = big.NewInt(0)
	ip2country.Ip6ToInt("0000:0000:0000:0000:0000:0000:0000:00ff", b)
	t.Logf("result: %v", b)
}

func TestIP4and6(t *testing.T) {
	ips, err := ip2country.Load("./dbip-country-test.csv")
	if err != nil {
		t.Fatalf("got err: %v", err)
	}

	v4orv6 := ips.GetCountry("1.1.1.1")
	if v4orv6 != "V4" {
		t.Fatalf("expected V4, got %v", v4orv6)
	}

	v4orv6 = ips.GetCountry("0fff:1fff:2fff:3fff:4fff:5fff:ffff:ffff")
	if v4orv6 != "V6" {
		t.Fatalf("expected V4, got %v", v4orv6)
	}
}
