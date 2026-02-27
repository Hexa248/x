# Hotel Booking Workspace

This repository contains a Go-based hotel booking sample application in `hotel-booking/`.

## Run locally

1. Change into the app directory:
   ```bash
   cd hotel-booking
   ```
2. Start one of the three app entry points:
   - **User app** (public browsing/booking):
     ```bash
     go run .
     ```
   - **Admin app** (admin dashboard):
     ```bash
     go run ./cmd/admin
     ```
   - **Staff app** (staff dashboard/rooms):
     ```bash
     go run ./cmd/staff
     ```
3. Open the app URL:
   - User: `http://localhost:8080` (or `PORT`)
   - Admin: `http://localhost:8081` (or `ADMIN_PORT`)
   - Staff: `http://localhost:8082` (or `STAFF_PORT`)

## Quick checks

Run package checks with:

```bash
go test ./...
```
