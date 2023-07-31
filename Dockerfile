FROM golang:1.18

ARG GITLAB_TOKEN
ENV GITLAB_TOKEN=${GITLAB_TOKEN}
ENV GOPRIVATE=gitlab.hotel.tools

# Set up netrc for private git authentication
RUN echo -e "machine gitlab.hotel.tools\nlogin gitlab-ci-token\npassword ${GITLAB_TOKEN}" > ~/.netrc

# Set the correct permissions on the .netrc file
RUN chmod 600 ~/.netrc

# Copy application data into image
COPY . /app
WORKDIR /app

RUN go mod download

# Compile the application and make it executable
RUN go build -o ./allocator ./cmd/allocator
RUN chmod +x ./allocator

# Clear the Go module cache before getting the packages
RUN go clean -modcache

# Install our third-party application for hot-reloading capability.
RUN go mod tidy
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

EXPOSE 8080

ENTRYPOINT CompileDaemon -polling -log-prefix=false -build="go build -o ./allocator ./cmd/allocator" -command="./allocator" -directory="./"
