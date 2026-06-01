# Demo Data

A SQL script populates the app with sample data so you can see it in action right away:

- **4 teams** — The Park Rangers, Central Park Striders, Bethesda Brigade, and The Ramblers
- **12 players** — 3 per team
- **1 admin** — full access to `/admin`
- **5 hunts** — Central Park challenges worth 100–500 points

## Quick Start

```sh
sqlite3 database.db < demo_data/data.sql
```

## Login Credentials

| Role | personal_id | Notes |
|------|-------------|-------|
| Admin | `admin` | Access the admin panel at `/admin` |
| Player | `alex`, `maya`, `james` | Team: The Park Rangers |
| Player | `sofia`, `liam`, `emma` | Team: Central Park Striders |
| Player | `noah`, `olivia`, `ethan` | Team: Bethesda Brigade |
| Player | `isabella`, `lucas`, `charlotte` | Team: The Ramblers |

## Hunts

| # | Title | Points | Description |
|---|-------|--------|-------------|
| 1 | The Imagine Mosaic | 100 | Snap a photo of the "Imagine" mosaic in Strawberry Fields |
| 2 | Bow Bridge Selfie | 200 | Selfie on the Bow Bridge with The Lake behind you |
| 3 | Belvedere Castle View | 300 | Photo from the top showing the Great Lawn and skyline |
| 4 | Bethesda Fountain & Angel | 400 | Photograph the Angel of the Waters statue |
| 5 | The Mall Literary Walk | 500 | Find Hans Christian Andersen statue and read a story with him |

## For a Real Event

Skip the sample data. Insert an admin member manually, then log in at `/admin` and create teams, members, and hunts through the web interface.

```sh
sqlite3 database.db "INSERT INTO members (name, personal_id, is_admin, created_at, updated_at) VALUES ('Admin', 'admin', 1, datetime('now'), datetime('now'));"
```