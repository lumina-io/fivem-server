export COMPOSE_BAKE=true
# Server Config
export USE_DEV="false"

# Auto restart
# Use for Production Environment
# export RESTART_POLICY="unless-stopped"

## NOTE: DO NOT EDIT (use for directory name)
export SERVER_ID="lumina"
export SERVER_PATH="/app/fivem/txData/${SERVER_ID}"

# Database Config
#export MARIADB_VERSION="lts"
export MYSQL_ROOT_PASSWORD=""
export MYSQL_DATABASE="lumina_gta"
export MYSQL_USER="app"
export MYSQL_PASSWORD=""

export MYSQL_ADDRESS="localhost"
export MYSQL_PORT="23306"

export DATABASE_URL="mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@${MYSQL_ADDRESS}:${MYSQL_PORT}/${MYSQL_DATABASE}?charset=utf8mb4"

export FXSERVER_LICENSE_KEY=""

export STEAM_API_KEY=""
export SCREENSHOT_BASIC_TOKEN=""
export FIVEMANAGE_API_KEY=""
export TEBEX_SECRET=""

# rcon
export RCON_ENABLED="YES"
export RCON_ADDRESS="localhost"
export RCON_PORT="30120"
export RCON_PASSWORD="changeme_for_fxc0n"
