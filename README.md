# 🌳 Menu Tree System - Backend API

Production-ready RESTful API for managing hierarchical menu tree structures with unlimited nesting depth, written in **Go (Golang)** using **Gin Web Framework**, **GORM**, and **PostgreSQL**.

---

## 🚀 Tech Stack & Architecture Decisions

| Component | Technology | Rationale |
| :--- | :--- | :--- |
| **Language** | Go (1.22+) | High performance, type safety, low memory footprint (Preferred in assessment) |
| **Framework** | Gin Gonic | Lightweight, ultra-fast HTTP web framework with robust routing & middleware |
| **ORM** | GORM | Powerful ORM with Auto-Migration, hooks, and clean relational mappings |
| **Database** | PostgreSQL 16 | ACID-compliant with native UUID support and recursive querying capabilities |
| **API Docs** | Swagger / OpenAPI | Auto-generated interactive documentation via `swaggo` |
| **Architecture**| Clean Layered | Strict separation of concerns: **Handler ➜ Service ➜ Repository ➜ Database** |

---

## 📐 Architecture & Key Design Highlights

### 1. In-Memory Tree Construction ($O(n)$)
Instead of executing expensive recursive SQL queries on every read, all menu items are fetched in a single optimized query and transformed into an in-memory hierarchical tree using a hash map lookup in **$O(n)$ time complexity**.

### 2. Infinite Nesting & Relational Integrity
- **Self-referencing Foreign Key**: `parent_id` points to `menus(id)` with `ON DELETE CASCADE`. Deleting any parent automatically purges all descendants.
- **Circular Dependency Guard**: Prevents cyclic hierarchy traps (e.g. a node cannot become a child of itself or any of its descendants).
- **Depth Calculation**: Automatically computes the `depth` level upon node creation or subtree movement.

---

## 🗄️ Database Schema

```sql
CREATE TABLE menus (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    parent_id  UUID REFERENCES menus(id) ON DELETE CASCADE,
    depth      INT NOT NULL DEFAULT 0,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_menus_parent_id ON menus(parent_id);
```

---

## 📡 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/menus` | Get full hierarchical menu tree |
| `GET` | `/api/menus/:id` | Get single menu item with parent details |
| `POST` | `/api/menus` | Create a new menu item (root or child) |
| `PUT` | `/api/menus/:id` | Update menu item name |
| `DELETE` | `/api/menus/:id` | Delete menu item and all nested children |
| `PATCH` | `/api/menus/:id/move` | Move menu item to a different parent |
| `PATCH` | `/api/menus/:id/reorder` | Reorder item among siblings |
| `GET` | `/health` | Server health check endpoint |
| `GET` | `/swagger/index.html`| Interactive Swagger UI |

---

## 🛠️ How to Run

### Option 1: Run with Docker Compose (Recommended)
Make sure Docker is running on your machine:

```bash
# Start PostgreSQL & Backend services
docker compose up --build

# Backend will be available at http://localhost:8080
# Swagger docs: http://localhost:8080/swagger/index.html
```

### Option 2: Run Locally
1. **Ensure PostgreSQL is running** with a database named `menus_db`.
2. Copy environment file:
   ```bash
   cp .env.example .env
   ```
3. Run the application:
   ```bash
   go mod download
   go run main.go
   ```

---

## 📖 Interactive API Documentation (Swagger)
Once the server is running, visit:
👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**
