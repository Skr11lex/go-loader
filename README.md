1、建议使用garble编译
macos&linux用户：GOOS=windows GOARCH=amd64 garble -tiny -literals -seed=random build -ldflags="-s -w -buildid=" -trimpath
windows用户：garble -tiny -literals -seed=random build -ldflags="-s -w -buildid=" -trimpath
2、如果360提示QVM202，可以自行寻找QVM bypass工具进行混淆
