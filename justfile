# Justfile for the Hytale Mod Publisher application

# Publishes the mods
publish:
    go run main.go publish

# Builds and publishes the mods
build-and-publish:
    go run main.go publish --build

# Updates the mods to specific version and publishes them
update version:
    go run main.go update --version {{version}}
