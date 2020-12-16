package net

import (
	"net"
	"net/http"
	"testing"
)

func TestHasLocalIpAddr(t *testing.T) {
	for ip, ok := range map[string]bool{
		"":                   false,
		"invalid ip address": false,
		"127.0.0.1":          true,
		"::1":                true,
		"182.56.9.18":        false,
		"192.168.9.18":       true,
		"10.168.9.18":        true,
		"11.168.9.18":        false,
		"172.16.9.18":        true,
		"172.17.9.18":        true,
		"172.18.9.18":        true,
		"172.19.9.18":        true,
		"172.20.9.18":        true,
		"172.21.9.18":        true,
		"172.22.9.18":        true,
		"172.23.9.18":        true,
		"172.24.9.18":        true,
		"172.25.9.18":        true,
		"172.26.9.18":        true,
		"172.27.9.18":        true,
		"172.28.9.18":        true,
		"172.29.9.18":        true,
		"172.30.9.18":        true,
		"172.31.9.18":        true,
		"172.32.9.18":        false,
	} {
		if HasLocalIpAddr(ip) != ok {
			t.Errorf("ip %s",ip)
		}
	}
}

func TestHasLocalIP(t *testing.T) {
	for ip,ok := range map[string]bool{
		"":                   false,
		"invalid ip address": false,
		"127.0.0.1":          true,
		"::1":                true,
		"182.56.9.18":        false,
		"192.168.9.18":       true,
		"10.168.9.18":        true,
		"11.168.9.18":        false,
		"172.16.9.18":        true,
		"172.17.9.18":        true,
		"172.18.9.18":        true,
		"172.19.9.18":        true,
		"172.20.9.18":        true,
		"172.21.9.18":        true,
		"172.22.9.18":        true,
		"172.23.9.18":        true,
		"172.24.9.18":        true,
		"172.25.9.18":        true,
		"172.26.9.18":        true,
		"172.27.9.18":        true,
		"172.28.9.18":        true,
		"172.29.9.18":        true,
		"172.30.9.18":        true,
		"172.31.9.18":        true,
		"172.32.9.18":        false,
	} {
		if HasLocalIP(net.ParseIP(ip)) !=ok {
			t.Errorf("ip %s",ip)
		}
	}
}

func TestRemoteIP(t *testing.T) {
	for _,v := range []struct{
		addr string
		expect string
	} {
		{"101.1.0.4:100", "101.1.0.4"},
		{"101.1.0.4:", "101.1.0.4"},
		{"101.1.0.4", ""},
		{":100", ""},
	} {
		if ok := RemoteIP(&http.Request{RemoteAddr: v.addr}); ok != v.expect {
			t.Errorf("RemoteAddr:%s actual: %s, expect %s", v.addr, ok, v.expect)
		}
	}
}

func TestClientIP(t *testing.T) {
	r := &http.Request{Header: http.Header{}}
	r.Header.Set("X-Real-IP", " 10.10.10.10  ")
	r.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	r.RemoteAddr = "  40.40.40.40:42123 "

	if ip := ClientIP(r); ip != "20.20.20.20" {
		t.Errorf("actual: 20.20.20.20, expected:%s", ip)
	}

	r.Header.Del("X-Forwarded-For")
	if ip := ClientIP(r); ip != "10.10.10.10" {
		t.Errorf("actual: 10.10.10.10, expected:%s", ip)
	}

	r.Header.Set("X-Forwarded-For", "30.30.30.30  ")
	if ip := ClientIP(r); ip != "30.30.30.30" {
		t.Errorf("actual: 30.30.30.30, expected:%s", ip)
	}

	r.Header.Del("X-Forwarded-For")
	r.Header.Del("X-Real-IP")
	if ip := ClientIP(r); ip != "40.40.40.40" {
		t.Errorf("actual: 40.40.40.40, expected:%s", ip)
	}

	r.RemoteAddr = "50.50.50.50"
	if ip := ClientIP(r); ip != "" {
		t.Errorf("ip: 50.50.50.50")
	}
}

