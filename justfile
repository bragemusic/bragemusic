
clean-data:
    rm -r ../data/img/albums/*
    rm -r ../data/img/artists/*
    rm -r ../data/music/*
    rm ../data/config/*

import:
    go run -tags fts5 cmd/importer/main.go

migrate:
    dbmate -d ./db/migrations up

generate-docs:
    go run cmd/docs/main.go

call-api route method="GET":
    curl -X "{{method}}" --header "Authorization: Bearer $(cat ~/.local/state/brage/users/11111111-1111-1111-1111-111111111111/token)" "http://localhost:3000/{{route}}"

login email password long_lived_token="true":
    jq -n \
      --arg email "{{email}}" \
      --arg password "{{password}}" \
      --argjson long_lived_token {{long_lived_token}} \
      '{email: $email, password: $password, long_lived_token: $long_lived_token}' \
    | curl -X POST http://localhost:3000/auth/login \
      -H "Content-Type: application/json" \
      -d @- \
    | jq

build-frontend:
    rm -rf ./assets/frontend/*
    mkdir -p ./assets/frontend
    npm install
    npm install --prefix web/
    npm run build --prefix web/
    cp -r ./web/dist/* ./assets/frontend
    cp ./web/static/* ./assets/frontend/assets

serve:
    mkdir -p ${BM_PATHS_MUSIC_DIR}
    mkdir -p ${BM_PATHS_IMAGE_DIR}
    mkdir -p ${BM_PATHS_CONFIG_DIR}
    mkdir -p ${BM_PATHS_IMPORT_DIR}
    mkdir -p ${BM_PATHS_MANUAL_IMPORT_DIR}
    mkdir -p ${BM_PATHS_BACKUP_IMPORT_DIR}

    DATABASE_URL=sqlite:${BM_PATHS_CONFIG_DIR}/data.db hermox --service=migrations -- dbmate -d ./db/migrations up
    go run -tags fts5 cmd/server/main.go

dev-client:
    wails dev -tags fts5

build-linux-client version:
    wails build -tags fts5 --ldflags="-X 'github.com/bragemusic/core/internal/vars.VERSION={{version}}' -X 'github.com/bragemusic/core/internal/vars.PLATFORM=linux'"
