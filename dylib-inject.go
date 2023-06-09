package main

import (
	"fmt"
	"os"
	"syscall"
)

// injectDylib injects a dylib into a target process with the given process ID.
func injectDylib(pid int, dylibPath string) error {
	// Open the target process.
	targetProcess, err := os.Open(fmt.Sprintf("/proc/%d/mem", pid))
	if err != nil {
		return err
	}
	defer targetProcess.Close()

	// Load the dylib file.
	dylibFile, err := os.Open(dylibPath)
	if err != nil {
		return err
	}
	defer dylibFile.Close()

	// Get the size of the dylib file.
	dylibStat, err := dylibFile.Stat()
	if err != nil {
		return err
	}
	dylibSize := int(dylibStat.Size())

	// Allocate memory within the target process.
	targetMemory, err := syscall.Mmap(int(targetProcess.Fd()), 0, dylibSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
	defer syscall.Munmap(targetMemory)

	// Read the dylib into the allocated memory.
	_, err = dylibFile.ReadAt(targetMemory, 0)
	if err != nil {
		return err
	}

	// Perform the dylib injection.
	err = syscall.Mprotect(targetMemory, syscall.PROT_READ|syscall.PROT_EXEC)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	pid := 1234                                // Replace with the target process ID.
	dylibPath := "/path/to/InjectedCode.dylib" // Replace with the actual path to your dylib.

	err := injectDylib(pid, dylibPath)
	if err != nil {
		fmt.Printf("Failed to inject dylib: %v\n", err)
		return
	}

	fmt.Println("Dylib injected successfully!")
}
