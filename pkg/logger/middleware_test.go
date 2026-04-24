package logger

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetLocalIPHex(t *testing.T) {
	ip := getLocalIPHex()

	assert.Equal(t, 8, len(ip), "IP十六进制字符串应该是8位")

	matched, _ := regexp.MatchString("^[0-9a-f]{8}$", ip)
	assert.True(t, matched, "IP应该是合法的十六进制格式")
}

func TestGenerateTraceId(t *testing.T) {
	traceId1 := GenerateTraceId()
	traceId2 := GenerateTraceId()

	assert.Equal(t, 30, len(traceId1), "TraceId应该是30位")
	assert.NotEqual(t, traceId1, traceId2, "连续生成的TraceId应该不同")

	matched, _ := regexp.MatchString("^[0-9a-f]{8}[0-9]{13}[0-9]{4}[0-9]{5}$", traceId1)
	assert.True(t, matched, "TraceId格式应该是: IP(8位hex)+时间(13位)+序列(4位)+PID(5位)")

	ipPart := traceId1[:8]
	timePart := traceId1[8:21]
	seqPart := traceId1[21:25]
	pidPart := traceId1[25:]

	assert.NotEqual(t, "00000000", ipPart, "IP部分不应该是默认值")
	assert.Greater(t, len(timePart), 0, "时间戳部分应该有值")

	seq, err := strconv.Atoi(seqPart)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, seq, 1000)
	assert.LessOrEqual(t, seq, 9000)

	pid := os.Getpid()
	expectedPidPart := fmt.Sprintf("%05d", pid)
	assert.Equal(t, expectedPidPart, pidPart, "PID部分应该匹配当前进程ID")
}

func TestGenerateSpanId(t *testing.T) {
	tests := []struct {
		name     string
		parentId string
		expected string
	}{
		{"空字符串生成根节点", "", "0"},
		{"根节点的第一个子节点", "0", "0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSpanId(tt.parentId)
			assert.Equal(t, tt.expected, result)
		})
	}

	spanId1 := GenerateSpanId("")
	spanId2 := GenerateSpanId(spanId1)
	assert.Equal(t, "0", spanId1)
	assert.Equal(t, "0.1", spanId2)

	spanId3 := GenerateSpanId(spanId2)
	assert.Equal(t, "0.2", spanId3)

	spanId4 := GenerateSpanId("0.1")
	assert.Equal(t, "0.2", spanId4)

	spanId5 := GenerateSpanId("0.1.1")
	assert.Equal(t, "0.1.2", spanId5)
}

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)

	middleware := Middleware()
	middleware(c)

	traceId := c.GetString("traceId")
	spanId := c.GetString("spanId")

	assert.NotEmpty(t, traceId, "TraceId不应该为空")
	assert.Equal(t, "0", spanId, "根节点SpanId应该是0")

	assert.Equal(t, 30, len(traceId), "TraceId应该是30位")

	responseTraceId := w.Header().Get("X-Trace-Id")
	responseSpanId := w.Header().Get("X-Span-Id")

	assert.Equal(t, traceId, responseTraceId, "响应头中的TraceId应该与context中一致")
	assert.Equal(t, spanId, responseSpanId, "响应头中的SpanId应该与context中一致")
}

func TestMiddlewareMultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	traceIds := make(map[string]bool)

	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)

		middleware := Middleware()
		middleware(c)

		traceId := c.GetString("traceId")
		assert.False(t, traceIds[traceId], "TraceId应该唯一")
		traceIds[traceId] = true

		assert.Equal(t, "0", c.GetString("spanId"))
	}

	assert.Equal(t, 10, len(traceIds), "应该生成10个不同的TraceId")
}

func TestTraceIdSequenceRollback(t *testing.T) {
	oldSeq := traceSequence

	traceSequence = 8999

	id1 := GenerateTraceId()
	id2 := GenerateTraceId()
	id3 := GenerateTraceId()

	seq1 := id1[21:25]
	seq2 := id2[21:25]
	seq3 := id3[21:25]

	assert.Equal(t, "9000", seq1)
	assert.Equal(t, "1000", seq2)
	assert.Equal(t, "1001", seq3)

	traceSequence = oldSeq
}
