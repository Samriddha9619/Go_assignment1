# Hotel Price Tracker

A real-time hotel price tracking system built with Go that scrapes hotel listings from MakeMyTrip and Goibibo, stores historical price data, and exposes a RESTful API for accessing hotel information across multiple Indian cities.

## Table of Contents

- Features  
- Tech Stack  
- Architecture  
- Prerequisites  
- Installation & Setup  
- Running the Application  
- API Overview  
- Scraper & Scheduler  
- Environment Variables  
- Project Structure  
- Future Enhancements  
- License  
- Developer  

## Features

- Real-time scraping from dynamic hotel websites  
- Historical price tracking with timestamped entries  
- Automated scheduler for periodic scraping  
- JWT-based authentication and protected routes  
- User registration, login, and profile access  
- Multi-city support with sorting and filtering  
- Scraping logs and fallback to mock data  
- RESTful API with clean endpoints  
- PostgreSQL database with GORM ORM  

## Tech Stack

Backend:  
- Go (Golang)  
- Gin (HTTP web framework)  
- GORM (ORM for PostgreSQL)  

Database:  
- PostgreSQL (hosted on Aiven)  

Scraping:  
- Colly (initial static scraper)  
- ChromeDP (headless browser automation for dynamic content)  
- cdproto (Chrome DevTools Protocol)
  
Authentication:  
- JWT (golang-jwt/jwt)  
- bcrypt (password hashing)  

Scheduling:  
- robfig/cron (cron job scheduler)  

Environment Management:  
- godotenv (load .env files)  

## Architecture

- Gin server handles HTTP requests and routes  
- Handlers process API logic and interact with the database  
- JWT middleware protects private routes  
- Scraper fetches hotel data using ChromeDP  
- Scheduler triggers scraping every 2 hours  
- PostgreSQL stores hotel data, price history, and logs  

## Prerequisites

- Go 1.21 or higher  
- PostgreSQL 15 (local or Aiven cloud)  
- Git  
- Chrome or Chromium (required for scraping)  

## Installation & Setup

1. Clone the repository  
2. Install dependencies using `go mod download`  
3. Create a `.env` file with database and server configuration  
4. Run the application using `go run main.go`  

## Running the Application

- On startup, the app connects to PostgreSQL, runs migrations, and performs an initial scrape  
- The scheduler runs every 2 hours to update hotel prices  
- API is served on the configured port (default: 8080)  

## API Overview

- `/api/auth/register` – Register a new user  
- `/api/auth/login` – Login and receive JWT  
- `/api/profile` – Get user profile (protected)  
- `/api/hotels` – Get all hotels with sorting  
- `/api/hotels/city/:city` – Get hotels by city  
- `/api/hotels/:id/history` – Get price history for a hotel  
- `/api/logs` – View scraping logs (protected)  
- `/api/scrape/trigger` – Manually trigger scraping (protected)  

## Scraper & Scheduler

The scraper was initially built using Colly, a fast and lightweight Go library for scraping static websites. However, hotel platforms like MakeMyTrip and Goibibo load listings dynamically using JavaScript, which Colly cannot handle. To overcome this, the scraper was rebuilt using ChromeDP, which automates a headless Chrome browser and can interact with JavaScript-rendered content.

The scraping flow includes:
- Launching headless Chrome
- Navigating to hotel listing pages
- Waiting for dynamic content to load
- Extracting hotel name, price, rating, and location
- Storing results in the database with upsert logic
- Recording price history for each hotel
- Logging each scraping operation

If scraping fails due to network issues or site structure changes, the system falls back to mock data to ensure frontend continuity.

The scheduler runs every 2 hours using robfig/cron. It triggers scraping for all supported cities and updates the database with fresh prices. Manual scraping can also be triggered via a protected API endpoint.


## Environment Variables

- `DB_HOST` – PostgreSQL host  
- `DB_PORT` – PostgreSQL port  
- `DB_USER` – Database username  
- `DB_PASSWORD` – Database password  
- `DB_NAME` – Database name  
- `SERVER_PORT` – API server port  
- `JWT_SECRET` – JWT signing secret  
- `GIN_MODE` – Gin framework mode (debug/release)  

## Project Structure

```
Go_1/
├── backend/
│   ├── auth/                  # JWT middleware & admin authorization
│   ├── config/                # Configuration loader
│   ├── database/              # Database connection & migrations
│   ├── handlers/              # API endpoint logic
│   ├── models/                # Database models
│   ├── scraper/               # Scraper and scheduler
│   ├── main.go                # Application entry point
│   ├── .env                   # Environment variables (not committed)
│   └── .env.example           # Example environment file
├── frontend/                  # Frontend application (optional)
├── .gitignore                 # Git ignore rules
└── README.md                  # Project documentation
```

## Future Enhancements

- Add more cities  
- Email notifications for price drops  
- Price comparison charts  
- Hotel booking integration  
- Mobile app  
- Admin dashboard  
- User favorites/watchlist  
- Advanced filtering  
- AI-powered price prediction  
- Multi-language support  

