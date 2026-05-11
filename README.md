# HytaleModding-mod-manager
This repository contains a small tool that uploads selected mods to newest versions. The tool is currently in early development and only supports CurseForge, but support for ModTale and other platforms will be added in the future. The tool is written in Go and uses the CurseForge API to upload mods. The tool is designed to be run as a command line application, but a GUI may be added in the future.

## Prerequisites
- Go 1.18 or higher
- A CurseForge account and API key

## Configuration
It can be configured with a YAML file named `publish-config.yaml`. The configuration file should be placed in the same directory as the executable. The configuration file should contain the following fields:

```yaml
curseforge:
  mods:
    - projectId: 123456
      name: "My Mod"
      repoLocation: "local/path/to/your/mod/repository"
    - projectId: 789012
      name: "Another Mod"
      repoLocation: "local/path/to/your/other/mod/repository"
```

## Environment Variables
The following environment variables are used for authentication:
- `CURSEFORGE_API_KEY`: Your CurseForge API key.

> **Note**: Make sure to keep your API key secure and do not share it publicly. Put it in a `.env` file or set it in your system environment variables.

## Usage
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
