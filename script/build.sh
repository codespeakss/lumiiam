set -e

current_dir="$(cd "$(dirname "$0")" && pwd)"
echo "$current_dir：$current_dir"
verFile=$current_dir"/../version"
Ver=$(cat $verFile) ;
BuildT=$(date -u +'%Y%m%dT%H%M%SZ')
GitBranch=$(git rev-parse --abbrev-ref HEAD)
GitCommit=$(git rev-parse --short HEAD)

rm -rf out/*

AppName=lumiiam
CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o out/$AppName       -ldflags "-X main.Version=$Ver -X main.BuildTime=$BuildT -X main.GitBranch=$GitBranch -X main.GitCommit=$GitCommit" cmd/version.go cmd/$AppName.go ;
echo "module      # build end with code: "$?
