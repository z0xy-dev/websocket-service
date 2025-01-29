# This script updates the repository and redeploys the Docker service for the WebSocket application.
# It pulls the latest changes from the remote repository.
# Then it rebuilds and redeploys the services in the background, skipping unnecessary dependencies.
git pull
docker compose up -d --no-deps --build
