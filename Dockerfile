# Dockerfile for Server Application
FROM golang:1.19

# Set the working directory
WORKDIR /app
# Copy server application code
COPY Golang .

# Install dependencies

RUN go get


# Set environment variables, if any
ENV PORT=8090

# Expose the server port
EXPOSE $PORT
# Specify the command to run the server application
CMD ["go", "run", "main.go"]