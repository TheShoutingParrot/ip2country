// DB-IP.com
// You are free to use this database in your application, provided you give attribution to DB-IP.com for the data.
// In the case of a web application, you must include a link back to DB-IP.com on the pages displaying or using results from the database.

package ip2country

import (
	"bufio"
	"errors"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
)

// ErrInvalidLine when csv line is invalid
var ErrInvalidLine = errors.New("invalid line structure")

// ErrInvalidIPv4 when invalid ip (v4) address provided
var ErrInvalidIPv4 = errors.New("invalid IPv4 address")

// ErrInvalidIPv4 when invalid ip v6 address provided
var ErrInvalidIPv6 = errors.New("invalid IPv4 address")

// Range for IPv4
type ip4Range struct {
	start   uint
	end     uint
	country string
}

// Range for v6
type ip6Range struct {
	start   *big.Int
	end     *big.Int
	country string
}

type LoadedIPs struct {
	arrV4 []ip4Range
	arrV6 []ip6Range
}

// Load db-ip.com csv file.
// It must be called only once
//
// Example usage:
//
//	ips, _ := ip2country.Load("db.csv")
//	fmt.Println(ips.GetCountry("1234::"))
//	fmt.Println(ips.GetCountry("255.255.255.255"))
func Load(filepath string) (loaded LoadedIPs, err error) {
	loaded.arrV4 = make([]ip4Range, 0)
	loaded.arrV6 = make([]ip6Range, 0)

	err = loaded.loadFile(filepath)

	return
}

// GetCountry returns the country which ip belongs to
func (ips LoadedIPs) GetCountry(ip string) string {
	if isIPv6(ip) {
		ipNumb := big.NewInt(0)

		err := Ip6ToInt(ip, ipNumb)
		if err != nil {
			return ""
		}

		index := bigBinarySearch(ips.arrV6, ipNumb, 0, len(ips.arrV6)-1)
		if index == -1 {
			return ""
		}

		return ips.arrV6[index].country
	}

	ipNumb, err := Ip4ToInt(ip)
	if err != nil {
		return ""
	}

	index := binarySearch(ips.arrV4, ipNumb, 0, len(ips.arrV4)-1)
	if index == -1 {
		return ""
	}

	return ips.arrV4[index].country
}

// ips.GetCountryMulti is a batch version of GetCountry function
// It allows you to pass many ip addresses as input, and will return countries as output
// the first index of slice is the answer for the first input , the second index for the second input and so on
func (ips LoadedIPs) GetCountryMulti(addrs ...string) []string {
	size := len(addrs)
	answers := make([]string, size)
	var wg sync.WaitGroup
	wg.Add(size)

	for i := 0; i < size; i++ {
		go func(index int) {
			answers[index] = ips.GetCountry(addrs[index])
			wg.Done()
		}(i)
	}
	wg.Wait()

	return answers
}

func (ips *LoadedIPs) loadFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = ips.addRaw(scanner.Text())
		if err != nil {
			return err
		}
	}

	err = scanner.Err()
	return err
}

// function assumes that if the string has the ':' character it is ipv6 and not v4
func isIPv6(ip string) bool {
	return strings.Contains(ip, ":")
}

// accept input string as follows:
//
// "{ip}","{ip}","{country}"
//
// # Supports both IPv4 and IPv6 addresses
//
// IPv6 addresses that have "::" (and not at the end)
// (for example ffff::eeee) will not work (as of yet)
//
// This shouldn't be a problem with dbip csv files
func (ips *LoadedIPs) addRaw(line string) error {
	//replace all double quotations
	line = strings.Replace(line, "\"", "", -1)

	startIP, endIP, country, err := extract(line)
	if err != nil {
		return err
	}

	// Checks if first ip is v4 or v6
	if isIPv6(startIP) {
		startIPnum := big.NewInt(0)
		endIPnum := big.NewInt(0)

		err := Ip6ToInt(startIP, startIPnum)
		if err != nil {
			return err
		}

		err = Ip6ToInt(endIP, endIPnum)
		if err != nil {
			return err
		}

		ips.arrV6 = append(ips.arrV6, ip6Range{startIPnum, endIPnum, country})
		ensureV6Sorted(ips.arrV6)

		return nil
	}

	startIPnum, err := Ip4ToInt(startIP)
	if err != nil {
		return err
	}

	endIPnum, err := Ip4ToInt(endIP)
	if err != nil {
		return err
	}

	ips.arrV4 = append(ips.arrV4, ip4Range{startIPnum, endIPnum, country})
	ensureV4Sorted(ips.arrV4)

	return nil
}

func ensureV4Sorted(arr []ip4Range) {

	i := len(arr) - 1
	temp := arr[i]
	for {
		if i == 0 || arr[i].start >= arr[i-1].start {
			break
		}

		arr[i] = arr[i-1]
		i--
	}
	arr[i] = temp
}

func ensureV6Sorted(arr []ip6Range) {
	i := len(arr) - 1

	temp := arr[i]

	for {
		if i == 0 || arr[i].start.Cmp(arr[i-1].start) >= 0 {
			break
		}

		arr[i] = arr[i-1]
		i--
	}

	arr[i] = temp
}

// Convert an ipv4 ip address to a uint32
func Ip4ToInt(ip string) (uint, error) {

	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0, ErrInvalidIPv4
	}

	var result uint
	var index uint = 3
	for i := 3; i >= 0; i-- {
		ipNumb, err := strconv.Atoi(parts[index])
		if err != nil {
			return 0, err
		}

		result |= uint(ipNumb) << ((uint(3) - index) * uint(8))
		index--
	}

	return result, nil
}

// Convert an ipv6 ip address to big.Int
func Ip6ToInt(ip string, result *big.Int) error {
	rawParts := strings.Split(ip, ":")

	parts := []string{}

	// remove empty strings
	for _, part := range rawParts {
		if part != "" {
			parts = append(parts, part)
		}
	}

	// if "filtered" list is 0 then it is equals to zero
	// ipv6 address "::" should be 0
	if len(parts) == 0 {
		big.NewInt(0)
		return nil
	}

	var index uint = 0

	for i := 0; i < len(parts); i++ {
		ipNumber, err := strconv.ParseInt(parts[i], 16, 64)
		if err != nil {
			return err
		}

		b := big.NewInt(ipNumber)
		b.Lsh(b, (7-index)*16)
		result.Or(result, b)

		index++
	}

	return nil
}

func extract(line string) (string, string, string, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 3 {
		return "", "", "", ErrInvalidLine
	}

	return parts[0], parts[1], parts[2], nil
}

func binarySearch(arr []ip4Range, key uint, start, end int) int {
	for {
		if start > end {
			return -1 //not found
		}

		mid := (start + end) / 2
		if key >= arr[mid].start && key <= arr[mid].end {
			return mid
		}

		if key < arr[mid].start {
			end = mid - 1
		} else if key > arr[mid].end {
			start = mid + 1
		}

	}
}

func bigBinarySearch(arr []ip6Range, key *big.Int, start, end int) int {
	for {
		if start > end {
			return -1
		}

		mid := (start + end) / 2
		// key >= arr.start && key <= arr.end
		if key.Cmp(arr[mid].start) >= 0 && key.Cmp(arr[mid].end) <= 0 {
			return mid
		}

		// key < arr.start
		if key.Cmp(arr[mid].start) < 0 {
			end = mid - 1
		} else if key.Cmp(arr[mid].end) > 0 { // key > arr.end
			start = mid + 1
		}
	}
}
