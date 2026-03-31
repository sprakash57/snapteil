# Snapteil

Snapteil is photo sharing app that lets you upload images, tag them and filter them by tags. Other users will be notified about new posts live.

## Demo

https://drive.google.com/file/d/1CHt5Jh7uj9rgezwA_iaU1NcKZhzz2qaO/view?usp=sharing

## Prerequisites

- Docker and Docker Compose
- For local development without Docker:
- Go `1.26+`
- Node.js `22+`

## Technology Stack

- Frontend: React 19, TypeScript, Vite, Tailwind CSS
- Backend: Go, Fiber v3, WebSocket, Swagger docs
- Image processing: `disintegration/imaging`, AVIF/WebP decoders
- Containerization: Docker, Docker Compose, Nginx
- Testing: Go test, Vitest

## Quick Run

The fastest way to run the project is with Docker Compose:

```bash
docker compose up --build
```

Services:

- Client: `http://localhost:5173`
- Backend API: `http://localhost:4000`
- Swagger: `http://localhost:4000/swagger`

To stop everything:

```bash
docker compose down
```

If you want Docker to recreate the uploads volume from scratch:

```bash
docker compose down -v
```

## Run, Build, and Test

### Client

Install dependencies:

```bash
cd client
npm install
```

Run locally:

```bash
npm run dev
```

Build:

```bash
npm run build
```

Test:

```bash
npm run test:run
```

### Backend

Run locally:

```bash
cd backend
air
# or
go run .
```

`air` is configured through `backend/.air.toml` for hot reloading during local development.

Build:

```bash
go build -o server .
```

Test:

```bash
go test ./...
```

## API Routes

The backend routes are organized under `/api/v1`.

For a complete list of endpoints, request/response shapes, and parameters, review the Swagger UI:

- `http://localhost:4000/swagger`

## Features

- Image upload with title and tag validation
- Mobile first neumorphic UI design
- Accessible upload form
- Infinite scroll feed with multi-select tag filtering
- Supported uploads: JPEG, PNG, GIF, WebP, AVIF
- Server-side image type detection
- Image normalization with max-dimension resizing
- Real-time upload notifications over WebSockets
- Swagger API documentation
- Dockerized frontend and backend setup

## Project Structure

```text
.
├── backend
│   ├── config
│   ├── data
│   ├── docs
│   ├── handlers
│   ├── models
│   ├── routes
│   ├── services
│   └── uploads
├── client
│   ├── src
│   └── public
└── docker-compose.yml
```

## Notes

- Initial images are read from `backend/data/seed.json`.
- Uploaded images are written to `backend/uploads` locally or to the Docker volume when using Compose.
- A separate design note is available in [ARCHITECTURE_NOTES.md](https://github.com/sprakash57/snapteil/blob/main/ARCHITECTURE_NOTES.md).
