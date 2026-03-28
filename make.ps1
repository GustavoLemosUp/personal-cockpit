# ============================================================
# Personal Cockpit — Script de build para Windows
# Uso: .\make.ps1 <comando>
#
# Comandos disponíveis:
#   install       Instala dependências do frontend
#   dev           Inicia em modo desenvolvimento (hot-reload)
#   build         Build completo (frontend + executável)
#   build-frontend  Só o build do React
#   test          Roda os testes
#   test-v        Roda os testes com saída detalhada
#   test-cover    Roda os testes e gera coverage.html
#   clean         Remove artefatos de build
# ============================================================

param([string]$cmd = "help")

function Invoke-Install {
    Write-Host "Instalando dependencias do frontend..." -ForegroundColor Cyan
    Set-Location frontend
    npm install
    Set-Location ..
}

function Invoke-BuildFrontend {
    Write-Host "Buildando frontend..." -ForegroundColor Cyan
    Set-Location frontend
    npm run build
    if ($LASTEXITCODE -ne 0) { Write-Host "Erro no build do frontend" -ForegroundColor Red; exit 1 }
    Set-Location ..
}

function Invoke-Build {
    Invoke-BuildFrontend
    Write-Host "Buildando executavel..." -ForegroundColor Cyan
    wails build
}

function Invoke-Dev {
    Write-Host "Iniciando modo desenvolvimento..." -ForegroundColor Cyan
    wails dev
}

$TEST_PKGS = @("./config/...", "./database/...", "./services/...")

function Invoke-Test {
    Write-Host "Rodando testes..." -ForegroundColor Cyan
    go test $TEST_PKGS
}

function Invoke-TestV {
    Write-Host "Rodando testes (verbose)..." -ForegroundColor Cyan
    go test -v $TEST_PKGS
}

function Invoke-TestCover {
    Write-Host "Rodando testes com cobertura..." -ForegroundColor Cyan
    go test -coverprofile=coverage.out $TEST_PKGS
    go tool cover -html=coverage.out -o coverage.html
    Write-Host "Relatorio gerado: coverage.html" -ForegroundColor Green
    Start-Process coverage.html
}

function Invoke-Clean {
    Write-Host "Limpando artefatos..." -ForegroundColor Cyan
    Remove-Item -ErrorAction SilentlyContinue coverage.out, coverage.html
    Remove-Item -Recurse -ErrorAction SilentlyContinue build\bin\
    Write-Host "Pronto." -ForegroundColor Green
}

function Show-Help {
    Write-Host ""
    Write-Host "Uso: .\make.ps1 <comando>" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  install        Instala dependencias do frontend (npm install)"
    Write-Host "  dev            Inicia modo desenvolvimento com hot-reload"
    Write-Host "  build          Build completo (frontend + executavel)"
    Write-Host "  build-frontend Apenas build do React"
    Write-Host "  test           Roda os testes"
    Write-Host "  test-v         Testes com saida detalhada"
    Write-Host "  test-cover     Testes + relatorio de cobertura"
    Write-Host "  clean          Remove artefatos de build"
    Write-Host ""
}

switch ($cmd) {
    "install"        { Invoke-Install }
    "dev"            { Invoke-Dev }
    "build"          { Invoke-Build }
    "build-frontend" { Invoke-BuildFrontend }
    "test"           { Invoke-Test }
    "test-v"         { Invoke-TestV }
    "test-cover"     { Invoke-TestCover }
    "clean"          { Invoke-Clean }
    default          { Show-Help }
}
