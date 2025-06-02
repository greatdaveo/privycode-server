
# 🔐 PrivyCode

This is the **Go backend** of the PrivyCode project.

PrivyCode is a secure platform that allows developers to share **read-only access to their private GitHub repositories** with recruiters or collaborators - without making them public or exposing secrets.

---

## 🚀 Features

- Generate expiring viewer links to private GitHub repositories
- Allow recruiters to browse your code - no GitHub login required
- Read-only access - no forking or editing
- Track view limits and expiration per link
- Developer dashboard to manage links
- Copy, edit, delete links with ease

---

## 🧰 Tech Stack

| Frontend               | Backend              | Database      |
|------------------------|----------------------|---------------|
| React + TypeScript     | Go (net/http + GORM) | PostgreSQL    |
| TailwindCSS            | GitHub OAuth2 API    |               |
| Monaco Editor          | JWT-style token auth |               |

---

## 🖼️ Live Demo
[https://privycode.com](https://privycode.com)

---

## 🔐 Authentication Flow

- Users log in via GitHub OAuth
- A secure token is stored in localStorage
- Authenticated users can generate, view, edit, and delete their viewer links

---

## 🔧 Getting Started

```bash
# Clone the repository
git clone https://github.com/greatdaveo/privycode-server 
cd privycode-server

# Install dependencies
go mod tidy

# Run the server
go run main.go
````
Create a `.env` file in the project root:

```
PORT=8080
DATABASE_URL=your_postgres_url
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
GITHUB_CALLBACK_URL=http://localhost:8080/github/callback
GO_ENV=development # For Production only
FRONTEND_URL=http://localhost:5173 or your frontend url

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

### 📋 Viewer Links (Auth Required)

| Method | Endpoint           | Description                       |
| ------ | ------------------ | --------------------------------- |
| GET    | `/dashboard`       | Get all viewer links for the user |
| POST   | `/generate-link`   | Create a new viewer link          |
| PUT    | `/update-link/:id` | Update an existing viewer link    |
| DELETE | `/delete-link/:id` | Soft delete a viewer link         |

---

### 🌐 Public Access (No Auth Required)

| Method | Endpoint                        | Description                        |
| ------ | ------------------------------- | ---------------------------------- |
| GET    | `/view/:token`                  | View repo contents (public access) |
| GET    | `/view-files/:token/file?path=` | View a specific file content       |
| GET    | `/view-folder/:token?path=`     | Browse inside folders & subfolders |
| GET    | `/view-info/:token`             | Get repo & owner info for display  |

> ✅ Recruiters only need the `/view/:token` link - no login required.

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

## 🤝 Contributing
Contributions are welcome!
If you'd like to suggest features or report bugs, feel free to fork the project, open an issue and possibly submit a pull request.

---

## 👨‍💻 Developed By
> Olowomeye David [GitHub](https://github.com/greatdaveo) [LinkedIn](https://linkedin.com/in/greatdaveo)

---

```
