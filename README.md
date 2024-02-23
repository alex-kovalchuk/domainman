# DomainMan

DomainMan is a tool for checking and sorting domains from .md files.

## Build

Run `make build` to build *domainman* file.

## Using

```
Usage of ./domainman:
  -check
        Check busy domains (default true)
  -in string
        Directory with files for processing (default "./var/in")
  -out string
        Directory for results (default "./var/out")
  -skip
        Skip busy domains (default true)
  -sort
        Sort records (default true)
  -sort-len
        Sort records by length (default true)
  -zones string
        Check busy domain zones (separated by commas) (default "com")
```

## Examples:

```bash
mkdir -p var/in
mkdir -p var/out

echo "# Domains
- Group1
    - google
    - qwerty
    - xyz
    - qwerty1234567890
- Group2
    - bb
    - bbbcc
    - aaa
    - bbb
    - aa" > var/in/test.md
```

### Case (default)

```bash
./domainman
```

Log
```
Processing: test.md
Domain xyz.com is busy
Domain google.com is busy
Domain qwerty.com is busy
Domain aa.com is busy
Domain bb.com is busy
Domain aaa.com is busy
Domain bbb.com is busy
Domain bbbcc.com is busy
Done!
```

Result in var/out/test.md
```markdown
# Domains
- Group1
    - qwerty1234567890
- Group2
```

### Case (don't skip busy domains)

```bash
./domainman -skip=false
```

Log
```
Processing: test.md
Domain xyz.com is busy
Domain google.com is busy
Domain qwerty.com is busy
Domain aa.com is busy
Domain bb.com is busy
Domain aaa.com is busy
Domain bbb.com is busy
Domain bbbcc.com is busy
Done!
```

Result in var/out/test.md
```markdown
# Domains
- Group1
    - xyz
    - google
    - qwerty
    - qwerty1234567890
- Group2
    - aa
    - bb
    - aaa
    - bbb
    - bbbcc
```

### Case (no domains check and disable sort by length)

```bash
./domainman -check=false -sort-len=false
```

Log
```
Processing: test.md
Done!
```

Result in var/out/test.md
```markdown
# Domains
- Group1
    - google
    - qwerty
    - qwerty1234567890
    - xyz
- Group2
    - aa
    - aaa
    - bb
    - bbb
    - bbbcc
```
