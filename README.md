# 📈 **Price Tracker - Group 19**

Full-stack application for price tracking. This project uses a unified Docker-based environment to ensure that all developers work with the same configuration without the need to install complex local dependencies.

## 🛠 Technologies

*   **Backend:** Go (Framework Gin) + GORM
*   **Frontend:** React + TypeScript + Vite
*   **Database:** PostgreSQL
*   **Infrastructure:** Docker & Docker Compose
---

## 🚀 Prerequisites

To run this project on your machine, **you only need to have installed**:

1.  [Git](https://git-scm.com/)
2.  [Docker Desktop](https://www.docker.com/products/docker-desktop/) (o Docker Engine + Docker Compose)

*Note: You do not need to install Go, Node.js, or PostgreSQL on your computer. Docker takes care of everything.*

---

## ⚙️ **Installation and Setup**

Follow this steps to run the project for the first time:

1. **Clone the repository:**
```bash
git clone <URL_DEL_REPOSITORIO>
cd Grupo-19
```

2. **Start the containers:**
   Run the following command in the project root. This will download the necessary images, install dependencies, and automatically create the database with its tables.
```bash
docker compose up --build
```

3. **Verification:**
   The backend will be ready when you see these messages in the terminal:
   * `✅ PostgreSQL connection established correctly!`
   * `INFO Tables synchronized correctly in PostgreSQL`

---

## 🔗 URLs for Local Access 

Once the containers are running, you can access the services at the following addresses:

| Service | Local URL | Description |
| :--- | :--- | :--- |
| **Frontend** | [http://localhost:5173](http://localhost:5173) | User Interface (Vite/React). |
| **Backend API** | [http://localhost:8080/api/hello](http://localhost:8080/api/hello) | Go server test endpoint. |
| **Health Check**| [http://localhost:8080/health](http://localhost:8080/health) | Server status check. |

---

## 💻 Development guide

The environment is configured with **Hot Reload**. This means you do not need to restart Docker constantly:

*   **Frontend:** If you modify any file inside `frontend/src/`, the browser will refresh automatically.
*   **Backend:** If you modify any `.go` file in `backend/`, the *Air* tool will recompile and restart the server instantly.

### Useful Commands

**See database tables:**
To enter PostgreSQL and list the created tables (`\dt`), open a new terminal and run:
```bash
docker compose exec db psql -U admin -d rastreador_precios -c "\dt"
```

**Install new dependencies in Go:**
If you need to add a new library to the backend, **do not use** your local Go. Ask the container to update the `go.mod`:
```bash
docker compose exec backend go get <package-name>
```

**Install new dependencies in Frontend (npm):**
Similarly, to install Node packages:
```bash
docker compose exec frontend npm install <package-name>
```

**Test Execution**
To run the tests, navigate to the folder where they are located and run:
```bash
go test -v
```

## 📚 Technical Documentation

All detailed documentation regarding architecture, data models, and design decisions is organized in the `/docs` directory. Click on the following links to read more:

* 🏛️ **[System Architecture](./docs/arquitecture.md)**: Data flow and Client-Server layer diagram.
* 💾 **[Data Models](./docs/data-models.md)**: Entity-Relationship diagram and PostgreSQL database dictionary.
* 🧠 **[Design Decisions](./docs/design-decisions.md)**: Technical justification for the use of Go, React, and the Scraping strategy.
* 📖 **[User Guide](./docs/user-guide.md)**: User manual and frequent troubleshooting.
