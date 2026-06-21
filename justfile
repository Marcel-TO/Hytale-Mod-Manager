# Justfile for the Hytale Mod Publisher application

# Publishes the mods
publish:
    go run main.go publish

# Builds and publishes the mods
build-and-publish:
    go run main.go publish --build

# Updates the mods to specific version and publishes them
update version copy="false" commit="false" publish="false":
    go run main.go update --version {{version}} --copy={{copy}} --commit={{commit}} --publish={{publish}}

# Updates the mods to specific version and publishes them
update-and-publish version copy="false" commit="false":
    go run main.go update --version {{version}} --copy={{copy}} --commit={{commit}} --publish=true
