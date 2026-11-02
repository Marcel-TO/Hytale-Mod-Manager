# HytaleModding-mod-manager
This repository contains a small tool that uploads selected mods to newest versions. The tool is currently in early development and only supports CurseForge, but support for ModTale and other platforms will be added in the future. The tool is written in Go and uses the CurseForge API to upload mods. The tool is designed to be run as a command line application, but a GUI may be added in the future.

## Prerequisites
- Go 1.18 or higher
- A CurseForge account and API key

## Configuration
By default, the tool reads `publish-config.yaml` from the current working directory. Use `--config <file>` after the command to select another YAML file. The configuration file should contain the following fields:

```yaml
mods:
  - name: "My Mod"
    repoLocation: "local/path/to/your/mod/repository"
    releaseType: "pre-release"
    platforms:
      curseforge:
        projectId: 123456
  - name: "Another Mod"
    repoLocation: "local/path/to/your/other/mod/repository"
    releaseType: "pre-release"
    platforms:
      curseforge:
        projectId: 789012

```

## Environment Variables
The following environment variables are used for authentication:
- `CURSEFORGE_API_KEY`: Your CurseForge API key.

The tool reads `.env` from the current working directory by default. This file is optional if the variables are already set in the environment. Use `--env <file>` after the command to select another environment file. A file selected with `--env` must exist and be valid. Existing environment variables take precedence over values in the file.

Keep your API key secure and do not share it publicly. Files named `.env` or `.env.<name>` are ignored by Git, except `.env.example`. If you use another name, make sure Git ignores that file too.

## Usage
Both `publish` and `update` accept `--config` and `--env`. The flags are independent, and relative paths are resolved from the current working directory. Quote paths that contain spaces.

```shell
go run main.go publish --config publish-config.alternate.yaml --env .env.alternate
go run main.go update --version 1.2.3 --copy=false --commit=false --publish=false --config publish-config.alternate.yaml --env .env.alternate
```

Use the separate `publish-custom` recipe to select a configuration file and an optional environment file:

```shell
just publish-custom publish-config.alternate.yaml
just publish-custom "folder with spaces\publish-config.yaml" "folder with spaces\.env.alternate"
```

To update and publish with custom files, use `update-and-publish-custom`:

```shell
just update-and-publish-custom 1.2.3 publish-config.alternate.yaml
just update-and-publish-custom 1.2.3 "folder with spaces\publish-config.yaml" "folder with spaces\.env.alternate" false false
```

The environment file defaults to `.env` for both custom recipes and must exist. For `update-and-publish-custom`, copying and committing default to `true`. The last two arguments can override these defaults.

The original recipes forward any extra flags to the Go command. For `just update` and `just update-and-publish`, list every positional value before the flags:

```shell
just publish --config publish-config.alternate.yaml
just update 1.2.3 false false false --config publish-config.alternate.yaml --env .env.alternate
```

Run `go run main.go publish --help` or `go run main.go update --help` to see the available options.

The tool expects the mod repositories to be structured in a specific way, with a `gradle.properties` file containing the mod metadata and a built JAR file located in the `build/libs` directory. The following fields are expected in the `gradle.properties` file:

- `mod_name`: The name of the mod.
- `version`: The version of the mod.
- `hytale_version`: The version of Hytale that the mod is compatible with.

They are used to generate the display name for the mod file on CurseForge, which follows the format: `mod_name-version+hytale_version`. For example, if your `gradle.properties` file contains:

```
mod_name=MyMod
version=1.0.0
hytale_version=1.2.3
```

The display name for the mod file on CurseForge would be `MyMod-1.0.0+1.2.3`.

> **Note**: It is expected that the mod file is built and located in the `build/libs` directory of the mod repository, with the name following the same format: `displayName.jar`. For example, if your display name is `MyMod-1.0.0+1.2.3`, the expected mod file would be located at `build/libs/MyMod-1.0.0+1.2.3.jar`.
