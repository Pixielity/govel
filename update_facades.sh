#!/bin/bash

# Script to batch update facade files to use interface imports
# This script updates facade files to use the new interface import pattern

# Define the facade updates in associative arrays
declare -A facade_updates=(
    ["hash_facade.go"]="hashInterfaces:hashing:hasher_interface:HashInterface:HASH_TOKEN"
    ["http_facade.go"]="httpInterfaces:http:http_interface:HttpInterface:HTTP_TOKEN"
    ["lang_facade.go"]="langInterfaces:lang:language_interface:LanguageInterface:LANG_TOKEN"
    ["mail_facade.go"]="mailInterfaces:mail:mail_interface:MailInterface:MAIL_TOKEN"
    ["orm_facade.go"]="ormInterfaces:orm:orm_interface:OrmInterface:ORM_TOKEN"
    ["view_facade.go"]="viewInterfaces:view:view_interface:ViewInterface:VIEW_TOKEN"
)

# Base directory
BASE_DIR="/Users/akouta/Projects/govel/packages/support/src/facades"

# Function to update a facade file
update_facade() {
    local file=$1
    local config=$2
    local file_path="$BASE_DIR/$file"
    
    echo "Updating $file..."
    
    # Parse the configuration
    IFS=':' read -r alias_name interface_package interface_file interface_name token_name <<< "$config"
    
    # Check if file exists
    if [[ ! -f "$file_path" ]]; then
        echo "Warning: $file_path does not exist"
        return 1
    fi
    
    # Create backup
    cp "$file_path" "${file_path}.backup"
    
    # Extract the service name from the facade file (e.g., http_facade.go -> http)
    service_name="${file%_facade.go}"
    
    # Update import statement
    sed -i.bak "s/${service_name}Tokens \"govel\/packages\/interfaces\/${interface_package}\"/${alias_name} \"govel\/packages\/interfaces\/${interface_package}\"/g" "$file_path"
    
    # Update function return types and token references
    # This is a simplified pattern - may need manual adjustment for complex cases
    sed -i.bak "s/func ${service_name^}() interface{}/func ${service_name^}() ${alias_name}.${interface_name}/g" "$file_path"
    sed -i.bak "s/func ${service_name^}WithError() (interface{}, error)/func ${service_name^}WithError() (${alias_name}.${interface_name}, error)/g" "$file_path"
    
    # Update token references
    sed -i.bak "s/${service_name}Tokens\..*_TOKEN/${alias_name}.${token_name}/g" "$file_path"
    
    # Update facade.Resolve calls
    sed -i.bak "s/facade\.Resolve\[interface{}\]/facade.Resolve[${alias_name}.${interface_name}]/g" "$file_path"
    sed -i.bak "s/facade\.TryResolve\[interface{}\]/facade.TryResolve[${alias_name}.${interface_name}]/g" "$file_path"
    
    # Clean up backup files
    rm "${file_path}.bak"
    
    echo "Updated $file successfully"
    return 0
}

# Process each facade file
for file in "${!facade_updates[@]}"; do
    update_facade "$file" "${facade_updates[$file]}"
done

echo "Facade update script completed!"