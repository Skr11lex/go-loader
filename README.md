# 项目编译说明

## 🚀 编译建议
推荐使用 [garble](https://github.com/burrowers/garble) 进行混淆编译。

### 📦 MacOS & Linux 用户
GOOS=windows GOARCH=amd64 garble -tiny -literals -seed=random build \
-ldflags="-s -w -buildid=" -trimpath

### 📦 Windows 用户
garble -tiny -literals -seed=random build `
  -ldflags="-s -w -buildid=" -trimpath

⚠️ 关于 360 提示
如果 360 报告 QVM202，可以尝试自行寻找 QVM bypass 工具 进行额外混淆处理。
