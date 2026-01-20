# Fleming — Architecture Overview

**Open-source · DeSci · Medical data sovereignty · Patient-owned timeline & consent engine**

Fleming consists of a modular monolith backend written in Go, a React frontend, minimal blockchain components, and supporting infrastructure. The project is structured as a monorepo with clear separation between application code, reusable libraries, generated artifacts, and documentation.

The layout prioritizes:
- Single-binary backend simplicity
- Encapsulation of private logic
- Clear distinction between reusable/public code and generated/tool output
- Self-hostable deployment via Docker Compose
- Easy onboarding for contributors

## Folder Structure

```text
fleming/
├── apps/                           # All runnable / deployable applications
│   ├── backend/                    # Golang modular monolith (core API & business logic)
│   │   ├── cmd/
│   │   │   └── fleming/            # Main executable entry point
│   │   │       └── main.go         # Dependency wiring & server startup
│   │   ├── internal/               # Private application logic — cannot be imported externally
│   │   │   ├── auth/               # Wallet-based authentication & session handling
│   │   │   ├── consent/            # Consent state machines & authorization logic
│   │   │   ├── encryption/         # Hybrid encryption implementation & key management
│   │   │   ├── timeline/           # Timeline event modeling & business rules
│   │   │   ├── graph/              # Medical graph relationships & queries
│   │   │   ├── handlers/           # HTTP handler implementations
│   │   │   ├── service/            # Use-case orchestration & domain services
│   │   │   └── ...                 # additional domain packages
│   │   ├── router.go               # HTTP route definitions & middleware
│   │   ├── Dockerfile              # Multi-stage build for production binary
│   │   └── go.mod                  # Backend module & dependencies
│   └── web/                        # React + Vite single-page application
│       ├── public/                 # Static assets (favicon, manifest, robots.txt…)
│       ├── src/
│       │   ├── components/         # Reusable UI building blocks
│       │   ├── features/           # Feature-specific slices (components, hooks, types)
│       │   ├── hooks/              # Custom React hooks
│       │   ├── lib/                # Shared utilities & helpers
│       │   ├── types/              # TypeScript interfaces & shared types
│       │   └── App.tsx             # Application root & routing
│       ├── Dockerfile              # Vite build → static file serving
│       ├── vite.config.ts          # Vite configuration
│       └── package.json            # Frontend dependencies & scripts
│
├── contracts/                      # Smart contract source code
│   └── ConsentAnchor.sol           # Consent anchoring contract
│
├── pkg/                            # Public / reusable Go packages (intended for external use)
│   ├── protocol/                   # Core protocol definitions & primitives
│   │   ├── consent/                # Consent-related types & validation
│   │   ├── timeline/               # Timeline event models & serialization
│   │   ├── crypto/                 # Cryptographic interfaces & types
│   │   └── types/                  # Shared domain types
│   └── eth/                        # Ethereum-related utilities & bindings
│       └── bindings/               # Generated contract bindings (abigen output)
│           └── consentanchor/      # ConsentAnchor contract bindings
│
├── internal/                       # Shared private Go code (rarely used — prefer per-app internal/)
│
├── migrations/                     # Database schema & migration scripts
│
├── gen/                            # Machine-generated code — never edit manually
│   ├── sqlc/                       # sqlc generated queries & models
│   ├── oapi/                       # OpenAPI-generated types & stubs (if used)
│   └── ...                         # other generated artifacts (mocks, enums…)
│
├── docs/                           # Project documentation & architecture knowledge base
│   ├── ARCHITECTURE.md             # High-level architecture & folder responsibilities (this file)
│   └── ...                         # additional decision logs, diagrams, guides
│
├── scripts/                        # Development & maintenance helper scripts
│
├── compose.yml                     # Docker Compose configuration — defines dev & deployment environment
├── .env.example                    # Example environment variables template
├── .gitignore                      # Git ignore rules
└── README.md                       # Project entry point & quick-start guide
