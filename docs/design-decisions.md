# Design Decisions

This document outlines the most important architectural and technical decisions made during the project's development.

## 1. Three-Tier Architecture
* **Decision:** Separation of concerns between the Frontend (React), Backend (Go), and Database (PostgreSQL).
* **Motivation:** To ensure scalability, allow independent service deployment, and facilitate maintenance. Communication via REST API ensures the frontend remains decoupled from backend implementation details.

## 2. Backend Language: Go
* **Decision:** Using Go (Golang) with the Gin framework.
* **Motivation:** High performance and concurrency. Go's ability to handle background processes (`goroutines`) is ideal for running the price scraper without blocking user requests.

## 3. Infrastructure: Docker
* **Decision:** Using Docker and Docker Compose for the entire development environment.
* **Motivation:** To eliminate the "it works on my machine" problem. It standardizes the versions of Go, Node.js, and PostgreSQL for all contributors.

## 4. Scraping Strategy (Limitations)
* **Decision:** On-demand and scheduled scraping.
* **Known Limitation:** Scraping complex websites is subject to security blocks. We have implemented a system that handles errors gracefully to ensure server stability is not compromised.