
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
    - [ ] Registration (macOS-only)
- [x] CDN
    - [x] Serving uploaded files
    - [x] Serving thumbnails
- [ ] Web
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
    - [ ] Account page functionality
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
        - [ ] Exporting pools

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

### Docker Setup

The project supports separate Docker containers for each service component:

- **web** - Web interface container
- **api** - API service container
- **cdn** - CDN service container

Here's an example `docker-compose.yml` to get you started:

```yaml
services:
  web:
    build:
      context: .
      dockerfile: Dockerfile
      target: web
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/.data
    env_file:
      - ".env"
    restart: unless-stopped

  api:
    build:
      context: .
      dockerfile: Dockerfile
      target: api
    ports:
      - "8081:8081"
    volumes:
      - ./data:/app/.data
    env_file:
      - ".env"
    restart: unless-stopped

  cdn:
    build:
      context: .
      dockerfile: Dockerfile
      target: cdn
    ports:
      - "8082:8082"
    volumes:
      - ./data:/app/.data
    env_file:
      - ".env"
    restart: unless-stopped
```
