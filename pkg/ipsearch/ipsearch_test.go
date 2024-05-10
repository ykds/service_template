package ipsearch

import (
	"fmt"
	"testing"
)

func TestGetIPLocation(t *testing.T) {
	InitIpDatabase("C:\\Users\\1\\Documents\\国内精华版-202404-354425\\qqzeng-ip-china-utf8.dat")
	ipaddress := GetIPAddress("203.72.97.75")
	fmt.Printf("%+v\n", ipaddress)
}
