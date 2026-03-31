# Architecture Notes

## Overview

Snapteil is split into a React client service and a Go/Fiber backend.

- The client is responsible for rendering the feed, handling uploads, and responding to realtime events.
- The backend is responsible for validation, image processing, persistence of uploaded files, feed state, and websocket broadcast.

## Main Architectural Decisions

### 1. Why choose Fiber over other Go frameworks?

I have worked a lot with Express in the past. I believe the fastest way to learn and implement something is to have an analogy. Fiber is lightweight, has an Express-like API that keeps route definitions readable, and performs well in benchmarks. For an assignment of this scope it provides enough structure through middleware chaining and route grouping without the ceremony of heavier frameworks like Gin or Echo. Its first-class WebSocket support via `contrib/websocket` also simplified the real-time requirement.

### 2. Layered backend structure

The backend is organized into:

- `routes`: HTTP route registration
- `handlers`: request/response logic
- `services`: domain logic and image processing
- `models`: shared response/data structures
- `config`: environment-driven configuration

Why:

- It keeps Fiber-specific request code out of the image-processing logic.
- It makes unit testing easier because the most important logic lives in services and small handlers.
- It fits the scale of the assignment without introducing heavy abstractions.

### 3. In-memory feed state plus file-system storage

Uploaded files are stored on disk under `backend/uploads`, while image metadata is kept in memory after loading `backend/data/seed.json`.

Why:

- It is simple and fast to implement.
- It avoids database setup overhead for an assignment.
- It supports the main flows: initial feed, upload, pagination, and live notifications.

Tradeoff:

- Uploaded metadata is not durable across backend restarts unless it is also present in seed data.
- This would not be enough for a real multi-user production system.

### 4. Seed file for initial content

The application loads initial images from `backend/data/seed.json`.

Why:

- It gives the app a non-empty initial state for demos and review.
- It avoids introducing a database migration or seeding flow.
- It makes Docker startup deterministic.

Tradeoff:

- The seed file must stay in sync with the actual files in `backend/uploads`.

### 5. Byte-based image validation

The backend validates uploads by reading file bytes and detecting the real image type instead of trusting the multipart `Content-Type`.

Why:

- It is a meaningful security improvement.
- It prevents obvious spoofing of non-image files.
- It keeps validation aligned with the actual decoders used by the backend.

Tradeoff:

- The backend reads the uploaded file into memory for type detection and processing.
- For larger or more numerous uploads, a streaming or staged approach would be better.

### 6. WebSocket broadcast for realtime updates

When an upload succeeds, the backend broadcasts the new image (and the count if there are more than 1 uploads are available) to connected clients.

Why:

- It makes the assignment more interactive.
- It demonstrates simple realtime behavior without adding external infrastructure.

### 7. Docker-first evaluation path

The project can run locally without Docker, but `docker compose up --build` is the intended quick-start.

Why:

- It gives reviewers a single command to run the full stack.
- It keeps Go, Node, and Nginx environments predictable.
- It ships the initial dataset inside the backend image.

## Security and Hardening Decisions

Authentication was intentionally left out, but the backend still includes a basic hardening layer:

- request recovery middleware
- request IDs
- CORS allowlist
- upload and websocket rate limiting
- upload timeout
- body size limit
- `X-Content-Type-Options: nosniff`
- websocket origin checks
- Swagger behind a configuration flag

Why this level:

- It improves the project materially without pushing it beyond assignment scope.
- It shows awareness of common API risks even without full auth.

What is still intentionally limited:

- no user authentication or authorization
- no persistent DB-backed audit trail
- no antivirus or malware scanning
- no object-storage isolation for uploaded files

## Frontend Decisions

### React + Vite

The frontend uses React with Vite for a fast local workflow and a lightweight build setup.

Why:

- It is productive for assignment work.
- I am well versed with React ecosystem.
- Assignment mentions using React and TypeScript.

## Assumptions

- Upload titles are limited to 100 characters.
- Uploads can include at most 5 tags, and each tag can be at most 24 characters long.
- The current maximum upload size is 10 MB.
- The allowed MIME types are `image/jpeg`, `image/png`, `image/gif`, `image/webp`, and `image/avif`.
- The maximum image dimension is 1920 pixels for either width or height after normalization.
- Uploads are limited to one image per request.
- The assignment only requires image uploads. Other media types such as video or PDF are out of scope.
- A simple in-memory metadata store is acceptable.
- A reviewer values clarity and runnability over production-scale infrastructure.
- Docker volume persistence for uploads is acceptable for demo purposes.

## What I Would Improve With More Time

### 1. Replace in-memory state with persistent storage

I would move image metadata into a database such as Postgres. That would make uploads durable across restarts and remove the need to keep `seed.json` synchronized manually.

### 2. Introduce storage abstraction

Right now the filesystem is used directly. I would add a storage interface so local disk, S3-compatible object storage, or another backend could be swapped without touching handler logic. A more production-ready design would store the original upload untouched and generate resized display variants or thumbnails separately.

### 3. Expand automated testing

The current tests focus on high-value service and handler behavior. With more time I would add:

- end-to-end automation tests
- websocket behavior tests

### 4. Improve deployment separation

Serving uploaded files directly from the application is acceptable for this assignment, but I would normally separate API hosting from file delivery and use object storage plus CDN-style caching.

### 5. Add a caching layer for repeated queries

The feed endpoint currently rebuilds its response on every request. Adding an in-memory or HTTP-level cache (e.g. `Cache-Control` / `ETag` headers, or a lightweight LRU cache in the service layer) would reduce work for repeated page and tag combinations.

## Known Tradeoffs

- `seed.json` is simple but manual.
- The backend service keeps current feed state in memory.
- The websocket layer is intentionally minimal.
- File uploads are validated strongly enough for the assignment, but not to production security standards.
