# JT1078 Streaming Media

<h2 id="rtvs-dev"> RTVS Terminal Emulator </h2>

```
rtvsdev (1078 terminal emulator docker version)
Run on command line
docker run --restart always -p 5288:80 -d vanjoge/rtvsdevice
Then access your http://IP:5288
```

<h2 id="lal"> LAL Streaming Media Service </h2>

1. Use the emulator's default data to continuously push to the LAL service
2. Online playback address http://49.234.235.7:8080/live/295696659617_1.flv
- [LAL Official Documentation](https://pengrl.com/lal/#/streamurllist)
- [Code Reference](./lal/main.go)

<h2 id="sky-java"> JT1078 sky-java </h2>

1. Start the service
2. Use RTVS Terminal Emulator to connect to the service
3. Call sky-java's JT1078 HTTP interface to send a request (need to pull the stream within 10 seconds by default)
- [sky-java Official Address](https://gitee.com/hui_hui_zhou/open-source-repository)
- [sky-java HTTP Documentation](http://222.244.144.181:9991/doc.html)
- [Code Reference](./sky-java/main.go)
