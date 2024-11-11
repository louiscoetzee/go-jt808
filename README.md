# go-jt808

- This project is better supported with the goal of facilitating secondary development. You can accomplish corresponding functions through various custom events. Common use cases are listed below:
3. jt1078 Video [Details](./example/jt1078/README.md)

```txt
jt808 server jt1078 server emulator on the local machine
Platform sends 0x9101 command, emulator starts streaming
```
| Streaming Service | Language | Online Playback URL | Description |
|-------------------|----------|----------------------|-------------|
| LAL | Go | http://49.234.235.7:8080/live/295696659617_1.flv | [Click for details](./example/jt1078/README.md#lal) |
| sky-java | Java | Requires deployment, HTTP request within 10 seconds to pull the stream. Reference format as follows: <br/> http://222.244.144.181:7777/video/1001-1-0-0.live.mp4 | [Click for details](./example/jt1078/README.md#sky-java) |

2. Storing Latitude and Longitude [Code Reference](./example/simulator/server/main.go)
```txt
jt808 server, emulator, message queue, and database all running on a 2-core 4GB Tencent Cloud server
Tested saving 5000 records per second, saved nearly 100 million latitude and longitude in about 5.5 hours
```

3. Platform Sends Commands to Terminal [Code Reference](./example/protocol/active_reply/main.go)
```txt
Active commands sent to device to receive responses
```

4. Protocol Interaction Details [Code Reference](./example/protocol/register/main.go)
```txt
Using a custom emulator, you can easily generate test messages with detailed descriptions
```

5. Custom Protocol Extensions [Code Reference](./example/protocol/custom_parse/main.go)
```txt
Handle custom additional information to obtain desired extension content
```

---

