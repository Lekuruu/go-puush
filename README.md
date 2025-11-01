
<p align="center">
    <img src="https://raw.githubusercontent.com/Lekuruu/go-puush/refs/heads/main/web/static/img/toplogo.png">
</p>

---

**go-puush** is a go implementation of the [puush](https://puush.me) file sharing service, providing file upload, sharing, and management capabilities.  
The service consists of three main components:

- **API** - Handles communication with the puush client  
- **CDN** - Serves uploaded files and thumbnails  
- **Web Interface** - Provides a web-based user interface, closely matching [the original](https://puush.me)

A key advantage is the ease of use, as this project uses only Go and SQLite. No external dependencies are required!

## Progress

The following features have been implemented or are planned for the future:

- [x] API
    - [x] Authentication
    - [x] Upload
    - [x] History
    - [x] Delete
    - [x] Thumbnail
<!-- - [ ] Registration (macOS-only & unused) -->
- [x] CDN
    - [x] Serving uploaded files
    - [x] Serving thumbnails
- [x] Web
    - [x] Public pages (Home, Login, Register, About, ...)
    - [x] Account page(s)
    - [x] Gallery page
    - [x] Login functionality
    - [x] Logout functionality
    - [x] Registration functionality
        - [x] Registration
        - [x] Email activation
        - [x] Invite codes
    - [x] Password reset functionality
    - [x] Account page functionality
        - [x] Move uploads
        - [x] Delete uploads
        - [x] Update default pool
        - [x] Reset API key
        - [x] Switch between views
        - [x] Search for uploads
        - [x] Change password
        - [x] Username check
        - [x] Username claiming
        - [x] "Stop asking about my username"

## Setup

### Quick Start

If you want to quickly try out go-puush without setting up a development environment:

1. Go to the [GitHub Actions page](https://github.com/Lekuruu/go-puush/actions)
2. Click on the latest successful workflow run
3. Download the artifact for your platform
4. Extract and run the binary

It will download all the required files off of GitHub automatically.

### Development Setup

For development with hot reloading support:

1. **Install Air**

   ```bash
   go install github.com/air-verse/air@latest
   ```

2. **Clone the repository**

   ```bash
   git clone https://github.com/Lekuruu/go-puush.git
   cd go-puush
   ```

3. **Install dependencies**

   ```bash
   go mod download
   ```

4. **Run with hot reloading**

   ```bash
   air
   ```

   Air will automatically rebuild and restart the application when you make changes to the code.

### Docker Compose

For a docker deployment, you can run the prebuilt container published to [ghcr.io](https://ghcr.io):

1. Copy the example environment file and tweak it as needed:

    ```bash
    cp .example.env .env
    ```

2. Create a `docker-compose.yml` file:

    ```yaml
    services:
        puush:
            image: ghcr.io/lekuruu/go-puush:latest
            container_name: puush
            restart: unless-stopped
            env_file:
                - .env
            ports:
                - "${WEB_HOST}:${WEB_PORT}:${WEB_PORT}"
            volumes:
                - ./data:/app/.data
                - ./web:/app/web
    ```

3. Start the stack:

    ```bash
    docker compose up -d
    ```

This setup *won't* require you to have the repository downloaded. The `web` directory will be populated automatically upon first launch, allowing you to edit web assets.
