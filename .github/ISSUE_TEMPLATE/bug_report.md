---
name: Bug report
about: Create a report to help us improve
title: '[Bug] '
labels: 'bug'
assignees: ''

---

**Describe the bug**

> A clear and concise description of what the bug is.

**To Reproduce**

> Steps to reproduce the behavior:

1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**

> A clear and concise description of what you expected to happen.

**Screenshots**

> If applicable, add screenshots to help explain your problem.

**Desktop (please complete the following information):**

```
 - OS: [e.g. iOS]
 - Browser [e.g. chrome, safari]
 - Version [e.g. 22]
```

**System Time**

> Run `timedatectl` and share the output

```
(timedatectl output here)
```

**Logs**

> Run following command to collect corresponding logs:

```bash
sudo journalctl -xef -u casaos-gateway
sudo journalctl -xef -u casaos-user-service
sudo journalctl -xef -u casaos-local-storage
sudo journalctl -xef -u casaos-app-management
sudo journalctl -xef -u casaos.service
```


**Additional context**

> Add any other context about the problem here.
> 
> If you are a Zimaboard user, make it explicit with when you got your Zimaboard.
