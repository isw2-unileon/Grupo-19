# Arquitectura del Sistema

The system is designed using a three-tier architecture (Client-Server), separating the user interface, the business logic, and the storage system. 

Communication between the Frontend and the Backend is handled via a REST API. Client and Server are decoupled and exchange information through HTTP requests, ensuring a secure, agile, and scalable data flow.

## Tecnologías utilizadas:

**Frontend**: React (with Vite and TypeScript) to build a fast and structured user interface.
**Backend**: Go (with Gin framework) to handle server-side logic with high performance.
**Database**: PostgreSQL as a relational system to store information securely.
**Infrastructure**: Docker for containerization and Render for automated cloud deployment.