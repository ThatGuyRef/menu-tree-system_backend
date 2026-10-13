# 📋 Technical Assessment Submission Report
## Project: Fullstack Menu Tree System (Hierarchical CRUD & Responsive UI)

---

### 1. 🔗 GitHub Repositories (Public)
The source code has been structured into clean, decoupled repositories following microservices best practices:

* 🌐 **Frontend Repository (Next.js 14 + Tailwind + Zustand)**:  
  👉 **[https://github.com/ThatGuyRef/menu-tree-system_frontend](https://github.com/ThatGuyRef/menu-tree-system_frontend)**

* ⚙️ **Backend Repository (Go / Gin / GORM / PostgreSQL)**:  
  👉 **[https://github.com/ThatGuyRef/menu-tree-system_backend](https://github.com/ThatGuyRef/menu-tree-system_backend)**

---

### 2. 📦 Necessary Files & Artifacts Included

Both repositories are completely self-contained and include all files required for local development and containerized deployment:

| Category | Files Included | Description |
| :--- | :--- | :--- |
| **Backend Core** | `main.go`, `config/`, `models/`, `repositories/`, `services/`, `handlers/`, `routes/`, `middleware/` | Clean 4-tier layered architecture with full CRUD, move, and reorder |
| **Backend Testing** | `services/menu_service_test.go` | Automated unit test suite verifying tree construction & circular dependency guards |
| **Frontend Core** | `src/app/`, `src/components/layout/`, `src/components/menu/`, `src/components/ui/`, `src/store/`, `src/services/`, `src/types/` | Pixel-accurate responsive tree UI matching Figma design |
| **Database & Seeds** | `backend/scripts/seed.sql` | Auto-initialization schema and seed data with exact Figma menu hierarchy |
| **Docker & Infra** | `backend/Dockerfile`, `backend/docker-compose.yml`, `frontend/Dockerfile`, `frontend/.dockerignore` | Multi-stage Docker builds for minimal production image footprint |
| **Environment** | `backend/.env.example`, `frontend/.env.example` | Clear templates for all configurable environment variables |
| **API Docs** | `backend/docs/` (`swagger.json`, `swagger.yaml`, `docs.go`) | Full OpenAPI 3.0 / Swagger documentation |

---

### 3. 📖 Setup & Running Instructions

#### Option A: Running with Docker Compose (Recommended - 1 Command)
```bash
# 1. Clone both repositories into your workspace
git clone https://github.com/ThatGuyRef/menu-tree-system_backend.git
git clone https://github.com/ThatGuyRef/menu-tree-system_frontend.git

# 2. Start PostgreSQL and Backend via Docker Compose
cd menu-tree-system_backend
docker compose up --build -d

# 3. Start Frontend
cd ../menu-tree-system_frontend
docker build -t menu-frontend .
docker run -p 3000:3000 menu-frontend
```

#### Option B: Running Locally for Development

**1. Backend (Go & PostgreSQL)**:
```bash
cd menu-tree-system_backend
cp .env.example .env

# Run database seed
psql "postgres://localhost:5432/menus_db" -f scripts/seed.sql

# Start server
go run main.go
# ➜ Running on http://localhost:8080
# ➜ Swagger Docs on http://localhost:8080/swagger/index.html
```

**2. Frontend (Next.js 14)**:
```bash
cd menu-tree-system_frontend
cp .env.example .env.local

npm install
npm run dev
# ➜ Running on http://localhost:3000 (or available port)
```

---

### 4. 📸 UI & Figma Comparison

The user interface was crafted to faithfully replicate the design specified in the Figma mockup:
* **Figma Reference Link**: [Figma Design Canvas](https://www.figma.com/design/Pa1vc4POl4DGsvDTutwM6k/Fullstack---ReactJS-VueJS?node-id=0-1&p=f)
* **Visual Match Highlights**:
  - **Royal Blue Sidebar (`#0055D6`)** with smooth rounded edges, brand mark, and active white menu capsule (`Menus`).
  - **Tree View Section**: Clean connecting vertical/horizontal guide lines, expandable chevrons, and inline blue `+` action button.
  - **Form Details Panel**: Seamless borderless fields with soft background tint (`#EAECF0`), displaying Menu ID, Depth, Parent Data, and blue pill `Save` action.
  - **Responsive Design**: Mobile drawer navigation with identical single-column stacked layout on mobile screens.

*(To capture live screenshots for demo submissions, launch the frontend at `http://localhost:3000` or `http://localhost:3002` and take full-page captures for desktop & mobile).*

---

### 5. 🧪 Verification & Automated Test Results

#### Automated Unit Tests (Backend)
All core business logic features test suites with zero failures:
```bash
$ go test ./... -v
=== RUN   TestGetMenuTree_Hierarchy
--- PASS: TestGetMenuTree_Hierarchy (0.00s)
=== RUN   TestCreateMenuItem_DepthCalculation
--- PASS: TestCreateMenuItem_DepthCalculation (0.00s)
=== RUN   TestMoveMenuItem_CircularDependencyGuard
--- PASS: TestMoveMenuItem_CircularDependencyGuard (0.00s)
PASS
ok      menu-tree-system/backend/services       0.443s
```

#### API Integration Testing
- `GET /health` ➜ `200 OK`
- `GET /api/menus` ➜ `200 OK` (Hierarchical JSON tree structure)
- `POST /api/menus` ➜ `201 Created`
- `PUT /api/menus/:id` ➜ `200 OK`
- `PATCH /api/menus/:id/move` ➜ `200 OK` (With circular dependency prevention)
- `DELETE /api/menus/:id` ➜ `204 No Content` (With recursive cascading delete)
- `GET /swagger/index.html` ➜ `200 OK` (Interactive documentation)
