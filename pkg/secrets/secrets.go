package secrets


import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "syscall"
    
    "golang.org/x/term"
)

// getConfigPath returns the path where the encrypted config will be stored
func GetConfigPath() (string, error) {
    var configDir string
    
    switch runtime.GOOS {
    case "windows":
        configDir = os.Getenv("APPDATA")
    case "darwin":
        home, _ := os.UserHomeDir()
        configDir = filepath.Join(home, "Library", "Application Support")
    default: // linux and others
        configDir = os.Getenv("XDG_CONFIG_HOME")
        if configDir == "" {
            home, _ := os.UserHomeDir()
            configDir = filepath.Join(home, ".config")
        }
    }
    
    appConfigDir := filepath.Join(configDir, "myapp")
    if err := os.MkdirAll(appConfigDir, 0700); err != nil {
        return "", err
    }
    
    return filepath.Join(appConfigDir, "config.enc"), nil
}

// getMachineID generates a machine-specific key for encryption
// This prevents the config file from being copied to another machine
func GetMachineID() string {
    // Combine multiple system identifiers
    hostname, _ := os.Hostname()
    
    // You could add more identifiers here for stronger binding
    // For example: MAC address, system UUID, etc.
    
    hash := sha256.Sum256([]byte(hostname))
    return base64.StdEncoding.EncodeToString(hash[:])
}

// deriveKey creates an encryption key from the machine ID
func DeriveKey(machineID string) []byte {
    hash := sha256.Sum256([]byte(machineID))
    return hash[:]
}

// encrypt encrypts plaintext using AES-GCM
func Encrypt(plaintext, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts ciphertext using AES-GCM
func Decrypt(encodedCiphertext string, key []byte) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
    if err != nil {
        return nil, err
    }
    
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

// readAPIKey reads the API key from encrypted storage
func ReadAPIKey() (string, error) {
    configPath, err := GetConfigPath()
    if err != nil {
        return "", err
    }
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        return "", err
    }
    
    machineID := GetMachineID()
    key := DeriveKey(machineID)
    
    decrypted, err := Decrypt(string(data), key)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt (file may be from another machine)")
    }
    
    return string(decrypted), nil
}

// saveAPIKey encrypts and saves the API key
func saveAPIKey(apiKey string) error {
    configPath, err := GetConfigPath()
    if err != nil {
        return err
    }
    
    machineID := GetMachineID()
    key := DeriveKey(machineID)
    
    encrypted, err := Encrypt([]byte(apiKey), key)
    if err != nil {
        return err
    }
    
    return os.WriteFile(configPath, []byte(encrypted), 0600)
}

// promptForAPIKey securely prompts the user for their API key
func promptForAPIKey() (string, error) {
    fmt.Print("Enter your API key: ")
    
    // Read password without echoing to terminal
    bytePassword, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        return "", err
    }
    fmt.Println() // Print newline after hidden input
    
    apiKey := strings.TrimSpace(string(bytePassword))
    if apiKey == "" {
        return "", fmt.Errorf("API key cannot be empty")
    }
    
    return apiKey, nil
}

// getAPIKey retrieves the API key, prompting if not found
func GetAPIKey() (string, error) {
    apiKey, err := ReadAPIKey()
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("No API key found. Let's set one up.")
            apiKey, err = promptForAPIKey()
            if err != nil {
                return "", err
            }
            
            if err := saveAPIKey(apiKey); err != nil {
                return "", fmt.Errorf("failed to save API key: %w", err)
            }
            
            fmt.Println("✓ API key saved securely!")
            return apiKey, nil
        }
        return "", err
    }
    
    return apiKey, nil
}

// resetAPIKey allows the user to reset their API key
func ResetAPIKey() error {
    fmt.Println("Resetting API key...")
    
    apiKey, err := promptForAPIKey()
    if err != nil {
        return err
    }
    
    if err := saveAPIKey(apiKey); err != nil {
        return fmt.Errorf("failed to save API key: %w", err)
    }
    
    fmt.Println("✓ API key updated successfully!")
    return nil
}

func main() {
    // Check for --reset-key flag
    if len(os.Args) > 1 && os.Args[1] == "--reset-key" {
        if err := ResetAPIKey(); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
        return
    }
    
    // Get the API key (prompts if not found)
    apiKey, err := GetAPIKey()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    
    // Use the API key in your application
    fmt.Printf("Using API key: %s...\n", apiKey[:min(8, len(apiKey))])
    
    // Your actual application logic here
    fmt.Println("Running application...")
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}