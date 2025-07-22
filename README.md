# Sentiric API Gateway Service

**Description:** A unified API Gateway and/or Backend-for-Frontend (BFF) layer for all Sentiric microservices, providing a single entry point for external API consumers.

**Core Responsibilities:**
*   Receiving all incoming API requests from external clients and routing them to the appropriate internal microservices.
*   Enforcing centralized authentication and authorization (AuthN/AuthZ) for all API requests.
*   Managing cross-cutting concerns such as rate limiting, caching, and request/response transformations.
*   Simplifying client-side integration by aggregating multiple internal APIs into a single, cohesive interface.

**Technologies:**
*   Node.js (e.g., Express, Fastify) or Go (e.g., Fiber, Gin)
*   HTTP proxy libraries.
*   Authentication/Authorization libraries (e.g., JWT).
*   Fwe can test open-source Traefik [https://github.com/traefik/traefik/] 

**API Interactions (As an API Provider & Client):**
*   **As a Provider:** Exposes the public-facing APIs for `sentiric-dashboard-ui`, `sentiric-web-agent-ui`, `sentiric-embeddable-voice-widget-sdk`, `sentiric-cli`.
*   **As a Client:** Calls APIs of all other internal microservices (e.g., `sentiric-user-service`, `sentiric-dialplan-service`, `sentiric-agent-service`, `sentiric-stt-service`, `sentiric-tts-service`, `sentiric-cdr-service`).

**Local Development:**
1.  Clone this repository: `git clone https://github.com/sentiric/sentiric-api-gateway-service.git`
2.  Navigate into the directory: `cd sentiric-api-gateway-service`
3.  Install dependencies: `npm install` (Node.js) or `go mod tidy` (Go).
4.  Create a `.env` file from `.env.example` to configure routing rules, internal service URLs, and authentication settings.
5.  Start the service: `npm start` (Node.js) or `go run main.go` (Go).

**Configuration:**
Refer to `config/` directory and `.env.example` for service-specific configurations, including routing tables, security policies, and caching settings.

**Deployment:**
Designed for containerized deployment (e.g., Docker, Kubernetes). Often deployed with load balancers for high availability. Refer to `sentiric-infrastructure`.

**Contributing:**
We welcome contributions! Please refer to the [Sentiric Governance](https://github.com/sentiric/sentiric-governance) repository for coding standards and contribution guidelines.

**License:**
This project is licensed under the [License](LICENSE).
