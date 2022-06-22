
FROM golang:1.16-alpine

# Adds argument for API key
ARG NASA_API_KEY

# Set working directory
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY main.go ./
COPY static/ ./static

# Compile binaries
RUN go build -o ./apod

# Specify internal port
EXPOSE 8080

# Run the built app!
CMD [ "./apod" ]
