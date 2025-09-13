#!/usr/bin/env python3
"""
Script to update all tokens.go files to use symbol.For() pattern with var declarations.

This script will:
1. Find all tokens.go files
2. Add the symbol import if missing
3. Change const (...) to var (...)
4. Convert string literals to symbol.For() calls
5. Preserve comments and formatting
"""

import os
import re
import glob
from typing import List, Tuple


def find_tokens_files(root_dir: str) -> List[str]:
    """Find all tokens.go files in the directory tree."""
    pattern = os.path.join(root_dir, "**/tokens.go")
    return glob.glob(pattern, recursive=True)


def has_symbol_import(content: str) -> bool:
    """Check if the file already has the symbol import."""
    return 'import "govel/packages/support/src/symbol"' in content or \
           '"govel/packages/support/src/symbol"' in content


def add_symbol_import(content: str) -> str:
    """Add the symbol import to the file."""
    lines = content.split('\n')
    
    # Find the package declaration
    package_line = -1
    for i, line in enumerate(lines):
        if line.startswith('package '):
            package_line = i
            break
    
    if package_line == -1:
        raise ValueError("No package declaration found")
    
    # Check if there's already an import block
    import_start = -1
    import_end = -1
    has_imports = False
    
    for i in range(package_line + 1, len(lines)):
        line = lines[i].strip()
        if line == '':
            continue
        if line.startswith('import'):
            has_imports = True
            if line.startswith('import ('):
                # Multi-line import
                import_start = i
                for j in range(i + 1, len(lines)):
                    if lines[j].strip() == ')':
                        import_end = j
                        break
                break
            elif line.startswith('import '):
                # Single import - convert to multi-line
                import_line = line
                lines[i] = 'import ('
                lines.insert(i + 1, '\t' + import_line[7:])  # Remove 'import '
                lines.insert(i + 2, ')')
                import_start = i
                import_end = i + 2
                break
        else:
            break
    
    if not has_imports:
        # No imports exist, create new import block
        lines.insert(package_line + 1, '')
        lines.insert(package_line + 2, 'import "govel/packages/support/src/symbol"')
        lines.insert(package_line + 3, '')
    else:
        # Add to existing import block
        if import_start != -1 and import_end != -1:
            # Add the symbol import
            lines.insert(import_end, '\t"govel/packages/support/src/symbol"')
        else:
            # Single import case already handled above
            pass
    
    return '\n'.join(lines)


def convert_const_to_var(content: str) -> str:
    """Convert const (...) to var (...)."""
    # Replace const ( with var (
    content = re.sub(r'\bconst\s*\(', 'var (', content)
    return content


def convert_string_to_symbol(content: str) -> str:
    """Convert string literals like 'govel.xxx' to symbol.For('govel.xxx')."""
    # Pattern to match lines like: TOKEN_NAME = "govel.something"
    # But not lines that already have symbol.For()
    pattern = r'(\s*)([A-Z_]+)\s*=\s*"(govel\.[^"]+)"'
    
    def replace_func(match):
        indent = match.group(1)
        token_name = match.group(2)
        token_value = match.group(3)
        return f'{indent}{token_name} = symbol.For("{token_value}")'
    
    # Only replace if the line doesn't already contain symbol.For
    lines = content.split('\n')
    new_lines = []
    
    for line in lines:
        if 'symbol.For(' not in line and re.search(r'[A-Z_]+\s*=\s*"govel\.[^"]+"', line):
            new_line = re.sub(pattern, replace_func, line)
            new_lines.append(new_line)
        else:
            new_lines.append(line)
    
    return '\n'.join(new_lines)


def update_tokens_file(file_path: str) -> bool:
    """Update a single tokens.go file. Returns True if file was modified."""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            original_content = f.read()
        
        content = original_content
        
        # Step 1: Add symbol import if missing
        if not has_symbol_import(content):
            content = add_symbol_import(content)
        
        # Step 2: Convert const to var
        content = convert_const_to_var(content)
        
        # Step 3: Convert string literals to symbol.For() calls
        content = convert_string_to_symbol(content)
        
        # Check if content changed
        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            return True
        
        return False
    
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False


def main():
    """Main function to update all tokens.go files."""
    root_dir = "/Users/akouta/Projects/govel"
    
    # Find all tokens.go files
    tokens_files = find_tokens_files(root_dir)
    
    if not tokens_files:
        print("No tokens.go files found!")
        return
    
    print(f"Found {len(tokens_files)} tokens.go files:")
    for file_path in tokens_files:
        print(f"  {file_path}")
    
    print("\nUpdating files...")
    
    updated_count = 0
    for file_path in tokens_files:
        print(f"Processing: {file_path}")
        if update_tokens_file(file_path):
            print(f"  ✅ Updated")
            updated_count += 1
        else:
            print(f"  ⏭️  No changes needed")
    
    print(f"\nCompleted! Updated {updated_count} out of {len(tokens_files)} files.")


if __name__ == "__main__":
    main()