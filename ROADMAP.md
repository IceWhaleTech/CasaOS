# Roadmap

## Completed

### Security (Phase 1-2)
- [x] Path traversal protection on all file endpoints
- [x] SSH credentials sent over WebSocket (not URL)
- [x] WebSocket token validation
- [x] Soft-delete (trash) for file operations
- [x] XSS prevention in error responses
- [x] Configurable localhost auth bypass
- [x] Configurable CORS origins
- [x] Token refresh race condition fix

### UI/UX (Phase 3-5)
- [x] Memory leak fix in AppPanel
- [x] Error handling typo fixes (6x)
- [x] Tag validation bug fix
- [x] Router guard double-next fix
- [x] Hardcoded strings → i18n
- [x] Extended legacy container config form
- [x] Multi-container log viewer with selector

### Code Quality (Phase 7)
- [x] Dead code removal
- [x] Shared utilities (constants.js, volumeHelper.js)
- [x] Driver vet fixes

## In Progress

### Data Protection (Phase 6)
- [ ] Trash cleanup cron job (30-day retention)
- [ ] Database backup before upgrades
- [ ] Config export/import

### Documentation (Phase 9)
- [x] SECURITY.md
- [x] CONTRIBUTING.md
- [x] ROADMAP.md
- [ ] CHANGELOG.md

## Planned

### High Priority
- [ ] Vue 2 → Vue 3 migration (Vue 2 EOL Dec 2023)
- [ ] Buefy → Oruga migration (Buefy unmaintained)
- [ ] Socket.IO v2 → v4 upgrade
- [ ] `yamljs` → `js-yaml` replacement

### Medium Priority
- [ ] httpOnly cookie token storage
- [ ] Docker privilege validation
- [ ] Samba password encryption at rest
- [ ] Structured logging

### Low Priority
- [ ] Container config editing (full)
- [ ] File operation checksum verification
- [ ] Config export/import (JSON/YAML)

## Known Issues

- Dependabot alerts for `rardecode` and `archiver/v3` (transitive, not exploitable)
- `container.launcher` endpoint may have command injection (external dep)
- No sudo access on dev machine (limits system-level testing)
