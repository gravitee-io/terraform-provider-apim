import re
import json
import sys
import subprocess
import shutil

def remove_nulls_and_empty_lists(data):
    """Recursively removes keys with null values or empty lists from a Python object."""
    if isinstance(data, dict):
        return {
            k: remove_nulls_and_empty_lists(v)
            for k, v in data.items()
            if v is not None and not (isinstance(v, list) and len(v) == 0)
        }
    elif isinstance(data, list):
        # Filter out nulls and empty lists from the list itself
        new_list = []
        for item in data:
            if item is None:
                continue
            if isinstance(item, list) and len(item) == 0:
                continue
            new_list.append(remove_nulls_and_empty_lists(item))
        return new_list
    return data

def format_hcl_val(obj, indent=2):
    """Formats a Python object into HCL-style (keys with =)."""
    dumped = json.dumps(obj, indent=indent)
    return re.sub(r'(".*?")\s*:', r'\1 =', dumped)

def transform_json_string(match):
    """Try to convert a matched string into jsonencode() blocks."""
    original_content = match.group(1)
    try:
        unescaped = original_content.encode().decode('unicode_escape')
        data = json.loads(unescaped)

        if not isinstance(data, (dict, list)) or not data:
            return f'"{original_content}"'

        clean_data = remove_nulls_and_empty_lists(data)
        hcl_body = format_hcl_val(clean_data)
        return f"jsonencode({hcl_body})"

    except (json.JSONDecodeError, UnicodeDecodeError, ValueError):
        return f'"{original_content}"'

def process_file_logic(content):
    # 1. First Pass: Remove multi-line and single-line HCL nulls/empty arrays
    # This regex handles:
    # key = null
    # key = []
    # key = [
    # ]
    # The [ \t]* matches spaces/tabs, and (?:\r?\n)? optionally matches the newline
    hcl_cleanup_pattern = r'^\s*[a-zA-Z0-9_-]+\s*=\s*(?:null|\[\s*\])\s*,?\s*$'

    # We apply this line by line for safety, but we'll pre-process the content
    # to collapse empty multi-line arrays into single lines first.
    # After 'terraform fmt', an empty multi-line array looks like: key = [\n]\n
    content = re.sub(r'([a-zA-Z0-9_-]+\s*=\s*\[)\s*\n\s*(\])', r'\1\2', content)

    lines = content.splitlines()
    processed_lines = []
    string_pattern = r'"((?:[^"\\]|\\.)*)"'

    for line in lines:
        stripped = line.lstrip()

        # Skip comments
        if stripped.startswith('#') or stripped.startswith('//'):
            processed_lines.append(line)
            continue

        # Remove native HCL nulls and (now collapsed) empty arrays
        if re.match(hcl_cleanup_pattern, line):
            continue

        # Process JSON strings
        new_line = re.sub(string_pattern, transform_json_string, line)
        processed_lines.append(new_line)

    return "\n".join(processed_lines)

def run_fmt(tf_path, target_file):
    if tf_path:
        print(f"--- Running 'terraform fmt' on {target_file} ---")
        subprocess.run([tf_path, "fmt", target_file], check=True)

def main():
    if len(sys.argv) < 3:
        print("Usage: python3 cleanup.py <input.tf> <output.tf>")
        sys.exit(1)

    input_path = sys.argv[1]
    output_path = sys.argv[2]
    tf_path = shutil.which("terraform")

    try:
        shutil.copy2(input_path, output_path)

        # Format 1: Ensure predictable structure
        run_fmt(tf_path, output_path)

        with open(output_path, 'r', encoding='utf-8') as f:
            formatted_content = f.read()

        final_content = process_file_logic(formatted_content)

        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(final_content)

        # Format 2: Final cleanup
        run_fmt(tf_path, output_path)
        print("✅ Refactor complete: Copied, formatted, cleaned, and re-formatted.")

    except Exception as e:
        print(f"❌ Error: {e}")

if __name__ == "__main__":
    main()