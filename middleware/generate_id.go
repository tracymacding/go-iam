package middleware

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/go-iam/context"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GenerateRequestIdMiddleware struct {
}

var serverIp = ""

func init() {
	interfaces, err := net.InterfaceAddrs()
	if err == nil {
		for _, iface := range interfaces {
			if ipnet, ok := iface.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					serverIp = ipnet.IP.String()
					break
				}
			}
		}
	}
}

var (
	randLen = 8
	tsLen   = 8
)

func ipV4ToInt32(ip string) int32 {
	parts := strings.Split(ip, ".")
	b0, _ := strconv.Atoi(parts[0])
	b1, _ := strconv.Atoi(parts[1])
	b2, _ := strconv.Atoi(parts[2])
	b3, _ := strconv.Atoi(parts[3])

	var sum int32
	sum += int32(b0) << 24
	sum += int32(b1) << 16
	sum += int32(b2) << 8
	sum += int32(b3)
	return sum
}

func int32ToIPV4(ip int32) string {
	var bytes [4]byte
	bytes[0] = byte(ip & 0xFF)
	bytes[1] = byte((ip >> 8) & 0xFF)
	bytes[2] = byte((ip >> 16) & 0xFF)
	bytes[3] = byte((ip >> 24) & 0xFF)
	return fmt.Sprintf("%d.%d.%d.%d", bytes[3], bytes[2], bytes[1], bytes[0])
}

func generateRequestId(remoteIp string) string {
	buf := new(bytes.Buffer)
	buf.Grow(randLen + tsLen + len([]byte(remoteIp)))

	rand.Seed(time.Now().UTC().UnixNano())
	binary.Write(buf, binary.LittleEndian, rand.Int63())
	binary.Write(buf, binary.LittleEndian, time.Now().Unix())
	binary.Write(buf, binary.LittleEndian, ipV4ToInt32(remoteIp))
	return hex.EncodeToString(buf.Bytes())
}

func (m *GenerateRequestIdMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	context.Set(req, "request_id", generateRequestId(serverIp))
	context.Set(req, "request_start", time.Now().Format(time.RFC3339))
	context.Set(req, "request_start_unix", time.Now().UnixNano())
	next(w, req)
}
