#!/bin/bash

# Script to fix all lowercase type names in GoVel types package
echo "🔧 Fixing type names in GoVel types package..."

# Find all .go files with lowercase type declarations and fix them
find /Users/akouta/Projects/govel/packages/types/src/types -name "*.go" -type f | while read -r file; do
    # Check if file contains lowercase type declarations
    if grep -q "type [a-z].*Type interface" "$file"; then
        echo "📝 Processing: $file"
        
        # Create backup
        cp "$file" "${file}.backup"
        
        # Fix common lowercase type patterns using sed
        sed -i '' \
            -e 's/type cipherType interface/type CipherType interface/g' \
            -e 's/type payloadType interface/type PayloadType interface/g' \
            -e 's/type pipelinecallbackType interface/type PipelineCallbackType interface/g' \
            -e 's/type hasherType interface/type HasherType interface/g' \
            -e 's/type instancecreatorType interface/type InstanceCreatorType interface/g' \
            -e 's/type drivercreatorType interface/type DriverCreatorType interface/g' \
            "$file"
        
        echo "  ✅ Fixed type names in $file"
    fi
done

# Also check for other possible lowercase type patterns (not just ending with Type)
echo ""
echo "🔍 Checking for other lowercase type patterns..."

# Check for any type definitions starting with lowercase
find /Users/akouta/Projects/govel/packages/types/src/types -name "*.go" -type f | while read -r file; do
    if grep -q "type [a-z].*Interface interface\|type [a-z].*Struct struct\|type [a-z][a-zA-Z]*[[:space:]]*interface[[:space:]]*{" "$file"; then
        echo "📝 Processing additional patterns in: $file"
        
        # Create backup if not already created
        if [[ ! -f "${file}.backup" ]]; then
            cp "$file" "${file}.backup"
        fi
        
        # Fix additional patterns that might exist
        sed -i '' \
            -e 's/type maintenancemodeType interface/type MaintenanceModeType interface/g' \
            -e 's/type hookcallbackType interface/type HookCallbackType interface/g' \
            -e 's/type providercallbackType interface/type ProviderCallbackType interface/g' \
            -e 's/type shutdowncallbackType interface/type ShutdownCallbackType interface/g' \
            -e 's/type serviceidentifierType interface/type ServiceIdentifierType interface/g' \
            "$file"
        
        echo "  ✅ Fixed additional patterns in $file"
    fi
done

echo ""
echo "🔍 Verifying changes..."

# Check if any lowercase type declarations remain
remaining_files=$(find /Users/akouta/Projects/govel/packages/types/src/types -name "*.go" -exec grep -l "type [a-z].*Type interface\|type [a-z].*Interface interface\|type [a-z].*Struct struct" {} \; 2>/dev/null | wc -l)

if [[ $remaining_files -eq 0 ]]; then
    echo "✅ All type names successfully capitalized!"
else
    echo "⚠️  $remaining_files files still have lowercase type names:"
    find /Users/akouta/Projects/govel/packages/types/src/types -name "*.go" -exec grep -l "type [a-z].*Type interface\|type [a-z].*Interface interface\|type [a-z].*Struct struct" {} \; 2>/dev/null
    echo ""
    echo "Remaining lowercase patterns:"
    find /Users/akouta/Projects/govel/packages/types/src/types -name "*.go" -exec grep "type [a-z].*Type interface\|type [a-z].*Interface interface\|type [a-z].*Struct struct" {} \; 2>/dev/null
fi

echo ""
echo "🧪 Testing compilation..."
cd /Users/akouta/Projects/govel/packages/types/src/types

if go build ./... 2>/dev/null; then
    echo "✅ All types compile successfully!"
else
    echo "⚠️  Some compilation issues may exist. Checking specific errors..."
    go build ./... 2>&1 | head -20
fi

echo ""
echo "📊 Summary:"
echo "💾 Backups created with .backup extension"
echo "🔧 To remove backups: find /Users/akouta/Projects/govel/packages/types/src/types -name '*.backup' -delete"
echo "🎉 Type name fix completed!"