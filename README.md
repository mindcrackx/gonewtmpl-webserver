# gonewtmpl-webserver

My opinionated skeleton repository for a go web app.  


---
### Start new Project:  
`go install golang.org/x/tools/cmd/gonew@latest`  
`gonew github.com/mindcrackx/gonewtmpl-webserver github.com/myuser/myrepo`  


---
### Setup pre-commit
https://pre-commit.com/  
`pip install pre-commit`  
`pre-commit install`  
`pre-commit autoupdate`  


---
### Run the server locally with hot-reload:  
`air`  


---
### Run dev environment with docker:  
`cd ./deployments/dev`  
`docker compose up -d --build`  
The dev environment is a docker-compose with the containerized server plus preconfigured prometheus + grafana + dashboard.  
Grafana: `http://localhost:3000` -> user: `admin` passwd: `supersecret`


---
### Run sqlc to generate database code:  
`./tools/sqlc.sh`  


---
### For commits use "Conventional Commits"
https://www.conventionalcommits.org/  

---
### Description
This project is the skeleton template for a go web app.  
It is my default boilerplate:
* CI/CD with github-workflows and local testing with pre-commit
* cmd/server/main.go -> configuration and startup with graceful shutdown
* deployments -> docker-compose for local development with containers
* Dockerfile -> multi-stage builder for the webserver
* internal/metrics -> prometheus metrics (for grafana dashboard see deployments/dev/grafana_myapp_dashboard.json)
* internal/server -> "actual" server
* internal/db
    * migration -> db-migration-scripts (up and down) for [golang-migrate](https://github.com/golang-migrate/migrate)
    * query -> sql commands for [sqlc](https://sqlc.dev)
* tools/sqlc.sh -> generating go boilerplate from sql with [sqlc](https://sqlc.dev)
* ui -> structure for static and template files. Gets embedded into the go binary. My current preferences are [htmx](https://htmx.org/docs/), [tailwind](https://tailwindcss.com/docs/installation) and if necessary [alpinejs](https://alpinejs.dev/start-here) and some plain javascript.
