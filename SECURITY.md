# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| `main` (latest) | Yes |
| Older tags | No |

Only the current `main` branch receives security fixes. Users are encouraged to stay on the latest release.

## Reporting a Vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

Report vulnerabilities privately via a [GitHub Security Advisory](https://github.com/BenjamenMeyer-inspectiv/go-t3/security/advisories/new).

Include as much detail as possible:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Any suggested mitigations

## Response Timeline

| Action | Target |
|--------|--------|
| Acknowledgement | Within 3 business days |
| Status update | Within 10 business days |
| Patch / fix | Within 30 days of confirmation |

We will coordinate disclosure with the reporter before publishing any advisory.

## Out of Scope

- Vulnerabilities in dependencies (report those upstream)
- Issues requiring physical access to the machine
- Denial of service against the in-memory server (no persistence layer)
