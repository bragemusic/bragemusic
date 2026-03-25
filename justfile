
clean-data:
    rm -r ../data/img/albums/*
    rm -r ../data/img/artists/*
    rm -r ../data/music/*
    rm ../data/config/data.db

serve:
    go run -tags fts5 cmd/server/main.go

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
