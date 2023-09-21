package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // Check for the command-line argument
    if len(os.Args) != 2 {
        fmt.Println("Usage: extract_pngs <executable_path>")
        return
    }

    // Get the executable path from the command-line argument
    exePath := os.Args[1]

    // Read the entire executable file
    exeData, err := os.ReadFile(exePath)
    if err != nil {
        fmt.Printf("Error reading executable: %v\n", err)
        return
    }

    // Define PNG signature bytes
    pngStartSig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
    pngEndSig := []byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}

    // Initialize a counter for extracted PNGs
    var pngCount int

    // Create a "resources" folder if it doesn't exist
    resourcesFolder := "resources"
    if _, err := os.Stat(resourcesFolder); os.IsNotExist(err) {
        os.Mkdir(resourcesFolder, os.ModePerm)
    }

    // Search for PNG signatures in the executable
    for i := 0; i < len(exeData)-len(pngStartSig); i++ {
        if bytesEqual(exeData[i:i+len(pngStartSig)], pngStartSig) {
            // Found the start of a PNG
            for j := i + len(pngStartSig); j < len(exeData)-len(pngEndSig); j++ {
                if bytesEqual(exeData[j:j+len(pngEndSig)], pngEndSig) {
                    // Found the end of the PNG
                    pngCount++
                    pngData := exeData[i : j+len(pngEndSig)]

                    // Define the path to save the PNG image
                    pngFileName := filepath.Join(resourcesFolder, fmt.Sprintf("extracted_png_%d.png", pngCount))

                    // Save the PNG image
                    err := os.WriteFile(pngFileName, pngData, 0644)
                    if err != nil {
                        fmt.Printf("Error saving PNG file %s: %v\n", pngFileName, err)
                    } else {
                        fmt.Printf("Saved %s\n", pngFileName)
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
