# Base image for building the go project
FROM golang:1.14-alpine AS build

# Updates the repository and installs git
RUN apk update && apk upgrade && \
    apk add --no-cache git

# Switches to /tmp/app as the working directory, similar to 'cd'
WORKDIR /tmp/app

## If you have a go.mod and go.sum file in your project, uncomment lines 13, 14, 15

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Builds the current project to a binary file called nazwa
# The location of the binary file is /tmp/app/out/nazwa
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./out/nazwa .

#########################################################

# The project has been successfully built and we will use a
# lightweight alpine image to run the server 
FROM alpine:latest

# Adds CA Certificates to the image
RUN apk add ca-certificates

# Copies the binary file from the BUILD container to /app folder
COPY --from=build /tmp/app/out/nazwa /app/nazwa
COPY --from=build /tmp/app/setup /app/setup
COPY --from=build /tmp/app/statics /app/statics
COPY --from=build /tmp/app/setup/female.png /app/upload/profile/female.png
COPY --from=build /tmp/app/setup/thumb.female.png /app/upload/profile/thumbnail/female.png
COPY --from=build /tmp/app/setup/male.png /app/upload/profile/male.png
COPY --from=build /tmp/app/setup/thumb.male.png /app/upload/profile/thumbnail/male.png
COPY --from=build /tmp/app/setup/no-photo.png /app/upload/product/no-photo.png
COPY --from=build /tmp/app/setup/thumb.no-photo.png /app/upload/product/thumbnail/no-photo.png
COPY --from=build /tmp/app/.env /app/.env


# Switches working directory to /app
WORKDIR "/app"

# Exposes the 8080 port from the container
EXPOSE 8080

# Runs the binary once the container starts
CMD ["./nazwa"]