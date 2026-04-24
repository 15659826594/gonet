package logger

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	traceSequence uint32 = 1000
	localIP       string
)

func init() {
	localIP = getLocalIPHex()
}

// getLocalIPHex 获取本机IP并转换为8位十六进制字符串
func getLocalIPHex() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "00000000"
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip := ipNet.IP.To4()
				return fmt.Sprintf("%02x%02x%02x%02x", ip[0], ip[1], ip[2], ip[3])
			}
		}
	}
	return "00000000"
}

// GenerateTraceId 生成TraceId: IP(8位) + 时间戳(13位) + 自增序列(4位) + 进程ID(5位)
// 示例: 0ad1348f1403169275002100356696
func GenerateTraceId() string {
	ipPart := localIP

	timestamp := time.Now().UnixMilli()
	timePart := fmt.Sprintf("%013d", timestamp)

	seq := atomic.AddUint32(&traceSequence, 1)
	if seq > 9000 {
		atomic.StoreUint32(&traceSequence, 1000)
		seq = 1000
	}
	seqPart := fmt.Sprintf("%04d", seq)

	pid := os.Getpid()
	pidPart := fmt.Sprintf("%05d", pid)

	return ipPart + timePart + seqPart + pidPart
}

// GenerateSpanId 生成SpanId，表示调用链路树中的位置
// 根节点: 0
// 第一层调用: 0.1, 0.2, 0.3...
// 第二层调用: 0.1.1, 0.1.2, 0.2.1...
func GenerateSpanId(parentSpanId string) string {
	if parentSpanId == "" {
		return "0"
	}

	if parentSpanId == "0" {
		return "0.1"
	}

	parts := strings.Split(parentSpanId, ".")

	lastPart := parts[len(parts)-1]
	num, err := strconv.Atoi(lastPart)
	if err != nil {
		num = 0
	}
	num++

	parts[len(parts)-1] = strconv.Itoa(num)
	return strings.Join(parts, ".")
}

// Middleware 链路追踪中间件
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := GenerateTraceId()
		spanId := GenerateSpanId("")

		c.Set("traceId", traceId)
		c.Set("spanId", spanId)
		c.Writer.Header().Set("X-Trace-Id", traceId)
		c.Writer.Header().Set("X-Span-Id", spanId)
		c.Next()
	}
}
