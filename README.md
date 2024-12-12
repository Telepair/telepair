# Telepair

> ⚠️ **Notice**: This project is currently in early development stage. No user features are available yet.

A WebRTC-based remote management platform supporting remote desktop, terminal access, API proxy, Kubernetes management, and LLM service proxy.

## Project Objectives

### For Individuals

- Personal remote resource management platform
- Unified access to personal devices, services, and development resources
- Secure and efficient remote workspace management

### For Enterprise

- Comprehensive remote resource management platform for Dev, Ops, and DevOps teams
- Unified management of:
  - Host terminals
  - Kubernetes clusters
  - Internal APIs
  - LLM service hosting
- Streamlined workflow and enhanced team collaboration

## Features

### Remote Control

- 🖥️ Remote Desktop (Agent mode)
  - Multi-monitor support
  - Adaptive quality settings
  - Clipboard synchronization
- 📟 Remote Terminal (supports both Agent and SSH modes)
  - Full PTY support
  - Session recording
  - Command history
- 🔄 Multiple concurrent sessions
- 🔐 End-to-end encryption
  - WebRTC secure data channels
  - TLS for signaling

### Proxy Services

- 🌐 API Proxy
- ☸️ Kubernetes Multi-cluster Management
  - Pod logs viewing and filtering
  - Pod terminal access
  - Node terminal access
  - kubectl command execution
  - Resource monitoring
- 🤖 LLM Service Proxy
  - Support for popular LLM APIs
  - Rate limiting
  - Cost management

### Management Features

- 👥 User Management
  - Role-based access control (RBAC)
  - Multi-factor authentication
  - SSO integration
- 🔑 Access Control
  - Fine-grained permissions
  - Audit logging
  - Session management
- 📝 Agent Registration
  - Multiple agent registration support
  - Agent grouping and tagging
  - Agent status monitoring
  - Centralized agent management
- 🎯 Resource Aggregation
  - Resource usage metrics
  - Performance monitoring
  - Alert configuration

## Quick Start

More details are available [here](docs/development.md).

## TODO

- [ ] Implement user features for remote desktop and terminal access
- [ ] Enhance API proxy with additional security measures
- [ ] Complete Kubernetes management features
- [ ] Integrate more LLM service providers
- [ ] Improve user management with additional SSO options
- [ ] Optimize performance monitoring and alert configuration

## Similar Projects

- [RustDesk](https://github.com/rustdesk/rustdesk) - Open source virtual desktop infrastructure
- [Gotty](https://github.com/yudai/gotty), [Gotty 1](https://github.com/sorenisanerd/gotty) - Share your terminal over HTTP,
- [Screego](https://github.com/screego/server) - Screen sharing
- [SSHX](https://github.com/ekzhang/sshx) - Web-based SSH client
- [Natpass](https://github.com/lwch/natpass) - Work from home, remote development tool
- [WebTTY](https://github.com/maxmcd/webtty) - Web-based TTY client
- [WebSSH2](https://github.com/billchurch/webssh2) - Web-based SSH2 client
- [shell2http](https://github.com/msoap/shell2http) - Share your terminal over HTTP
- [Apache Guacamole](https://guacamole.apache.org/) - HTML5 web-based remote desktop gateway

## License

This project is licensed under the [MIT License](LICENSE)

## Contact

- Project Lead: [Telepair](mailto:me@telepair.online)
- Project Homepage: [GitHub](https://github.com/telepair/telepair)
