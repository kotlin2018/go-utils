package net

import (
	"errors"
	"math"
	"net"
	"net/http"
	"strings"
)

// 检测 IP地址字符串是否是内网地址
func HasLocalIpAddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// 检测 IP 地址是否是内网地址,默认是内网地址
func HasLocalIP(ip net.IP) bool {
	// 如果 ip = 127.0.0.1
	if ip.IsLoopback(){
		return true
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	// 返回常见内容地址
	return ip4[0] == 10 || 								 // 10.0.0.0./8
		(ip4[0] == 172 && ip4[1] >=16 && ip4[1] <=31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || 			 // 169.254.0.0./16
		(ip4[0] == 192 && ip4[1] == 168) 				 // 192.168.0.0/16
}

// 获取客户端的IP
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作
func ClientIP(r *http.Request) string {
	x := r.Header.Get("X-Forwarded-For")

	if ip := strings.TrimSpace(strings.Split(x, ",")[0]);ip != ""{
		return ip
	}

	if ip := strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip != "" {
		return ip
	}

	if ip, _,err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr));err == nil {
		return ip
	}
	return ""
}

// 获取客户端公网IP
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _,ip = range strings.Split(r.Header.Get("X-Forwarded-For"),",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIpAddr(ip) {
			return ip
		}
	}

	if ip = strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip !="" &&!HasLocalIpAddr(ip) {
		return ip
	}

	if ip,_ ,err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr));err == nil && !HasLocalIpAddr(ip){
		return ip
	}
	return ""
}

// 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func RemoteIP(r *http.Request) string {
	if ip,_,err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// 把ip字符串转为数值
func IpStringToLong(ip string) (uint,error) {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0, errors.New("ipv4 格式化错误")
	}
	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24,nil
}

// 把数值转为ip字符串
func LongToIpString(i uint)(string,error) {
	if i > math.MaxUint32 {
		return "", errors.New("超出了ipv4的范围")
	}
	ip := make(net.IP,net.IPv4len)
	ip[0] = byte(i >>24)
	ip[1] = byte(i >>16)
	ip[2] = byte(i >>8)
	ip[3] = byte(i)
	return ip.String(),nil
}

// 把net.IP转为数值
func IpToLong(ip net.IP) (uint,error) {
	b := ip.To4()
	if b == nil {
		return 0, errors.New("ipv4 格式化错误")
	}
	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// 把数值转为net.IP
func LongToIp(i uint)(net.IP,error) {
	if i > math.MaxUint32 {
		return nil, errors.New("超出了ipv4的范围")
	}
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)
	return ip, nil
}
