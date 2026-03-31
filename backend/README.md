# Backend

This directory contains the Go/Fiber backend for Snapteil.

Configuration lives in `backend/.env`. Start from `backend/.env.example`, then adjust any reviewer-facing flags such as `ENABLE_SWAGGER`.

For project setup, running instructions, Docker usage, and the overall feature summary, see the main README:

- [Main README](../README.md)

For architectural decisions, tradeoffs, and assumptions behind the backend and the full application, see:

- [Architecture Notes](../ARCHITECTURE_NOTES.md)

For API route details, reviewers should check the Swagger endpoint when the backend is running and Swagger is enabled:

- `http://localhost:4000/swagger`
