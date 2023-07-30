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

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

# Clear the Go module cache before getting the packages
RUN go clean -modcache

# Copy only `.go` files, if you want all files to be copied then replace `with `COPY . .` for the code below.
COPY . .

EXPOSE 8080

# Install our third-party application for hot-reloading capability.
RUN go mod tidy
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

# Compile the application and make it executable
RUN go build -o ./allocator ./cmd/allocator
RUN chmod +x ./allocator

ENTRYPOINT CompileDaemon -polling -log-prefix=false -build="go build -o ./allocator ./cmd/allocator" -command="./allocator" -directory="./"
