# PgInspector

```text
postgresql-inspector/
├── agent/                       # （可选）作为服务运行的客户端
│   └── main.go                  # 入口文件
├── pkg/                             # 核心库（可被外部导入）
│   ├── pg-inspect/                      # 核心功能包
│   │   ├── client.go                # 核心数据库连接及巡检逻辑
│   │   ├── metrics.go               # Prometheus 指标采集逻辑
│   │   ├── alert.go                 # 飞书通知或其他警报机制
│   │   ├── report.go                # 生成巡检报告
│   │   └── test/                    # 单元测试
│   ├── config/                      # 配置解析与管理
│   │   ├── config.go                # 配置加载、验证与解析
│   │   └── defaults.go              # 默认配置项
│   ├── dashboard/                   # Grafana 仪表盘相关
│   │   ├── dashboard.go             # 管理 Grafana 仪表盘导出与导入
│   │   └── template.json            # 默认仪表盘模板
│   └── utils/                       # 通用工具
│       ├── logger.go                # 日志工具
│       └── http.go                  # HTTP 请求工具
├── configs/                         # 配置文件示例
│   ├── config.example.yaml          # 配置文件示例（用户可以复制并修改）
│   └── config.yaml                  # 默认配置文件
├── dashboards/                      # Grafana 仪表盘 JSON 模板
│   └── postgresql_inspector.json    # 示例仪表盘模板
├── scripts/                         # 部署或辅助脚本
│   ├── build.sh                     # 构建脚本
│   └── deploy.sh                    # 部署脚本
├── Makefile                         # 构建和测试自动化
├── go.mod                           # Go 模块定义
└── README.md                        # 项目说明文档

```