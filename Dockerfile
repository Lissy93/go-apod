
FROM golang:1.18-alpine AS build
# Adds argument for API key
ARG NASA_API_KEY
# Set working directory
WORKDIR /app
# Copy over all files
COPY . .
# Download Go modules
RUN go mod download
# Compile binaries
RUN go build -o ./apod

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/apod ./
# Specify internal port
EXPOSE 8080
# Run the built app!
CMD [ "./apod" ]
