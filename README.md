# 📖 项目编译说明

本项目推荐使用 [garble](https://github.com/burrowers/garble) 进行编译与混淆

## 🚀 编译方法

### MacOS & Linux 用户
```bash
GOOS=windows GOARCH=amd64 garble -tiny -literals -seed=random build \
  -ldflags="-s -w -buildid=" -trimpath
```

### Windows 用户
```bash
garble -tiny -literals -seed=random build -ldflags="-s -w -buildid=" -trimpath
```


## 🛠️ 使用说明
- 示例用法已写在代码中，可直接运行  
- 可根据需要修改：  
  - 去除多余的控制台输出  
  - 调整参数名称或默认值
  - 修改默认密钥  

## ⚠️ 常见问题
- 编译完成后建议在目标环境中进行测试，确保兼容性
- 360报告QVM，自行寻找QVMbypass工具对编译后的程序进行额外混淆
