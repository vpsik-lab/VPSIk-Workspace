# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| Open Source (latest) | ✅ |
| VPSIk Pro | ✅ (license required) |
| WorkSpace SaaS | ✅ (managed by VPSIk) |

## Reporting a Vulnerability

We take security seriously. If you discover a security issue, please report it privately.

**Email:** opensource@vpsik.com

Do not open public GitHub issues for security vulnerabilities.

### What to include

- Description of the vulnerability
- Steps to reproduce
- Affected versions
- Potential impact
- Suggested fix (if any)

### Response time

We will acknowledge your report within 48 hours and provide a timeline for the fix.

## Security Best Practices

1. Always use the latest version
2. Change default admin password immediately after installation
3. Use Cloudflare Tunnel to avoid exposing ports 80/443
4. Keep your VPS OS and Docker up to date
5. Enable automatic backups via `vpsik backup --all`

## Responsible Disclosure

We appreciate researchers and users who report vulnerabilities responsibly. We will acknowledge your contribution in the release notes (unless you prefer to remain anonymous).
