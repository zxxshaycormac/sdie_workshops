# 课前环境准备

工作坊开始前，请确保你的机器满足以下所有要求。
以下步骤覆盖 macOS、Windows (WSL2) 和 Linux。

## 1. 前置条件

| 工具 | 最低版本 | 用途 |
|---|---|---|
| Git | 2.x | 版本控制 |
| Go | 1.22+ | 后端编译 |
| Node.js | 18+ | 前端构建 |
| Make | 任意 | 构建编排 |

### 按平台安装

**macOS (Homebrew)：**

```bash
xcode-select --install    # Git + Make
brew install go node
```

**Windows (WSL2 — Ubuntu)：**

打开 WSL2 终端（推荐 Ubuntu），然后执行：

```bash
sudo apt update
sudo apt install -y git make gcc
```

从 https://go.dev/dl/ 安装 Go — 下载 `linux-amd64` 压缩包：

```bash
curl -LO https://go.dev/dl/go1.22.10.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.10.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

通过 NodeSource 安装 Node.js：

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
```

> **WSL2 提示：** 如果 `make build` 因 CGO 错误失败，请安装 C 编译器：
> `sudo apt install -y gcc`

**Linux (Debian/Ubuntu)：**

与 WSL2 步骤相同（跳过 WSL2 本身的设置）。

**Linux (Fedora/RHEL)：**

```bash
sudo dnf install -y git make gcc
```

Go 和 Node.js — 参照上面的压缩包/NodeSource 说明安装，或使用：

```bash
sudo dnf install -y golang nodejs
```

### 验证（所有平台）

```bash
go version      # 预期：go1.22.x 或更高
node --version  # 预期：v18.x.x 或更高
git --version
make --version
```

## 2. 克隆与构建

```bash
# 克隆仓库（如果尚未克隆）
git clone $WORKSHOP
cd gitea

# 使用 SQLite 支持构建
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

首次构建需要 2-5 分钟（前端 + 后端），后续构建更快。

验证二进制文件：

```bash
./gitea --version
# 预期：Gitea version 1.22.x built with ... sqlite ...
```

## 3. 启动开发服务器

```bash
GITEA_RUN_MODE=dev ./gitea web
```

在浏览器中打开 **http://localhost:3000**。

- 首次访问会显示**安装页面** — 保持默认值（SQLite3），设置管理员用户名/密码，点击 "Install Gitea"。
- 此操作会创建 `custom/conf/app.ini` 和 SQLite 数据库。
- 使用 `Ctrl+C` 停止服务器。

### macOS / Linux

上述命令在任何终端中直接运行即可。

### Windows (WSL2)

在 WSL2 终端中运行相同命令。然后在 Windows 浏览器中访问 **http://localhost:3000** — WSL2 会自动转发 localhost 端口。

如果 localhost 转发未启用，查找 WSL2 IP：

```bash
hostname -I | awk '{print $1}'
# 在 Windows 浏览器中打开 http://<该IP>:3000
```

## 4. 验证测试

```bash
# 后端单元测试（应通过）
make test-backend

# 单个测试（示例）
go test -run TestBufWriter ./modules/git/pkg/

# SQLite 集成测试
make test-sqlite
```

集成测试耗时较长。至少确认 `make test-backend` 通过。

## 5. 安装 Claude Code

工作坊使用 Claude Code 作为 AI 编程工具链。

### macOS / Linux / WSL2

```bash
npm install -g @anthropic-ai/claude-code
```

### Windows（原生）

Claude Code 在 WSL2 内运行 — 请在 WSL2 中安装，而非 PowerShell。

验证：

```bash
claude --version
```

## 6. 推荐编辑器配置

- **VS Code** + [Go 扩展](https://marketplace.visualstudio.com/items?itemName=golang.Go)，用于代码跳转和重构
- **Claude Code VS Code 扩展**（可选），用于编辑器内 AI 辅助
- **WSL2 用户：** 通过 VS Code 的 "Remote - WSL" 扩展打开仓库文件夹，以获得完整 IDE 支持

## 常见问题

### `make build` 因 Go 版本错误失败

确保 Go >= 1.22：

```bash
go version
# macOS:  brew upgrade go
# Linux:  从 go.dev/dl 重新安装或通过包管理器升级
```

### `make build` 在前端阶段失败（webpack）

确保 Node >= 18：

```bash
node --version
# macOS:  brew upgrade node
# Linux:  sudo apt install -y nodejs（通过 NodeSource）
```

### CGO / C 编译器错误（Linux / WSL2）

SQLite 需要 CGO，CGO 需要 C 编译器：

```bash
# Debian/Ubuntu
sudo apt install -y gcc

# Fedora/RHEL
sudo dnf install -y gcc
```

### 3000 端口已被占用

```bash
# macOS / Linux / WSL2
lsof -i :3000       # 查找进程
kill <PID>           # 终止进程
```

或使用其他端口：

```bash
./gitea web -p 3001
```

### SQLite 构建标签缺失

二进制文件必须包含 `sqlite` 构建标签。验证：

```bash
./gitea --version
# 检查输出中是否包含 "sqlite"
```

如果缺失，重新构建：

```bash
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

### WSL2：`make build` 因 "cannot find package" 失败

确保 Go 和 Node 安装在 **WSL2 内部**，而非 Windows 中。验证：

```bash
which go    # 应为 /usr/local/go/bin/go 或类似的 Linux 路径
which node  # 应为 /usr/bin/node 或类似的 Linux 路径
```

如果指向 Windows 路径（如 `/mnt/c/...`），说明你使用了错误的 shell。

## 快速参考

```bash
# 构建
TAGS="bindata sqlite sqlite_unlock_notify" make build

# 运行
GITEA_RUN_MODE=dev ./gitea web

# 测试
make test-backend
go test -run TestFunctionName ./path/to/package

# 代码检查
make lint

# 热重载（替代手动重新构建）
make watch
```
