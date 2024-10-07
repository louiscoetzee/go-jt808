package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cuteLittleDevil/go-jt808/protocol"
	"github.com/cuteLittleDevil/go-jt808/protocol/jt808"
	"reflect"
	"testing"
)

func TestReply(t *testing.T) {
	type Handler interface {
		HasReply() bool
		ReplyBody(*jt808.JTMessage) ([]byte, error)
		ReplyProtocol() uint16
	}
	type args struct {
		Handler
		msg2011 string
		msg2013 string
		msg2019 string
	}
	type want struct {
		result2011 string
		result2013 string
		result2019 string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			// 测试的数据使用terminal.go中的CreateTerminalPackage生成
			// 终端和平台的流水号都使用0
			name: "T0x0002 终端-心跳",
			args: args{
				Handler: &T0x0002{},
				msg2013: "7e000200000123456789017fff0a7e",
				msg2019: "7e000240000100000000017299841738ffff027e",
			},
			want: want{
				result2013: "7e8001000501234567890100007fff0002008e7e",
				result2019: "7e8001400501000000000172998417380000ffff000200867e",
			},
		},
		{
			name: "T0x0001 终端-通用应答",
			args: args{
				Handler: &T0x0001{},
				msg2013: "7e000100050123456789017fff007b01c803bd7e",
			},
			want: want{
				result2013: "7e8001000501234567890100007fff0001008d7e",
			},
		},
		{
			name: "T0x0102 终端-鉴权",
			args: args{
				Handler: &T0x0102{},
				msg2013: "7e0102000b01234567890100003137323939383431373338b57e",
				msg2019: "7e0102402f010000000001729984173800000b3137323939383431373338313233343536373839303132333435332e372e31350000000000000000000000000000227e",
			},
			want: want{
				result2013: "7e80010005012345678901000000000102010e7e",
				result2019: "7e80014005010000000001729984173800000000010200877e",
			},
		},
	}
	checkReplyInfo := func(t *testing.T, msg string, handler Handler, expectedResult string) {
		if msg == "" {
			return
		}
		data, _ := hex.DecodeString(msg)
		jtMsg := jt808.NewJTMessage()
		if err := jtMsg.Decode(data); err != nil {
			t.Errorf("Parse() error = %v", err)
			return
		}
		jtMsg.Header.ReplyID = handler.ReplyProtocol()
		if ok := handler.HasReply(); !ok {
			return
		}
		body, _ := handler.ReplyBody(jtMsg)
		got := jtMsg.Header.Encode(body)
		if !reflect.DeepEqual(fmt.Sprintf("%x", got), expectedResult) {
			t.Errorf("ReplyInfo() got = [%x]\n want = [%s]", got, expectedResult)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkReplyInfo(t, tt.args.msg2011, tt.args.Handler, tt.want.result2011)
			checkReplyInfo(t, tt.args.msg2013, tt.args.Handler, tt.want.result2013)
			checkReplyInfo(t, tt.args.msg2019, tt.args.Handler, tt.want.result2019)
		})
	}
}

// 为了覆盖率100%增加的测试 ------------------------------------
func TestT0x0102Reply(t *testing.T) {
	msg := "7e0102402f010000000001729984173800000b3137323939383431373338313233343536373839303132333435332e372e31350000000000000000000000000000227e"
	data, _ := hex.DecodeString(msg)
	jtMsg := jt808.NewJTMessage()
	_ = jtMsg.Decode(data)
	handler := &T0x0102{}
	// 强制错误情况
	jtMsg.Body = nil
	if _, err := handler.ReplyBody(jtMsg); !errors.Is(err, protocol.ErrBodyLengthInconsistency) {
		t.Errorf("T0x0102 ReplyBody() err[%v]", err)
		return
	}
}

func TestT0x002Encode(t *testing.T) {
	handler := &T0x0002{}
	got := handler.Encode()
	if got != nil {
		t.Errorf("T0x002 Encode() got = [%x]", got)
	}
}