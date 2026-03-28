# ============================================================
# Personal Cockpit — Makefile
# ============================================================

.PHONY: dev build build-frontend install test test-v test-cover clean

## Instala dependências do frontend (rodar uma vez após clonar)
install:
	cd frontend && npm install

## Build apenas do frontend (gera frontend/dist/)
build-frontend:
	cd frontend && npm run build

## Gera o executável em build/bin/ (faz build do frontend automaticamente)
build: build-frontend
	wails build

## Inicia o app em modo desenvolvimento com hot-reload
dev:
	wails dev

# Pacotes testáveis (exclui o pacote main — requer frontend/dist do build)
TEST_PKGS := ./config/... ./database/... ./services/...

## Roda todos os testes
test:
	go test $(TEST_PKGS)

## Roda testes com saída detalhada (verbose)
test-v:
	go test -v $(TEST_PKGS)

## Roda testes com relatório de cobertura
test-cover:
	go test -coverprofile=coverage.out $(TEST_PKGS)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório gerado: coverage.html"

## Roda testes de um pacote específico (ex: make test-pkg PKG=./services)
test-pkg:
	go test -v $(PKG)

## Remove artefatos de build e cobertura
clean:
	rm -f coverage.out coverage.html
	rm -rf build/bin/
