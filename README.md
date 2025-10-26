# SubPaper Backend

Backend for the **SubPaper Android App**
Play Store: [SubPaper: Reddit Wallpapers](https://play.google.com/store/apps/details?id=in.saikat.subpaper)

## Overview

Backend for the SubPaper Android application. Migrated from Python Flask to a faster Go implementation using goroutines and channels. Stateless. No database.

## Features

* Go and Gin framework
* Concurrency via goroutines and channels
* Token bucket rate limiting
* Authorization middleware
* CORS support
* GZip compressed JSON responses

## Getting Started

#### Requirements

* Go 1.25+

#### Installation

```bash
git clone https://github.com/dis70rt/SubPaper-backend.git
cd SubPaper-backend
go mod download
```

#### Environment Variables

```env
REDDIT_CLIENT_ID=
REDDIT_CLIENT_SECRET=
REDDIT_USER_AGENT=
REDDIT_USERNAME=
REDDIT_PASSWORD=
PORT=8080
API_SECRET_KEY=
```

#### Run

```bash
go run main.go
```

#### API Endpoints

| Method | Endpoint                   | Description                  | Auth |
| ------ | -------------------------- | ---------------------------- | ---- |
| GET    | `/health`                  | Health check                 | No   |
| GET    | `/metrics`                 | Metrics output               | No   |
| GET    | `/api/v1/wallpaper/search` | Fetch wallpapers from Reddit | Yes  |

## Authorization

```
Authorization: Bearer <token>
```

## Contributing

PRs accepted.

