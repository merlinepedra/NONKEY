
# del generated code 
# Get-ChildItem .\enum\ -Recurse -Include *_gen.go | Remove-Item

################################################################################
# generate enum
Write-Output "generate enums"
genenum -typename=TokenType -packagename=tokentype -basedir=enum 
genenum -typename=ObjectType -packagename=objecttype -basedir=enum 
genenum -typename=Precedence -packagename=precedence -basedir=enum 

goimports -w enum

################################################################################
$DATESTR=Get-Date -UFormat '+%Y-%m-%dT%H:%M:%S%Z:00'
$GITSTR=git rev-parse HEAD
################################################################################
# build bin

$BIN_DIR="bin"
$SRC_DIR="."

mkdir -ErrorAction SilentlyContinue "${BIN_DIR}"

# build bin here
$BUILD_VER="${DATESTR}_${GITSTR}_release_windows"
echo "Build Version: ${BUILD_VER}"
echo ${BUILD_VER} > ${BIN_DIR}/BUILD_windows
go build -o "${BIN_DIR}\nonkey.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\nonkey.go"

$BUILD_VER="${DATESTR}_${GITSTR}_release_linux"
echo "Build Version: ${BUILD_VER}"
echo ${BUILD_VER} > ${BIN_DIR}/BUILD_linux
$env:GOOS="linux" 
go build -o "${BIN_DIR}\nonkey" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\nonkey.go"
$env:GOOS=""
