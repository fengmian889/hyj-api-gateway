package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/fengmian889/hyj-api-gateway/internal/pkg/contextx"
)

// ProxyHandler 处理API请求转发
type ProxyHandler struct {
	targetURL *url.URL
}

// NewProxyHandler 创建新的代理处理器
func NewProxyHandler(target string) (*ProxyHandler, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return &ProxyHandler{
		targetURL: targetURL,
	}, nil
}

// Handle 执行代理请求
func (p *ProxyHandler) Handle(c *gin.Context) {
	// 从上下文中获取 RequestID
	requestID := contextx.RequestID(c.Request.Context())
	if requestID == "" {
		requestID = "unknown"
	}

	// 记录请求日志
	zap.L().Info("["+requestID+"] Proxy API called",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("clientIP", c.ClientIP()),
		zap.String("target", p.targetURL.String()),
	)

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(p.targetURL)

	// 修改请求的Host头
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 修改请求URL以匹配目标服务器
		req.URL.Scheme = p.targetURL.Scheme
		req.URL.Host = p.targetURL.Host
		req.Host = p.targetURL.Host

		// 保留原始路径，确保正确映射到目标服务器
		// 这会将如 /api/voice/invite 的请求转发到目标服务器的 /api/voice/invite
		req.URL.Path = c.Request.URL.Path
		req.URL.RawQuery = c.Request.URL.RawQuery
	}

	// 修改响应处理
	proxy.ModifyResponse = func(resp *http.Response) error {
		// 可以在这里修改响应，例如添加额外的头信息
		return nil
	}

	// 修改错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		zap.L().Error("Proxy error",
			zap.String("error", err.Error()),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error": "Proxy error", "message": "Failed to reach target server"}`))
	}

	// 执行代理
	proxy.ServeHTTP(c.Writer, c.Request)
}

// VoiceInviteProxyHandler 处理 /api/voice/invite 请求转发
func VoiceInviteProxyHandler(c *gin.Context) {
	proxyHandler, err := NewProxyHandler("http://go-cti.duyansoft.com")
	if err != nil {
		zap.L().Error("Failed to create proxy handler", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy configuration error"})
		return
	}
	proxyHandler.Handle(c)
}

// VoiceCallLogProxyHandler 处理 /api/voice/call/log 请求转发
func VoiceCallLogProxyHandler(c *gin.Context) {
	proxyHandler, err := NewProxyHandler("http://go-cti.duyansoft.com")
	if err != nil {
		zap.L().Error("Failed to create proxy handler", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy configuration error"})
		return
	}
	proxyHandler.Handle(c)
}
