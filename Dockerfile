# Use the official Ubuntu base image
FROM ubuntu:latest

# Set the working directory
WORKDIR /app

# Install dependencies
RUN apt-get update && apt-get install -y wget curl \
    && wget https://golang.org/dl/go1.17.6.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz \
    && rm go1.17.6.linux-amd64.tar.gz \
    && curl -fsSL https://deb.nodesource.com/setup_16.x | bash - \
    && apt-get install -y nodejs

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"

# 设置数据库连接的环境变量
ENV DATABASE_URL=""

# Copy the Go module files and download dependencies
COPY ShortLinkWizard/go.mod ShortLinkWizard/go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY ShortLinkWizard/ .

# Build the Go application
RUN go build -o main .

# Build the Next.js frontend
WORKDIR /app/front_end
RUN npm install && npm run build

# Install a simple HTTP server to serve the frontend
RUN npm install -g serve

# Expose only the port for the frontend
EXPOSE 3000

# Set the working directory back to the root
WORKDIR /app

# Command to run both the backend and frontend
CMD ["sh", "-c", "DATABASE_URL=${DATABASE_URL} ./main & serve -s front_end/out -l 3000"]
