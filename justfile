# Justfile for the Hytale Mod Publisher application

# Publishes the mod
publish:
    go run main.go publish

# Builds and publishes the mod
build-and-publish:
    go run main.go publish --build
