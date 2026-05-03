# 📈 Rastreador de Precios - Grupo 19

Aplicación full-stack para el seguimiento de precios. Este proyecto utiliza un entorno unificado basado en Docker para garantizar que todos los desarrolladores trabajen con la misma configuración sin necesidad de instalar dependencias locales complejas.

## 🛠 Tecnologías

*   **Backend:** Go (Framework Gin) + GORM
*   **Frontend:** React + TypeScript + Vite
*   **Base de Datos:** PostgreSQL
*   **Infraestructura:** Docker & Docker Compose

---

## 🚀 Requisitos Previos

Para ejecutar este proyecto en tu máquina, **solo necesitas tener instalado**:

1.  [Git](https://git-scm.com/)
2.  [Docker Desktop](https://www.docker.com/products/docker-desktop/) (o Docker Engine + Docker Compose)

*Nota: No es necesario que instales Go, Node.js ni PostgreSQL en tu ordenador. Docker se encarga de todo.*

---

## ⚙️ Instalación y Puesta en Marcha

Sigue estos pasos para arrancar el proyecto por primera vez:

1. **Clonar el repositorio:**
   ```bash
   git clone <URL_DEL_REPOSITORIO>
   cd Grupo-19
   ```

2. **Levantar los contenedores:**
   Ejecuta el siguiente comando en la raíz del proyecto. Esto descargará las imágenes necesarias, instalará las dependencias y creará la base de datos con sus tablas automáticamente.
   ```bash
   docker compose up --build
   ```

3. **Verificación:**
   El backend estará listo cuando veas estos mensajes en la terminal:
   * `✅ PostgreSQL conexion established correctly!`
   * `INFO Tablas sincronizadas correctamente en PostgreSQL`

---

## 🔗 URLs de Acceso Local

Una vez que los contenedores estén corriendo, puedes acceder a los servicios en las siguientes direcciones:

| Servicio | URL Local | Descripción |
| :--- | :--- | :--- |
| **Frontend** | [http://localhost:5173](http://localhost:5173) | Interfaz de usuario (Vite/React). |
| **Backend API** | [http://localhost:8080/api/hello](http://localhost:8080/api/hello) | Endpoint de prueba del servidor Go. |
| **Health Check**| [http://localhost:8080/health](http://localhost:8080/health) | Comprobación de estado del servidor. |

---

## 💻 Guía de Desarrollo

El entorno está configurado con **Hot Reload**. Esto significa que no necesitas reiniciar Docker constantemente:

*   **Frontend:** Si modificas cualquier archivo dentro de `frontend/src/`, el navegador se actualizará automáticamente.
*   **Backend:** Si modificas cualquier archivo `.go` en `backend/`, la herramienta *Air* recompilará y reiniciará el servidor al instante.

### Comandos Útiles

**Ver las tablas de la base de datos:**
Para entrar a PostgreSQL y listar las tablas creadas (`\dt`), abre una nueva terminal y ejecuta:
```bash
docker compose exec db psql -U admin -d rastreador_precios -c "\dt"
```

**Instalar nuevas dependencias en Go:**
Si necesitas añadir una nueva librería al backend, **no uses** tu Go local. Pídeselo al contenedor para que actualice el `go.mod`:
```bash
docker compose exec backend go get <nombre-del-paquete>
```

**Instalar nuevas dependencias en Frontend (npm):**
De igual manera, para instalar paquetes de Node:
```bash
docker compose exec frontend npm install <nombre-del-paquete>
```

---
*Este README se irá actualizando a medida que se añadan nuevas funcionalidades, endpoints y guías de despliegue.*
