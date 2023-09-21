package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: extract_pngs <executable_path>")
        return
    }

    exePath := os.Args[1]

    exeData, err := os.ReadFile(exePath)
    if err != nil {
        fmt.Printf("Error reading executable: %v\n", err)
        return
    }

    pngStartSig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
    pngEndSig := []byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

    var pngCount int

    for i := 0; i < len(exeData)-len(pngStartSig); i++ {
        if bytesEqual(exeData[i:i+len(pngStartSig)], pngStartSig) {
            // Found a PNG start signature, look for the end signature
            for j := i + len(pngStartSig); j < len(exeData)-len(pngEndSig); j++ {
                if bytesEqual(exeData[j:j+len(pngEndSig)], pngEndSig) {
                    // Found a PNG image
                    pngCount++
                    pngData := exeData[i : j+len(pngEndSig)]

                    // Save the PNG image as a file
                    fileName := fmt.Sprintf("extracted_png_%d.png", pngCount)
                    err := os.WriteFile(fileName, pngData, 0644)
                    if err != nil {
                        fmt.Printf("Error saving PNG file %s: %v\n", fileName, err)
                    } else {
                        fmt.Printf("Saved %s\n", fileName)
                    }

                    // Move to the next position after the end signature
                    i = j + len(pngEndSig) - 1
                    break
                }
            }
        }
    }
}

func bytesEqual(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}