- Observing Fei Ge's single-machine TCP million concurrency, curious about performance with data, planning to test with data during the National Day holiday
- Performance testing on a single machine [2-core 4GB machine], concurrent 100k+, saving 400 million+ latitude and longitude daily [Details](./README.md#save)
- Supports JT808 (2011/2013/209) JT1078 (requires other streaming media services), supports segmentation and automatic retransmission

| Feature | Description |
|:-------:|-------------|
| Safe and Reliable | Core protocol components tested with 100% coverage, pure native Go implementation (no dependencies) |
| Simple and Elegant | Core code under 1000 lines, no locks used, utilizes only channels |
| Easy to Extend | Convenient for secondary development, includes JT1078 streaming media integration, storing latitude and longitude, and other use cases |

---

## Quick Start
```go
package main

import (
	"github.com/cuteLittleDevil/go-jt808/service"
	"log/slog"
	"os"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))
	slog.SetDefault(logger)
}

func main() {
	goJt808 := service.New(
		service.WithHostPorts("0.0.0.0:8080"),
		service.WithNetwork("tcp"),
	)
	goJt808.Run()
}
```

---
- Currently (2024-10-01), I personally find the Go language versions before this date unsatisfactory and do not recommend referencing them. Recommended reference materials are listed below:

## Reference Materials
| Project Name | Language | Date | Star Count | Link |
|--------------|----------|------|------------|------|
| JT808 | C# | 2024-10-01 | 534 | [JT808 C#](https://github.com/SmallChi/JT808.git) |
| jt808-server | Java | 2024-10-01 | 1.4k+ | [JT808 Java](https://gitee.com/yezhihao/jt808-server.git) |

- [Fei Ge's Development Kung Fu](https://github.com/yanfeizhang/coder-kung-fu?tab=readme-ov-file)
- [Protocol Documents (PDF Compilation)](https://gitee.com/yezhihao/jt808-server/tree/master/协议文档)
- [Protocol Documents (Official Website)](https://jtst.mot.gov.cn/hb/search/stdHBView?id=a3011cd31e6602ec98f26c35329e88e4)
- [Protocol Parsing Website](https://jttools.smallchi.cn/jt808)
- [BCD to DEC Encoding](https://github.com/deatil/lakego-admin/tree/main/pkg/lakego-pkg/go-encoding/bcd)
- [LAL Streaming Media Documentation](https://pengrl.com/lal/#/streamurllist)

## Performance Testing
- Java emulator (QQ group download 373203450)
- Go emulator [Click for details](./example/simulator/client/main.go#go-emulator)

### Connection Count Testing
[Click for details](./example/simulator/README.md#online)

- 2 cloud servers each running 50k+ clients, totaling 100k+

| Server Version | Scenario | Concurrent Count | Server Configuration | Server Resource Usage | Description |
|:--------------:|:--------:|:----------------:|:--------------------:|:----------------------:|:------------:|
| v0.3.0 | Connection Count Test | 100k+ | 2-core 4GB | 120%+ CPU, 1.7G Memory | Server and emulator started on 10.0.16.5 <br/> Emulator started on 10.0.16.14 |

<h3 id="save"> Latitude and Longitude Storage Simulation Test </h3>

[Click for details](./example/simulator/README.md#save)

- Save process lost some data due to channel queue overflow (tested channel queue size was 100)
- Saved 100 million lost 826 records, saved 432 million lost 1216 records (tested twice)

| Server Version | Client | Server Configuration | Server Resource Usage | Description |
|:--------------:|:------:|:--------------------:|:----------------------:|:-----------:|
| v0.3.0 | 10k Go emulators | 2-core 4GB | 35% CPU, 180.4MB Memory | Saving 100 million latitude and longitude at 5000 per second <br/> Actually saved 99,999,174 with a success rate of 99.999% |

| Service | CPU | Memory | Description |
|:-------:|:---:|:-----:|:-----------:|
| server | 35% | 180.4MB | 808 server |
| client | 23% | 196MB | Emulator clients |
| save | 18% | 68.8MB | Data storage service |
| nats-server | 20% | 14.8MB | Message queue |
| taosadapter | 37% | 124.3MB | TDengine database adapter |
| taosd | 15% | 124.7MB | TDengine database |

## Use Cases

### 1. Protocol Handling

#### 1.1 Protocol Parsing
```go
func main() {
	t := terminal.New(terminal.WithHeader(consts.JT808Protocol2013, "1")) // Custom emulator, 2013 version
	data := t.CreateDefaultCommandData(consts.T0100Register) // Generate default registration command 0x0100
	fmt.Println(fmt.Sprintf("Emulator generated [%x]", data))

	jtMsg := jt808.NewJTMessage()
	_ = jtMsg.Decode(data) // Parse fixed request header

	var t0x0100 model.T0x0100 // Parse body data into struct
	_ = t0x0100.Parse(jtMsg)
	fmt.Println(jtMsg.Header.String())
	fmt.Println(t0x0100.String())
}
```

Partial Output [Output Details](./example/protocol/README.md#register)
```txt
Emulator generated [7e010000300000000000010001001f006e63643132337777772e3830382e636f6d0000000000000000003736353433323101b2e2413132333435363738797e]
[0100] Message ID:[256] [Terminal - Registration]
Message Body Property Object: {
        [0000000000110000] Message Body Property Object:[48]
        Version Number:[JT2013]
        [bit15] [0]
        [bit14] Protocol Version Flag:[0]
        [bit13] Is Segmented:[false]
        [bit10-12] Encryption Flag:[0] 0-None, 1-RSA
        [bit0-bit9] Message Body Length:[48]
}
[000000000001] Terminal Phone Number:[1]
[0001] Message Serial Number:[1]
Data Body Object:{
        [001f006e6364313233] Manufacturer ID (5):[cd123]
        [7777772e3830382e636f6d000000000000000000] Terminal Model (20):[www.808.com]
        [37363534333231] Terminal ID (7):[7654321]
        [01] License Plate Color:[1]
        [b2e2413132333435363738] License Plate Number:[Test A12345678]
}
```

#### 1.2 Custom Protocol Extensions (Additional Information)

Custom parsing extension 0x33 as an example. Key code is as follows:
```go
func (l *Location) Parse(jtMsg *jt808.JTMessage) error {
	l.T0x0200AdditionDetails.CustomAdditionContentFunc = func(id uint8, content []byte) (model.AdditionContent, bool) {
		if id == uint8(consts.A0x01Mile) {
			l.customMile = 100
		}
		if id == 0x33 {
			value := content[0]
			l.customValue = value
			return model.AdditionContent{
				Data:        content,
				CustomValue: value,
			}, true
		}
		return model.AdditionContent{}, false
	}
	return l.T0x0200.Parse(jtMsg)
}

func (l *Location) OnReadExecutionEvent(message *service.Message) {
	if v, ok := tmp.Additions[consts.A0x01Mile]; ok {
		fmt.Println(fmt.Sprintf("Mileage [%d] Custom Auxiliary Mileage [%d]", v.Content.Mile, tmp.customMile))
	}
	id := consts.JT808LocationAdditionType(0x33)
	if v, ok := tmp.Additions[id]; ok {
		fmt.Println("Custom Unknown Information Extension", v.Content.CustomValue, tmp.customValue)
	}
}
```

Partial Output [Output Details](./example/protocol/README.md#custom)
```txt
Mileage [11] Custom Auxiliary Mileage [100]
Custom Unknown Information Extension 32 32
```

#### 1.3 Platform Sends Parameters to Terminal (8104 Query Terminal Parameters)

Key code is as follows:
```go
	replyMsg := goJt808.SendActiveMessage(&service.ActiveMessage{
		Key:              phone,                           // Defaults to using phone number as the unique key to find the corresponding terminal's TCP connection
		Command:          consts.P8104QueryTerminalParams, // Command to send
		Body:             nil,                             // Body data to send, 8104 has no body
		OverTimeDuration: 3 * time.Second,                 // Timeout duration; fails if the device does not respond within this time
	})
	var t0x0104 model.T0x0104
	if err := t0x0104.Parse(replyMsg.JTMessage); err != nil {
		panic(err)
	}
	fmt.Println(t0x0104.String())
```

Partial Output [Output Details](./example/protocol/README.md#active_reply)
```txt
Data Body Object:{
        [0003] Response Message Serial Number:[3]
        [5b] Number of Response Parameters:[91]
        Terminal Query Parameters:

        {
                [0001] Terminal Parameter ID:1 Terminal Heartbeat Interval, in seconds (s)
                Parameter Length [4] Exists [true]
                [0000000a] Parameter Value:[10]
        }
		...
        {
                [0110] Terminal Parameter ID:272 CAN Bus ID Collection Settings:
                Parameter Length [8] Exists [true]
                [0000000000000101] Parameter Value:[[0 0 0 0 0 0 1 1]]
        }
        Unknown Terminal Parameter IDs:[33 117 118 119 121 122 123 124]
}
```

### 2. Custom Latitude and Longitude Storage
[Click for details](./example/simulator/server/main.go)
- Custom implementation of 0x0200 message handling to send data to NATS. Key code is as follows:
```go
type T0x0200 struct {
	model.T0x0200
}

func (t *T0x0200) OnReadExecutionEvent(message *service.Message) {
	var t0x0200 model.T0x0200
	if err := t0x0200.Parse(message.JTMessage); err != nil {
		fmt.Println(err)
		return
	}
	location := shared.NewLocation(message.Header.TerminalPhoneNo, t0x0200.Latitude, t0x0200.Longitude)
	if err := mq.Default().Pub(shared.SubLocation, location.Encode()); err != nil {
		fmt.Println(err)
		return
	}
}

func (t *T0x0200) OnWriteExecutionEvent(_ service.Message) {}
```

- Start with custom 0x0200 message handling

```go
	goJt808 := service.New(
		service.WithHostPorts("0.0.0.0:8080"),
		service.WithNetwork("tcp"),
		service.WithCustomHandleFunc(func() map[consts.JT808CommandType]service.Handler {
			return map[consts.JT808CommandType]service.Handler{
				consts.T0200LocationReport:      &T0x0200{},
			}
		}),
	)
	goJt808.Run()
```

### 3. jt808 Attachment Upload

### 4. jt1078 Related

#### 4.1 Using LAL for Streaming Media Service

- Convert 1078 format stream to the corresponding format and insert it into the LAL service. Core code reference:
```go
func (j *jt1078) createStream(name string) chan<- *Packet {
	...
	ch := make(chan *Packet, 100)
	go func(session logic.ICustomizePubSessionContext, ch <-chan *Packet) {
		for v := range ch {
				...
				switch v.Flag.PT {
				case PTG711A:
					tmp.PayloadType = base.AvPacketPtG711A
				case PTG711U:
					tmp.PayloadType = base.AvPacketPtG711U
				case PTH264:
				case PTH265:
					tmp.PayloadType = base.AvPacketPtHevc
				default:
					slog.Warn("Unknown type",
						slog.Any("pt", v.Flag.PT))
				}
				if err := session.FeedAvPacket(tmp); err != nil {
					slog.Warn("session.FeedAvPacket",
						slog.Any("err", err))
				}
			}
		}
	}(session, ch)
}
```

## Protocol Integration Completion Status

### JT808 Terminal Communication Protocol Message Mapping Table

| No. | Message ID | Completion Status | Test Status | Message Body Name | 2019 Version | 2011 Version |
|-----|------------|--------------------|-------------|-------------------|--------------|--------------|
| 1   | 0x0001     | ✅ | ✅ | Terminal General Response | | |
| 2   | 0x8001     | ✅ | ✅ | Platform General Response | | |
| 3   | 0x0002     | ✅ | ✅ | Terminal Heartbeat | | |
| 5   | 0x0100     | ✅ | ✅ | Terminal Registration | Modified | Modified |
| 4   | 0x8003     | ✅ | ✅ | Request for Re-transmission and Segmentation | | Added |
| 6   | 0x8100     | ✅ | ✅ | Platform Registration Response | | |
| 9   | 0x8103     | ✅ | ✅ | Set Terminal Parameters | Modified and Added | Modified |
| 8   | 0x0102     | ✅ | ✅ | Terminal Authentication | Modified | |
| 10  | 0x8104     | ✅ | ✅ | Platform Query Terminal Parameters | | |
| 11  | 0x0104     | ✅ | ✅ | Query Terminal Parameters Response | | |
| 18  | 0x0200     | ✅ | ✅ | Location Information Report | Added Additional Information | Modified |
| 49  | 0x0704     | ✅ | ✅ | Bulk Upload of Positioning Data | Modified | Added |

### JT1078 Extended JT808 Protocol Message Mapping Table

| No. | Message ID | Completion Status | Test Status | Message Body Name |
|-----|------------|-------------------|-------------|-------------------|
| 13  | 0x1003     | ✅ | ✅ | Terminal Upload Audio/Video Properties |
| 14  | 0x1005     | ✅ | ✅ | Terminal Upload Passenger Flow |
| 15  | 0x1205     | ✅ | ✅ | Terminal Upload Audio/Video Resource List |
| 16  | 0x1206     | ✅ | ✅ | File Upload Completion Notification |
| 17  | 0x9003     | ✅ | ✅ | Platform Query Terminal Audio/Video Properties |
| 18  | 0x9101     | ✅ | ✅ | Platform Real-time Audio/Video Transmission Request |
| 19  | 0x9102     | ✅ | ✅ | Platform Audio/Video Real-time Transmission Control |
| 20  | 0x9105     | ✅ | ✅ | Platform Real-time Audio/Video Transmission Status Notification |
| 21  | 0x9201     | ✅ | ✅ | Platform Remote Video Playback Request |
| 22  | 0x9202     | ✅ | ✅ | Platform Remote Video Playback Control |
| 23  | 0x9205     | ✅ | ✅ | Platform Query Resource List |
| 24  | 0x9206     | ✅ | ✅ | Platform File Upload Command |
| 25  | 0x9207     | ✅ | ✅ | Platform File Upload Control |
