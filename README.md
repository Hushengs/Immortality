# Douyin Bot MVP

这是一个基于 Go 1.22+ 的抖音直播福袋执行端 MVP。当前版本以 Windows 环境为主，使用 `adb` 控制安卓手机，结合截图分析和 Tesseract OCR 识别福袋弹窗、倒计时和参与按钮，并通过状态机驱动主流程。

## 当前能力

- 识别可用 Android 设备并选择目标设备
- 执行截图、点击、滑动、返回等 ADB 操作
- 基于固定区域、像素特征和 OCR 检测福袋入口
- 解析福袋弹窗中的奖品、倒计时和按钮文案
- 根据白名单、黑名单和倒计时阈值决定是否参与
- 记录结构化运行日志和抽奖结果
- 在直播结束、长时间无福袋或出现验证提示时执行基础恢复

## 目录结构

```text
cmd/douyin-bot/           程序入口
internal/app/             领域模型、状态机、主运行器
internal/device/          ADB 客户端与设备管理
internal/capture/         截图和区域裁剪
internal/ocr/             Tesseract OCR 封装
internal/detectors/       福袋、布局、领奖、页面、验证检测
internal/actions/         导航、参与、领奖动作
internal/strategy/        参与和切房策略
internal/repository/      记录与设备仓储
internal/integration/     管理端预留接口
internal/utils/           日志、路径和公共工具
config.yaml               本地执行配置
```

## 依赖安装

### 1. 安装 Go

- 安装 Go 1.22 或更高版本
- 验证命令：

```bash
go version
```

### 2. 安装 ADB

- 安装 Android Platform Tools
- 确保 `adb` 可在命令行直接执行，或者把绝对路径写入 `config.yaml` 中的 `adb_path`
- 验证命令：

```bash
adb devices
```

### 3. 安装 Tesseract OCR

- Windows 可从 Tesseract 安装包完成安装
- 需要安装中文语言包 `chi_sim` 和英文语言包 `eng`
- 将 `tesseract.exe` 路径写入 `config.yaml` 中的 `tesseract_path`
- 验证命令：

```bash
tesseract --version
```

## 连接手机

1. 打开 Android 手机开发者选项与 USB 调试
2. 通过 USB 连接手机
3. 首次连接时在手机上确认调试授权
4. 使用以下命令验证设备：

```bash
adb devices
```

如果存在多台设备，可以在 `config.yaml` 中指定 `device_serial`。

## 配置说明

`config.yaml` 中重点关注：

- `safe_mode`: 默认开启。开启后允许识别中奖和记录截图，但不自动下单
- `auto_claim_reward`: 是否自动执行领奖动作
- `auto_place_order`: 是否允许自动下单，默认关闭
- `debug.dry_run`: 调试模式，不向真机发送点击、滑动、返回等动作
- `debug.static_screen_path`: 指定静态截图路径后，可不依赖真机直接跑检测与策略链路
- `debug.max_iterations`: 限制主循环次数，便于本地调试
- `adb.command_timeout_seconds` / `adb.retry_count`: ADB 超时与重试配置
- `strategy.prize_whitelist` / `strategy.prize_blacklist`: 奖品过滤规则
- `strategy.min_join_seconds` / `strategy.max_join_seconds`: 参与倒计时阈值
- `screen_regions`: 关键识别区域，建议按设备分辨率校准
- `thresholds.fudai_red_ratio`: 福袋红色像素阈值
- `ocr.preprocess_enabled`: 是否启用 OCR 预处理

## 运行方式

首次安装依赖：

```bash
go mod tidy
```

直接运行：

```bash
go run ./cmd/douyin-bot -config ./config.yaml
```

使用静态截图 dry-run 调试：

1. 在 `config.yaml` 中设置：
   - `debug.dry_run: true`
   - `debug.static_screen_path: D:/path/to/sample.png`
   - `debug.max_iterations: 1`
2. 运行：

```bash
go run ./cmd/douyin-bot -config ./config.yaml
```

这样会直接复制静态截图进入处理链路，不会向设备发送动作命令。

编译后运行：

```bash
go build -o douyin-bot.exe ./cmd/douyin-bot
./douyin-bot.exe -config ./config.yaml
```

## 安全模式

- `safe_mode=true`：记录中奖、保存截图、输出提醒，不执行自动下单
- `safe_mode=false` 且 `auto_claim_reward=true`：允许自动点击领奖
- `auto_place_order` 默认关闭，除非你明确理解风险并手动开启

## 建议调试方式

- 先在 `safe_mode=true` 下运行，观察 `data/logs/` 与 `data/records/`
- 通过保存的截图校准 `config.yaml` 中的区域坐标
- 如 OCR 识别不稳定，可先对静态截图单独调试
- 调试识别链路时，优先使用 `debug.dry_run + debug.static_screen_path + debug.max_iterations=1`

## 后续扩展

- 把本地配置源替换为远程 `ConfigProvider`
- 把本地 JSONL 记录替换为数据库或 HTTP API 上报
- 增加管理端 Web 界面；如果引入前端，默认采用 `layui`
