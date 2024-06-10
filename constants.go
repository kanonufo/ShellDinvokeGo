package main

// Constantes para la asignación de memoria y protección de página
const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
	MEM_RELEASE            = 0x8000
)
