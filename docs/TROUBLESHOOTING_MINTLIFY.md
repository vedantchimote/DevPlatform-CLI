# Mintlify Error 500 Troubleshooting

## Issue
Getting "Error 500 - Error loading page" when accessing Mintlify documentation at http://localhost:3000

## Possible Causes & Solutions

### 1. Browser Cache Issue
**Solution**: Hard refresh the browser
- **Windows/Linux**: Ctrl + F5 or Ctrl + Shift + R
- **Mac**: Cmd + Shift + R
- Or open in incognito/private mode

### 2. Mintlify Cache Issue
**Solution**: Clear Mintlify cache and restart
```bash
# Stop the dev server (Ctrl+C)

# Clear Mintlify cache (Windows)
Remove-Item -Path "$env:USERPROFILE\.mintlify" -Recurse -Force

# Clear Mintlify cache (Mac/Linux)
rm -rf ~/.mintlify

# Restart the server
cd docs
npx mintlify dev
```

### 3. Configuration Issue
**Solution**: Verify mint.json is valid
```bash
# Check if mint.json exists
ls docs/mint.json

# Validate JSON syntax
# Use an online JSON validator or:
cat docs/mint.json | python -m json.tool
```

### 4. File Path Issues
**Solution**: Verify all referenced files exist
```bash
# Check if all MDX files exist
cd docs
ls introduction.mdx quickstart.mdx installation.mdx
ls concepts/*.mdx
ls aws/*.mdx
ls azure/*.mdx
```

### 5. Port Conflict
**Solution**: Check if port 3000 is already in use
```bash
# Windows
netstat -ano | findstr :3000

# Mac/Linux
lsof -i :3000

# If port is in use, kill the process or use a different port
npx mintlify dev --port 3001
```

### 6. Node Modules Issue
**Solution**: Clear node modules cache
```bash
# Clear npm cache
npm cache clean --force

# Remove node_modules if it exists
rm -rf node_modules

# Restart server
cd docs
npx mintlify dev
```

### 7. MDX Syntax Errors
**Solution**: Check for syntax errors in MDX files
- Look for unclosed tags
- Check for invalid JSX syntax
- Verify frontmatter is properly formatted

Common issues:
- Missing closing tags: `</Card>`, `</Tab>`, `</Accordion>`
- Invalid frontmatter (must be at the very top)
- Mixing markdown and JSX incorrectly

### 8. Try Minimal Configuration
**Solution**: Test with a minimal mint.json

Create a backup and try this minimal config:
```json
{
  "name": "DevPlatform CLI",
  "logo": {
    "dark": "/logo/dark.svg",
    "light": "/logo/light.svg"
  },
  "navigation": [
    {
      "group": "Get Started",
      "pages": [
        "introduction"
      ]
    }
  ]
}
```

If this works, gradually add back sections to identify the problem.

## Current Status

### Server Status
- ✅ Mintlify dev server is running
- ✅ Accessible at http://localhost:3000
- ⚠️ Warnings about missing files (but files exist)
- ❌ Pages returning Error 500

### Files Verified
- ✅ `docs/mint.json` exists and is valid JSON
- ✅ `docs/introduction.mdx` exists
- ✅ `docs/quickstart.mdx` exists
- ✅ `docs/installation.mdx` exists
- ✅ All subdirectory MDX files exist

### Warnings (Non-Critical)
- "Legacy mint.json detected, auto-upgrading to docs.json"
- "aws/overview is referenced but file does not exist" (FALSE - file exists)
- "aws/kubernetes is referenced but file does not exist" (FALSE - file exists)

## Recommended Steps

1. **Hard refresh browser** (Ctrl+F5)
2. **Try incognito mode** to rule out browser cache
3. **Check browser console** (F12) for JavaScript errors
4. **Restart dev server** with cache clear:
   ```bash
   # Stop server (Ctrl+C)
   cd docs
   rm -f docs.json  # Remove auto-generated config
   npx mintlify dev
   ```

5. **If still failing**, try accessing specific pages:
   - http://localhost:3000/introduction
   - http://localhost:3000/quickstart
   - http://localhost:3000/concepts/architecture

6. **Check for specific error messages** in:
   - Browser console (F12 → Console tab)
   - Network tab (F12 → Network tab)
   - Terminal where Mintlify is running

## Alternative: Revert to Separate Directories

If the unified structure continues to cause issues, we can revert to the original structure:

```bash
# Create mintlify-docs directory
mkdir mintlify-docs

# Move interactive docs back
mv docs/*.mdx mintlify-docs/
mv docs/concepts mintlify-docs/
mv docs/aws mintlify-docs/
# ... etc

# Move mint.json back
mv docs/mint.json mintlify-docs/

# Run from mintlify-docs
cd mintlify-docs
npx mintlify dev
```

## Getting Help

If none of these solutions work:
1. Check Mintlify documentation: https://mintlify.com/docs
2. Check Mintlify GitHub issues: https://github.com/mintlify/mint
3. Join Mintlify Discord: https://mintlify.com/community

## Debug Information

**Environment**:
- OS: Windows
- Shell: PowerShell/Bash
- Node version: (run `node --version`)
- Mintlify version: (check terminal output)

**File Structure**:
```
docs/
├── mint.json ✅
├── introduction.mdx ✅
├── quickstart.mdx ✅
├── installation.mdx ✅
├── concepts/ ✅
├── aws/ ✅
├── azure/ ✅
└── ... (all other directories) ✅
```

**Server Output**:
- Server starts successfully
- Shows "preview ready" message
- Accessible at http://localhost:3000
- But pages return Error 500
