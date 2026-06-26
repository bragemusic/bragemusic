
# Configuration

Below is a description of the configuration parameters. You can see the toml and the env ways to enter the data.

## Server Config
This config is used to define the server. It only accepts ENV values at the moment.
### ACOUST_ID
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| api_key | string | BM_ACOUST_ID_API_KEY | API key to Acoust ID. Used to identify the files to not rely solely on ID3. |
### ADMIN
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| email | string | BM_ADMIN_EMAIL | Default user, with admin rights, email. Defaults to 'admin@example.com' |
| username | string | BM_ADMIN_USERNAME | Default user, with admin rights, username. Defaults to 'admin' |
| password | string | BM_ADMIN_PASSWORD | Default user, with admin rights, password. Defaults to 'password' |
### ANALYSER
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| base_url | string | BM_ANALYSER_BASE_URL | URL to the analysis service. If left blank, no analysis will be performed. |
### JOBS
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| importer | string | BM_JOBS_IMPORTER | How often the importer will look for new media files. Cron expression. Defaults to '*/3 * * * *' |
| meta_syncer | string | BM_JOBS_META_SYNCER | How often the meta-syncer will sync the needed metadata. Cron expression. Defaults to '*/3 * * * *' |
| search_items | string | BM_JOBS_SEARCH_ITEMS | How often the search items will be updated. Cron expression. Defaults to '*/3 * * * *' |
| analyser | string | BM_JOBS_ANALYSER | How often the track analysis items will be updated. Cron expression. Defaults to '*/3 * * * *' |
| token_cleanup | string | BM_JOBS_TOKEN_CLEANUP | How often token cleanup will be performed. Cron expression. Defaults to '*/10 * * * *' |
### PATHS
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| config_dir | string | BM_PATHS_CONFIG_DIR | Dir where server files are stored. |
| image_dir | string | BM_PATHS_IMAGE_DIR | Dir where image assets are stored. |
| music_dir | string | BM_PATHS_MUSIC_DIR | Dir where music files are stored. |
| import_dir | string | BM_PATHS_IMPORT_DIR | Dir where imported albums and tracks will be saved before processing. |
| manual_import_dir | string | BM_PATHS_MANUAL_IMPORT_DIR | Dir where the importer is looking for manually added bulk imports. |
| backup_import_dir | string | BM_PATHS_BACKUP_IMPORT_DIR | Dir where imported albums and tracks will be saved after processing. |
### WIKIPEDIA
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| email | string | BM_WIKIPEDIA_EMAIL | Used against the wikipedia API. They require a valid email to make sure you behave. |
| name | string | BM_NAME | Name of the server. Defaults to 'Brage Music Server' |
| port | int | BM_PORT | Port of the server. Defaults to 3000. |
## Client Config
The following configuration is used in the local client. The daemon and the GUI are both accepting the same config.
It is stored in ~/.config/brage/config.toml or use it with ENV variables.
### AUTH
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| token | string | BM_AUTH_TOKEN | Token to authenticate with the server. Required if running in daemon mode. |
### GENERAL
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| player_name | string | BM_GENERAL_PLAYER_NAME | The name of your device. This is what will be shown in the connect funtionality. |
| player_icon | string | BM_GENERAL_PLAYER_ICON | The icon of your player will have. Choose from 'laptop', 'computer', 'phone', 'speaker', 'tv', 'generic'. Defaults to 'laptop'. |
| disable_transitions | bool | BM_GENERAL_DISABLE_TRANSITIONS | Set this to true if you see weird artefacts on popups, overlays and more. Will disable all transitions. |
| client_type | string | BM_GENERAL_CLIENT_TYPE | Either sync or streaming. Sync client is syncing all files to the client and can run offline. The streaming client streams everything from the server. |
| log_level | string | BM_GENERAL_LOG_LEVEL | Set the log level of the client. Defaults to 'INFO'. |
### PATHS
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| config_dir | string | BM_PATHS_CONFIG_DIR |  |
| image_dir | string | BM_PATHS_IMAGE_DIR |  |
| music_dir | string | BM_PATHS_MUSIC_DIR |  |
### SERVER
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| base_url | string | BM_SERVER_BASE_URL |  |
### THEME.theme_name
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| name | string | BM_THEME_theme_name_NAME | The name of the theme. Used in the GUI for human readable selection. |
| color_scheme | string | BM_THEME_theme_name_COLOR_SCHEME | Define if your color scheme is light or dark. Nothing else is accepted. |
### THEME.theme_name.COLORS
| Parameter | Type | ENV Variable | Descripion |
| --------- | ---- | ------------ | ---------- |
| accent | string | BM_THEME_theme_name_COLORS_ACCENT | Main brand color |
| accent_foreground | string | BM_THEME_theme_name_COLORS_ACCENT_FOREGROUND | Foreground color on components with accent background. |
| background | string | BM_THEME_theme_name_COLORS_BACKGROUND | Background color of the main frame. |
| border | string | BM_THEME_theme_name_COLORS_BORDER | Border color. |
| danger | string | BM_THEME_theme_name_COLORS_DANGER | Error message color. |
| danger_foreground | string | BM_THEME_theme_name_COLORS_DANGER_FOREGROUND | Foreground color on components with error background. |
| default | string | BM_THEME_theme_name_COLORS_DEFAULT | Used in many places. Its used for cancel buttons, hovers, background on fallback icons and background in input fields. |
| default_foreground | string | BM_THEME_theme_name_COLORS_DEFAULT_FOREGROUND | Used for player buttons, excluding play/pause. |
| field_foreground | string | BM_THEME_theme_name_COLORS_FIELD_FOREGROUND | Text color in forms. |
| foreground | string | BM_THEME_theme_name_COLORS_FOREGROUND | Main text color. |
| overlay | string | BM_THEME_theme_name_COLORS_OVERLAY | Background in popup windows. |
| success | string | BM_THEME_theme_name_COLORS_SUCCESS | Success message color. |
| success_foreground | string | BM_THEME_theme_name_COLORS_SUCCESS_FOREGROUND | Foreground color on components with success background. |
| surface | string | BM_THEME_theme_name_COLORS_SURFACE | Background color of player, cards and notifications. |
| surface2 | string | BM_THEME_theme_name_COLORS_SURFACE2 | Background color of the menu. |
| warning | string | BM_THEME_theme_name_COLORS_WARNING | Warning message color. |
| warning_foreground | string | BM_THEME_theme_name_COLORS_WARNING_FOREGROUND | Foreground color on components with warning background. |