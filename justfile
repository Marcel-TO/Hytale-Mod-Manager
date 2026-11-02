# Justfile for the Hytale Mod Publisher application

# Publishes the mods
[positional-arguments]
publish *args:
    go run main.go publish "$@"

# Publishes the mods using custom configuration and environment files
[positional-arguments]
publish-custom config env=".env":
    go run main.go publish --config "$1" --env "$2"

# Builds and publishes the mods
[positional-arguments]
build-and-publish *args:
    go run main.go publish --build "$@"

# Updates the mods to specific version and publishes them
[positional-arguments]
update version copy="false" commit="false" publish="false" *args:
    shift 4; go run main.go update --version {{version}} --copy={{copy}} --commit={{commit}} --publish={{publish}} "$@"

# Updates the mods to specific version and publishes them
[positional-arguments]
update-and-publish version copy="true" commit="true" *args:
    shift 3; go run main.go update --version {{version}} --copy={{copy}} --commit={{commit}} --publish=true "$@"

# Updates and publishes the mods using custom configuration and environment files
[positional-arguments]
update-and-publish-custom version config env=".env" copy="true" commit="true":
    go run main.go update --version "$1" --config "$2" --env "$3" --copy="$4" --commit="$5" --publish=true