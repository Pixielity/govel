# Encryption Service Provider Integration Guide

This guide explains how to integrate the Encryption Service Provider into a GoVel application.

## Overview

The Encryption Service Provider registers all encryption services with the application container, including:

- **Encryption Manager**: Main service for managing multiple encryption drivers
- **Encryption Driver**: Default encrypter driver instance
- **Encrypter Interface**: Contract for encryption operations
- **Factory Interface**: Contract for creating encrypter instances

## Quick Start

### 1. Register the Service Provider

Add the encryption service provider to your application's service providers during bootstrap:

```go
package main

import (
    "context"
    "log"
    
    "govel/packages/application"
    encryptionProviders "govel/packages/encryption/src/providers"
    configProviders "govel/packages/config/providers"
)

func main() {
    // Create application instance
    app := application.NewApplication()
    
    // Register core service providers
    configProvider := configProviders.NewConfigServiceProvider()
    if err := app.RegisterProvider(configProvider); err != nil {
        log.Fatal("Failed to register config provider:", err)
    }
    
    // Register encryption service provider
    encryptionProvider := encryptionProviders.NewEncryptionServiceProvider()
    if err := app.RegisterProvider(encryptionProvider); err != nil {
        log.Fatal("Failed to register encryption provider:", err)
    }
    
    // Boot all providers
    ctx := context.Background()
    if err := app.BootProviders(ctx); err != nil {
        log.Fatal("Failed to boot providers:", err)
    }
    
    // Application is now ready to use encryption services
}
```

### 2. Configuration Requirements

The encryption service provider requires the following configuration keys:

```json
{
  "app": {
    "key": "your-32-character-encryption-key-here",
    "cipher": "AES-256-CBC"
  },
  "encryption": {
    "default": "default"
  }
}
```

#### Configuration Keys

- `app.key`: The encryption key (must be 16 or 32 bytes for AES-128/256)
- `app.cipher`: Default cipher algorithm (defaults to "AES-256-CBC")
- `encryption.default`: Default driver name (defaults to "default")

### 3. Using Encryption Services

Once registered, you can access encryption services from the container:

```go
// Get the main encryption manager
encryptionService, err := app.Make(encryptionInterfaces.ENCRYPTION_MANAGER_TOKEN)
if err != nil {
    log.Fatal("Failed to resolve encryption service:", err)
}

// Cast to encryption interface
encrypter := encryptionService.(encryptionInterfaces.EncrypterInterface)

// Encrypt data
encrypted, err := encrypter.EncryptString("sensitive data")
if err != nil {
    log.Fatal("Encryption failed:", err)
}

// Decrypt data
decrypted, err := encrypter.DecryptString(encrypted)
if err != nil {
    log.Fatal("Decryption failed:", err)
}
```

### 4. Available Service Tokens

The encryption service provider registers the following services:

```go
// Import the encryption interfaces
import encryptionInterfaces "govel/packages/interfaces/encryption"

// Available service tokens
encryptionInterfaces.ENCRYPTION_MANAGER_TOKEN          // Main encryption manager
encryptionInterfaces.ENCRYPTION_DRIVER_TOKEN           // Default encryption driver  
encryptionInterfaces.ENCRYPTION_CONTRACT_TOKEN         // Encrypter interface contract
encryptionInterfaces.ENCRYPTION_FACTORY_CONTRACT_TOKEN // Factory interface contract
```

## Advanced Usage

### Multiple Encryption Drivers

You can configure and use different encryption drivers:

```go
// Get encryption manager
manager := encryptionService.(*encryption.EncryptionManager)

// Use specific cipher drivers
cbcEncrypter := manager.Driver("AES-256-CBC")
gcmEncrypter := manager.Driver("AES-256-GCM")
ctrEncrypter := manager.Driver("AES-256-CTR")
```

### Supported Ciphers

The encryption package supports the following cipher algorithms:

- **AES-128-CBC**: AES-128 with CBC mode and HMAC-SHA256 MAC
- **AES-256-CBC**: AES-256 with CBC mode and HMAC-SHA256 MAC (Laravel default)
- **AES-128-GCM**: AES-128 with GCM mode (authenticated encryption)
- **AES-256-GCM**: AES-256 with GCM mode (authenticated encryption)
- **AES-128-CTR**: AES-128 with CTR mode and HMAC-SHA256 MAC
- **AES-256-CTR**: AES-256 with CTR mode and HMAC-SHA256 MAC

### Driver Configuration

You can configure specific drivers in your configuration:

```json
{
  "encryption": {
    "default": "default",
    "drivers": {
      "aes": {
        "cipher": "AES-256-GCM"
      }
    }
  }
}
```

## Service Provider Details

### Registration Process

1. **Singleton Registration**: The encryption manager is registered as a singleton service
2. **Driver Factory**: Default driver factory is registered for quick access
3. **Interface Contracts**: All encryption interfaces are bound to the manager
4. **Deferred Loading**: Services are loaded only when first requested

### Dependencies

The encryption service provider depends on:

- **Config Service**: For loading encryption configuration
- **Container Service**: For dependency injection support

### Service Priority

The encryption service provider has standard priority (0) and should be registered after core services like configuration.

## Troubleshooting

### Common Issues

1. **Invalid Key Length**: Ensure your `app.key` is exactly 16 or 32 bytes for AES-128/256
2. **Missing Config Service**: Register the config service provider before encryption
3. **Invalid Cipher**: Check that the configured cipher is supported

### Error Messages

- `"No application encryption key specified"`: Set the `app.key` configuration value
- `"Invalid key length for X: expected Y bytes, got Z bytes"`: Use correct key length for cipher
- `"Unsupported cipher: X"`: Use one of the supported cipher algorithms

### Debug Information

Enable debug logging to see provider loading and service registration details:

```go
app.GetLogger().SetLevel("debug")
```

## Laravel Compatibility

This implementation follows Laravel's encryption service patterns:

- Compatible payload format for cross-platform encryption/decryption
- Same cipher algorithms and modes
- Similar service provider registration pattern
- Equivalent configuration structure

## Security Considerations

1. **Key Management**: Store encryption keys securely (environment variables, key management services)
2. **Key Rotation**: Implement proper key rotation procedures
3. **Algorithm Choice**: Use AES-GCM for authenticated encryption when possible
4. **Configuration**: Validate configuration during application startup