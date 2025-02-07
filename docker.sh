#!/bin/bash

# Build the containers
build() {
    case "$1" in
        -f)
            docker build -f Dockerfile.frontend -t frontend:latest .
            ;;
        -b)
            docker build -f Dockerfile.backend -t backend:latest .
            ;;
        *)
            echo "Usage: $0 build {-f|-b}"
            exit 1
            ;;
    esac
}

# Run the containers
run() {
    case "$1" in
        -f)
            docker run -it --rm -p 8081:80 frontend:latest
            ;;
        -b)
            docker run -it --rm -p 3000:3000 backend:latest
            ;;
        *)
            echo "Usage: $0 run {-f|-b}"
            exit 1
            ;;
    esac
}

# Main script
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 {build|run} {-f|-b}"
    exit 1
fi

case "$1" in
    build)
        build "$2"
        ;;
    run)
        run "$2"
        ;;
    *)
        echo "Usage: $0 {build|run} {-f|-b}"
        exit 1
        ;;
esac
