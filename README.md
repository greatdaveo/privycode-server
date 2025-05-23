
# 🔐 PrivyCode

PrivyCode is a secure platform that allows developers to share **read-only access to their private GitHub repositories** with recruiters or collaborators — without making them public or exposing secrets.

---

## 🚀 Features

- 🔗 Generate expiring viewer links to private GitHub repositories
- 👁️ Allow recruiters to browse your code — no GitHub login required
- ✂️ Read-only access — no forking or editing
- 📦 Track view limits and expiration per link
- 🧑‍💻 Developer dashboard to manage links
- 🔄 Light/dark theme support
- 📝 Copy, edit, delete links with ease

---

## 🧰 Tech Stack

| Frontend               | Backend              | Database      |
|------------------------|----------------------|---------------|
| React + TypeScript     | Go (net/http + GORM) | PostgreSQL    |
| TailwindCSS            | GitHub OAuth2 API    |               |
| Monaco Editor          | JWT-style token auth |               |

---

## 🖼️ Live Demo
[privycode.com](https://privycode.com)

---

## 🔐 Authentication Flow

- Users log in via GitHub OAuth
- A secure token is stored in localStorage
- Authenticated users can generate, view, edit, and delete their viewer links

---

## 🔧 Getting Started

```bash
git clone https://github.com/greatdaveo/privycode-server 
cd privycode-server
go mod tidy
go run main.go
````

Create a `.env` file with:

```
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
GITHUB_CALLBACK_URL=http://localhost:8080/github/callback
DATABASE_URL=your_postgres_url
GO_ENV=development
```

---

## 🛣️ API Endpoints

### 👤 Auth

| Method | Endpoint           | Description              |
| ------ | ------------------ | ------------------------ |
| GET    | `/github/login`    | Redirect to GitHub OAuth |
| GET    | `/github/callback` | GitHub redirects here    |
| GET    | `/me`              | Get logged-in user info  |

---

### 📋 Viewer Links (auth required)

| Method | Endpoint           | Description                       |
| ------ | ------------------ | --------------------------------- |
| GET    | `/dashboard`       | Get all viewer links for the user |
| POST   | `/generate-link`   | Create a new viewer link          |
| PUT    | `/update-link/:id` | Update a viewer link's limits     |
| DELETE | `/delete-link/:id` | Delete (soft) a viewer link       |

---

### 🌐 Public Access

| Method | Endpoint                        | Description                        |
| ------ | ------------------------------- | ---------------------------------- |
| GET    | `/view/:token`                  | View repo contents (public access) |
| GET    | `/view-files/:token/file?path=` | View specific file content         |
| GET    | `/view-folder/:token?path=`     | Browse inside subfolders           |
| GET    | `/view-info/:token`             | Get repo + owner info for header   |

> ✅ Recruiters only need the `/view/:token` link — no login required.

---

## 🧠 Data Models

### User

```go
type User struct {
  ID              uint
  GitHubUsername  string
  GitHubToken     string
  Email           string
}
```

### ViewerLink

```go
type ViewerLink struct {
  ID         uint
  RepoName   string
  Token      string
  MaxViews   int
  ViewCount  int
  ExpiresAt  time.Time
  UserID     uint
}
```

---


## ✅ Future Improvements

* Analytics per link (view history, time opened)
* AI-powered repo summaries
* GitHub repo insights integration

---

## 👨‍💻 Developed By
> Olowomeye David [GitHub](https://github.com/greatdaveo)

---

```
