# CloudflareWarpSpeedTest

## 简介
Cloudflare Warp 速度测试是一个命令行工具，用于测试 Cloudflare Warp IP 地址的延迟和速度，并获取关于最低延迟和可用端口的信息。它提供了各种选项，以自定义测试参数并根据特定条件筛选结果。

Inspired by [CloudflareSpeedTest](https://github.com/XIU2/CloudflareSpeedTest)

## Download

  + [Releases](https://github.com/peanut996/CloudflareWarpSpeedTest/releases)
  + `go install github.com/peanut996/CloudflareWarpSpeedTest@latest`


## 用法

要使用 Cloudflare Warp Speed Test，您可以运行以下命令行选项

```bash
CloudflareWarpSpeedTest -n 200 -t 10 -tl 300 -tll 0 -tlr 0.2 -sl 5 -p 10 -f ip.txt -ip 1.1.1.1 -o result.csv -full
```

以下是主要可用选项的解释：

  + -n 200：指定延迟测试线程的数量。增加此值可以加快延迟测试过程，但不适合性能较低的设备，如路由器。默认值为 200，最大为 1000。
  + -t 10：设置对每个 IP 地址执行延迟测试的次数。默认值为 10 次。
  + -o result.csv：设置输出结果文件。默认文件为 "result.csv"。
  + -full：此标志表示应测试指定范围内的所有端口。
  + -h：显示帮助和用法说明。

更多使用说明请使用`-h`。

## 注意

请注意，调整测试参数可能会影响测试速度和结果。根据设备的性能和您希望应用的特定条件选择合适的设置至关重要。

**免责声明**： 本工具与 Cloudflare 无关，也未得到其认可。请负责任地使用并遵守其服务条款。

## License

此软件根据 [GPL v3 许可证](LICENSE) 发布