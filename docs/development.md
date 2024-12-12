# Development Guide

## Architecture

The project uses a distributed architecture with the following components:

1. Management Server

   - User authentication and authorization
   - Agent registration management
     - Multi-agent registration and discovery
     - Agent health monitoring
     - Agent metadata management
   - Resource management
   - WebRTC signaling service

2. Agent Nodes

   - Remote desktop service
   - Remote terminal service
   - Local resource management
   - Proxy services

3. Web Client
   - Responsive interface
   - Real-time communication
   - Remote control

### Directory Structure

```bash
/
├── cmd/                     # Executable entry points
│   ├── server.go            # Management server
│   ├── agent.go             # Agent program
│   ├── tools.go             # Tools
├── web/                     # Web frontend
│   ├── src/
│   ├── public/
│   └── package.json
├── internal/               # Internal shared code
│   ├── auth/               # Authentication
│   ├── models/             # Data models
│   ├── proto/              # Communication protocol definitions
│   └── common/             # Common utilities
├── pkg/                    # Exportable packages
│   ├── webrtc/             # WebRTC related
│   │   ├── connection/
│   │   └── signaling/
│   ├── remote/             # Remote control related
│   │   ├── desktop/
│   │   └── terminal/
│   ├── proxy/              # Proxy related
│   │   ├── kubernetes/
│   │   ├── api/
│   │   └── llm/
│   └── agent/              # Agent related
│       ├── registry/
│       └── session/
├── api/                    # REST API definitions
│   └── openapi/
├── configs/                # Configuration files
├── deployments/            # Deployment related
│   ├── docker/
│   └── kubernetes/
├── docs/                   # Documentation
├── scripts/                # Build scripts
├── Makefile                # Makefile
└── main.go                 # Main entry point
```

## Tech Stack

- Go 1.23
- WebRTC
- WebSocket

## Development Setup

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -s -m 'feat(server): Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Guidelines

- Use conventional commits format: `type(scope): message`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Keep commits atomic and focused
- Write clear, descriptive commit messages
- Use `git commit -s` to sign off on commits

### Code Style

- Follow the project's lint rules
- Write meaningful variable and function names
- Add comments for complex logic
- Include test cases for new features