func TestClientPublicIP(t *testing.T) {
	for _, v := range []struct {
		xForwardedFor string
		remoteAddr    string
		expected      string
	}{
		{"10.3.5.45, 21.45.9.1", "101.1.0.4:100", "21.45.9.1"},
		{"101.3.5.45, 21.45.9.1", "101.1.0.4:100", "101.3.5.45"},
		{"", "101.1.0.4:100", "101.1.0.4"},
		{"21.45.9.1", "101.1.0.4:100", "21.45.9.1"},
		{"21.45.9.1, ", "101.1.0.4:100", "21.45.9.1"},
		{"192.168.5.45, 210.45.9.1, 89.5.6.1", "101.1.0.4:100", "210.45.9.1"},
		{"192.168.5.45, 172.24.9.1, 89.5.6.1", "101.1.0.4:100", "89.5.6.1"},
		{"192.168.5.45, 172.24.9.1", "101.1.0.4:100", "101.1.0.4"},
		{"192.168.5.45, 172.24.9.1", "101.1.0.4:5670", "101.1.0.4"},
	} {
		if actual := ClientPublicIP(&http.Request{
			Header: http.Header{
				"X-Forwarded-For": []string{v.xForwardedFor},
			},
			RemoteAddr: v.remoteAddr,
		}); actual != v.expected {
			t.Errorf("IsxForwardedFor:%s, remoteAddr:%s, client ip Should Equal %s", v.xForwardedFor, v.remoteAddr, v.expected)
		}
	}

	r := &http.Request{Header: http.Header{}}
	r.Header.Set("X-Real-IP", " 10.10.10.10  ")
	r.Header.Set("X-Forwarded-For", " 172.17.40.152, 192.168.5.45")
	r.RemoteAddr = "  40.40.40.40:42123 "

	if ip := ClientPublicIP(r); ip != "40.40.40.40" {
		t.Errorf("actual:40.40.40.40, expected:%s", ip)
	}

	r.Header.Set("X-Real-IP", " 50.50.50.50  ")
	if ip := ClientPublicIP(r); ip != "50.50.50.50" {
		t.Errorf("actual:50.50.50.50, expected:%s", ip)
	}

	r.Header.Del("X-Real-IP")
	r.Header.Del("X-Forwarded-For")
	r.RemoteAddr = "  127.0.0.1:42123 "
	if ip := ClientPublicIP(r); ip != "" {
		t.Errorf("ip: 127.0.0.1")
	}
}

func TestIpStringToLong(t *testing.T) {
	for _, v := range []struct {
		ip   string
		long uint
	}{
		{"127.0.0.1", 2130706433},
		{"0.0.0.0", 0},
		{"255.255.255.255", 4294967295},
		{"192.168.1.1", 3232235777},
	} {
		expected, err := IpStringToLong(v.ip)
		if err != nil {
			t.Errorf("ip:%s long:%d err:%v", v.ip, v.long, err)
		}

		if expected != v.long {
			t.Errorf("ip:%s long:%d != expected:%d", v.ip, v.long, expected)
		}
	}

	for _, ip := range []string{
		"",
		"invalid ip address",
		"::1",
	} {
		_, err := IpStringToLong(ip)
		if err == nil {
			t.Errorf("ip:%s invalid IP passes", ip)
		}
	}
}

func TestLongToIpString(t *testing.T) {
	for _, v := range []struct {
		ip   string
		long uint
	}{
		{"127.0.0.1", 2130706433},
		{"0.0.0.0", 0},
		{"255.255.255.255", 4294967295},
		{"192.168.1.1", 3232235777},
	} {
		expected, err := LongToIpString(v.long)
		if err != nil {
			t.Errorf("ip:%s long:%d err:%v", v.ip, v.long, err)
		}

		if expected != v.ip {
			t.Errorf(" long:%d ip:%s != expected:%s", v.long, v.ip, expected)
		}
	}

	// 在64位机器上运行，否者输入值将超过限制
	if 32<<(^uint(0)>>63) == 64 {
		_, err := LongToIpString(1<<64 - 1)
		if err == nil {
			t.Errorf("long:%d out of range", uint64(1<<64-1))
		}
	}
}

func TestIpToLong(t *testing.T) {
	for _, v := range []struct {
		ip   string
		long uint
	}{
		{"127.0.0.1", 2130706433},
		{"0.0.0.0", 0},
		{"255.255.255.255", 4294967295},
		{"192.168.1.1", 3232235777},
	} {
		expected, err := IpToLong(net.ParseIP(v.ip))
		if err != nil {
			t.Errorf("ip:%s long:%d err:%v", v.ip, v.long, err)
		}

		if expected != v.long {
			t.Errorf("ip:%s long:%d != expected:%d", v.ip, v.long, expected)
		}
	}

	for _, ip := range []string{
		"",
		"invalid ip address",
		"::1",
	} {
		_, err := IpToLong(net.ParseIP(ip))
		if err == nil {
			t.Errorf("ip:%s invalid IP passes", ip)
		}
	}
}

func TestLongToIp(t *testing.T) {
	for _, v := range []struct {
		ip   string
		long uint
	}{
		{"127.0.0.1", 2130706433},
		{"0.0.0.0", 0},
		{"255.255.255.255", 4294967295},
		{"192.168.1.1", 3232235777},
	} {
		expected, err := LongToIp(v.long)
		if err != nil {
			t.Errorf("ip:%s long:%d err:%v", v.ip, v.long, err)
		}

		if expected.String() != v.ip {
			t.Errorf(" long:%d ip:%s != expected:%s", v.long, v.ip, expected.String())
		}
	}

	// 在64位机器上运行，否者输入值将超过限制
	if 32<<(^uint(0)>>63) == 64 {
		_, err := LongToIp(1<<64 - 1)
		if err == nil {
			t.Errorf("long:%d out of range", uint64(1<<64-1))
		}
	}
}