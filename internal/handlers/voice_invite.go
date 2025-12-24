// Package handlers 提供 HTTP 请求处理函数
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/fengmian889/hyj-calllog/internal/pkg/contextx"
)

// VoiceInviteHandler 处理 /api/voice/invite 请求，返回 XML 模板
func VoiceInviteHandler(c *gin.Context) {
	// 从上下文中获取 RequestID
	requestID := contextx.RequestID(c.Request.Context())
	if requestID == "" {
		requestID = "unknown" // 如果没有找到则使用默认值
	}

	// 记录请求日志
	zap.L().Info("["+requestID+"] Voice invite API called",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("clientIP", c.ClientIP()),
		zap.String("caller", c.Query("Caller-Caller-ID-Number")),
		zap.String("callee", c.Query("Caller-Destination-Number")),
	)

	uuid := uuid.New().String()
	// 设置响应头为 XML 类型
	caller := c.Query("Caller-Caller-ID-Number")
	callee := c.Query("Caller-Destination-Number")

	// 默认的 XML 模板内容
	xmlResponse := `<?xml version="1.0" encoding="UTF-8"?>
<document type="freeswitch/xml">
    <section name="dialplan" description="Dial Plan For FreeSwitch">
        <context name="public">
            <extension name="outbound">
                <condition>
                    <action application="set" data="LEG=A"/>
                    <action application="set" data="SKIP_XML_CDR=true" />
                    <action application="limit" data="hash outbound 100054 1000 !USER_BUSY"/>
                    <action application="export" data="rtp_autoflush_during_bridge=true" inline="true"/>
                    <action application="export" data="rtp_jitter_buffer_during_bridge=true"/>
                    <action application="export" data="jitterbuffer_msec=200:1000"/>
                    <action application="export" data="CALL_TYPE=OUTBOUND" inline="true"/>
                    <action application="export" data="AGENTID=10000345" inline="true"/>
                    <action application="export" data="TEAMID=10000029" inline="true"/>
                    <action application="export" data="CALLEE_SUFFIX=3965" inline="true"/>
                    <action application="export" data="CALLER=` + caller + `" inline="true"/>
                    <action application="export" data="CALLEE=` + callee + `" inline="true"/>
                    <action application="export" data="ORGID=100054" inline="true"/>

                        <action application="export" data="_DATE=${strftime(%Y-%m-%d)}" inline="true"/>
                        <action application="export" data="_TIME=${strftime(%Y-%m-%d %H:%M)}" inline="true"/>
                        <action application="export" data="sip_ph_X-UUID=` + uuid + `"/>
                        <action application="export" data="sip_ph_X-Local-Number=` + caller + `"/>
                    <action application="export" data="sip_h_X-LOCAL-UUID=` + uuid + `"/>
                    <action application="export" data="CA_LEG_UUID=` + uuid + `" inline="true"/>
                    <action application="export" data="sip_h_X-CTI_ORGID=100054" inline="true"/>
                    <action application="export" data="codec_string=PCMA"/>
                        <action application="set" data="sip_h_X-EarlyMediaRecordEnabled=true"/>
                        <action application="log" data="NOTICE Will record with PCMA format"/>
                    <action application="export" data="RECORD_HANGUP_ON_ERROR=true"/>

                    <action application="export"
                            data="nolocal:record_post_process_exec_app=bgsystem:/bin/dash /san/script/record_post_process_ob.sh 100054 2025-12-07 ` + uuid + ` "/>
                    <action application="export"
                            data="nolocal:execute_on_answer=record_session /tmp/recording/100054/2025-12-07/` + uuid + `"/>
                        <!-- Music on Hold -->
						<action application="playback" data="$${sounds_dir}/music/du.wav"/>
                        <action application="export" data="instant_ringback=true"/>
                        <action application="export" data="ringback=%(1000, 4000, 450.0)"/>
                        <action application="bind_digit_action" data="music_on_hold,A,exec:soft_hold,B"/>
                        <!-- Bridge -->
                        <action application="export" data="hangup_after_bridge=false"/>
                        <action application="sched_hangup" data="+7200"/>
                    <action application="export" data="effective_caller_id_number=cv1test"/>
                    <action application="bridge"
                            data="[origination_uuid=` + uuid + `-b,LEG=B,absolute_codec_string='PCMA',CALL_TYPE=${CALL_TYPE}]sofia/external/` + callee + `@172.19.0.69:60051"/>
                </condition>
            </extension>
        </context>
    </section>
</document>`

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, xmlResponse)
}
