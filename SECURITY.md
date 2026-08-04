# Security Policy

## Reporting Vulnerabilities

If you discover a security vulnerability, please report it responsibly:

1. **Do NOT open a public GitHub issue**
2. Email security details to the maintainers
3. Include reproduction steps if possible
4. Allow reasonable time for a fix before public disclosure

## Security Measures

### Authentication
- JWT tokens with RS256 signing
- Token validation on all API endpoints
- Configurable localhost bypass (can be disabled via `localhost_bypass = false` in config)
- WebSocket endpoints validate tokens on upgrade

### File Operations
- All file paths are sanitized to prevent path traversal
- Files confined to allowed directories: `/DATA`, `/var/lib/casaos`, `/tmp`, `/etc/samba`
- Soft-delete moves files to `.casaos-trash/` instead of permanent deletion
- SSRF protection blocks private/reserved IP ranges

### Network
- CORS origins configurable via `cors_origins` in config
- Rate limiting: 100 requests per minute per IP
- SSH credentials sent over WebSocket, not URL query params

### Input Validation
- HTML-escaping on all error messages embedded in responses
- Path sanitization on all file upload/download endpoints
- Docker command injection prevention in shell scripts

## Dependency Updates

Known vulnerabilities in dependencies:
- `nwaples/rardecode` and `mholt/archiver/v3`: Transitive deps from `CasaOS-Common`, not exploitable (archiver replaced with stdlib in our code)
- Frontend: `socket.io-client` v2, `yamljs` v0.3 have known CVEs - planned upgrades

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for security-related changes.
