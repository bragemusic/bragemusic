#!/usr/bin/env sh
set -e

mkdir -p ${BM_PATHS_MUSIC_DIR}
mkdir -p ${BM_PATHS_IMAGE_DIR}
mkdir -p ${BM_PATHS_CONFIG_DIR}
mkdir -p ${BM_PATHS_IMPORT_DIR}
mkdir -p ${BM_PATHS_MANUAL_IMPORT_DIR}
mkdir -p ${BM_PATHS_BACKUP_IMPORT_DIR}

DATABASE_URL=sqlite:${BM_PATHS_CONFIG_DIR}/data.db hermox --service=migrations -- dbmate -d ./migrations up
./bragemusic-server
