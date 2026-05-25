# 抖音福袋项目简版 Prompt

```text
请帮我生成一个 Golang 项目，用于通过 ADB 控制安卓手机，在抖音直播间自动识别并参与福袋抽奖。

一、项目目标
1. 连接安卓手机，识别可用设备。
2. 通过 adb 截图、点击、滑动、返回键控制手机。
3. 在抖音直播间内检测是否存在福袋入口。
4. 打开福袋详情后识别：
   - 奖品内容
   - 倒计时
   - 参与按钮状态
   - 是否需要任务/是否需要粉丝团/是否无法参与
5. 根据策略判断是否参与抽奖，必要时切换直播间。
6. 等待开奖后识别是否中奖。
7. 中奖后支持领奖流程，但默认安全模式下禁止自动下单。
8. 当直播关闭、长时间无福袋、页面跳转、人机验证出现时，能够自动恢复到可继续挂机的状态。

二、技术要求
1. 使用 Go 1.22+。
2. 使用 `adb` 与手机交互。
3. 使用 Go 的图像处理能力，必要时可配合 `disintegration/imaging` 之类的库处理截图。
4. 使用 `gosseract` + `Tesseract OCR` 识别文字，或封装调用本机 `tesseract.exe`。
5. 中文文案使用 `chi_sim`，数字倒计时可使用 `eng`。
6. 不要把所有逻辑写在一个文件中，必须模块化。
7. 必须提供 `go.mod`、`README.md`、`config.yaml`。

三、识别逻辑要求
1. 福袋入口优先用像素颜色 + 固定区域扫描判断。
2. 福袋详情弹窗需要先判断布局类型，因为任务数量不同会导致内容区域和倒计时区域位置不同。
3. 识别内容时要同时支持：
   - 奖品内容 OCR
   - 倒计时 OCR
   - 按钮文案 OCR
4. 倒计时识别后要转换为秒数。
5. 识别结果统一输出为结构化数据对象。

四、业务策略要求
1. 支持奖品白名单和黑名单关键字。
2. 支持最小参与倒计时和最大参与倒计时。
3. 如果奖品不符合要求，关闭弹窗并切换直播间。
4. 如果倒计时太短或太长，关闭弹窗并切换直播间。
5. 如果按钮显示“参与抽奖”或类似文案，则执行点击参与。
6. 如果已经参与成功，则关闭弹窗并等待开奖。
7. 如果条件不满足，则跳过当前福袋或切换直播间。

五、导航与恢复要求
1. 支持以下页面状态检测：
   - 关注页
   - 直播列表页
   - 直播间页面
   - 直播已结束页面
   - 红包弹窗页面
   - 人机验证页面
2. 提供恢复流程：
   - 返回到直播列表
   - 从关注页进入直播列表
   - 刷新直播列表
   - 进入第一个直播间
   - 直播关闭后重新切换
3. 如果长时间没有福袋，要自动回到列表重新进入直播间。

六、人机验证要求
1. 至少预留两类处理：
   - 滑块验证
   - 点击图片验证
2. 滑块验证支持根据截图估算滑动距离。
3. 点击图片验证如果无法自动处理，要保留人工接管接口。
4. 所有验证场景都要截图并记录日志。

七、架构预留要求
1. 这个项目后续还会增加一个管理系统，用于统一管理：
   - 抽奖记录
   - 直播间配置
   - 设备配置
   - 策略配置
2. 因此当前生成的 Go 项目必须预留“执行端”和“管理端”分离的架构边界，不要把配置、记录、设备状态全部写死在本地文件里。
3. 核心业务逻辑要围绕领域模型设计，至少包括：
   - Device
   - LiveRoom
   - LotteryTask
   - LotteryRecord
   - RuntimeConfig
   - StrategyConfig
4. 配置读取必须支持抽象层，至少预留两种来源：
   - 本地 `config.yaml`
   - 远程管理系统 API
5. 抽奖结果、运行日志、设备心跳、异常事件必须以结构化数据输出，方便后续直接落库或上报到管理系统。
6. 数据模型中要预留唯一标识字段，例如：
   - `device_id`
   - `room_id`
   - `task_id`
   - `record_id`
   - `run_id`
7. 请预留接口层或适配器层，例如：
   - `ConfigProvider`
   - `RecordRepository`
   - `DeviceRepository`
   - `LiveRoomRepository`
   - `ManagementClient`
8. 当前可以先实现本地版本，但代码结构必须允许后续无痛替换成：
   - HTTP API
   - gRPC
   - 数据库持久化

八、安全要求
1. 必须提供 `safe_mode` 配置，默认开启。
2. `safe_mode=true` 时：
   - 可以识别中奖
   - 可以保存中奖截图
   - 可以提醒用户
   - 不允许自动下单
3. 自动领奖和自动下单必须是显式配置项，默认关闭。

九、配置项要求
配置中至少包括：
1. 设备序列号
2. 手机高度分辨率
3. Y 轴偏移量
4. Tesseract 路径
5. 是否切换直播间
6. 最大等待分钟数
7. 奖品白名单和黑名单
8. 安全模式
9. 截图目录和日志目录

十、工程结构要求
请按模块化方式生成，建议结构如下：

douyin_bot/
  cmd/
    douyin-bot/
      main.go
  go.mod
  README.md
  config.yaml
  internal/
    app/
      runner.go
      state_machine.go
      models.go
    device/
      adb_client.go
      device_manager.go
    capture/
      screenshot_service.go
      image_cropper.go
    ocr/
      ocr_service.go
    detectors/
      fudai_detector.go
      popup_layout_detector.go
      reward_detector.go
      live_status_detector.go
      verification_detector.go
    actions/
      live_navigation.go
      lottery_actions.go
      reward_actions.go
    strategy/
      prize_filter.go
      room_switch_strategy.go
    repository/
      record_repository.go
      device_repository.go
      live_room_repository.go
    integration/
      management_client.go
    utils/
      logger.go
      paths.go

十一、主流程要求
请用状态机实现主流程，至少包含这些状态：
1. INIT
2. SELECT_DEVICE
3. CAPTURE_SCREEN
4. CHECK_FUDAI_ICON
5. OPEN_FUDAI_DETAIL
6. PARSE_POPUP_LAYOUT
7. OCR_CONTENT
8. DECIDE_JOIN
9. CLICK_JOIN
10. WAIT_RESULT
11. CHECK_RESULT
12. CLAIM_REWARD
13. SWITCH_ROOM
14. RECOVER_STATE
15. HANDLE_VERIFICATION

十二、输出要求
请直接生成：
1. 完整项目目录结构
2. 核心 Go 代码
3. go.mod
4. config.yaml 示例
5. README.md，说明：
   - 如何安装 Go 依赖
   - 如何安装并配置 Tesseract
   - 如何连接手机
   - 如何运行程序
   - 如何开启/关闭安全模式

十三、额外要求
1. 代码中关键步骤要有清晰日志。
2. 坐标和阈值不要全部硬编码，尽量做成配置或常量。
3. 识别结果尽量结构化，不要到处直接写分支判断。
4. 保留后续扩展空间，方便以后替换成更稳定的视觉识别方式。
5. 默认以 Windows 环境为主。
6. 本地执行端要能独立运行，但要天然支持后续被管理系统统一纳管。
```