
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
    - [ ] Password reset functionality
    - [ ] Registration functionality
    - [ ] Account page functionality
        - [x] Move uploads
        - [x] Delete uploads
        - [x] Update default pool
        - [x] Reset API key
        - [x] Switch between views
        - [x] Search for uploads
        - [ ] Change password
