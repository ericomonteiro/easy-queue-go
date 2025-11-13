# 🚀 Guia de Início Rápido

Este guia irá ajudá-lo a configurar e executar o EasyQueue em sua máquina local.

## 📋 Pré-requisitos

Antes de começar, certifique-se de ter instalado:

- **Go 1.25+** - [Download](https://golang.org/dl/)
- **Docker** - [Download](https://www.docker.com/get-started)
- **Docker Compose** - Geralmente incluído com Docker Desktop
- **Git** - Para clonar o repositório

## 📥 Instalação

### 1. Clone o Repositório

```bash
git clone https://github.com/ericomonteiro/easy-queue-go.git
cd easy-queue-go
```

### 2. Configure o Banco de Dados

Inicie o container PostgreSQL usando Docker Compose:

```bash
docker-compose up -d
```

Isso irá:
- Criar um container PostgreSQL 17
- Configurar o banco de dados `easyqueue`
- Expor a porta `5432` no localhost

**Credenciais padrão:**
```
Host: localhost
Port: 5432
Database: easyqueue
User: easyqueue
Password: easyqueue123
```

### 3. Configure as Variáveis de Ambiente

Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

Edite o arquivo `.env` se necessário:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=easyqueue
DB_PASSWORD=easyqueue123
DB_NAME=easyqueue
DB_SSLMODE=disable
```

### 4. Instale as Dependências

```bash
go mod download
```

### 5. Execute a Aplicação

```bash
go run src/internal/cmd/main.go
```

Ou compile e execute:

```bash
go build -o easyqueue src/internal/cmd/main.go
./easyqueue
```

## ✅ Verificação

Para verificar se tudo está funcionando:

### 1. Verifique o Status do PostgreSQL

```bash
docker ps
```

Você deve ver o container `easy-queue-go-postgres-1` rodando.

### 2. Teste a Conexão com o Banco

```bash
docker exec -it easy-queue-go-postgres-1 psql -U easyqueue -d easyqueue
```

### 3. Verifique os Logs da Aplicação

A aplicação deve exibir logs estruturados indicando:
- ✅ Conexão com o banco de dados estabelecida
- ✅ Pool de conexões inicializado
- ✅ Aplicação rodando

## 🛑 Parando a Aplicação

### Parar a Aplicação Go

Pressione `Ctrl+C` no terminal onde a aplicação está rodando.

### Parar o PostgreSQL

```bash
docker-compose down
```

Para remover também os volumes (dados do banco):

```bash
docker-compose down -v
```

## 🔧 Solução de Problemas

### Erro: "connection refused"

**Problema:** A aplicação não consegue conectar ao PostgreSQL.

**Solução:**
1. Verifique se o container está rodando: `docker ps`
2. Verifique as credenciais no arquivo `.env`
3. Certifique-se de que a porta 5432 não está sendo usada por outro processo

### Erro: "port already in use"

**Problema:** A porta 5432 já está em uso.

**Solução:**
1. Pare qualquer instância local do PostgreSQL
2. Ou altere a porta no `docker-compose.yml`:
```yaml
ports:
  - "5433:5432"  # Usa porta 5433 no host
```

### Erro: "go: module not found"

**Problema:** Dependências não instaladas.

**Solução:**
```bash
go mod tidy
go mod download
```

## 📚 Próximos Passos

Agora que você tem o EasyQueue rodando:

- 📖 Explore a [Estrutura do Projeto](project-structure.md)
- 🗄️ Veja o [Schema do Banco de Dados](database/schema.md)
- 🎯 Entenda a [Visão do Produto](product/overview.md)
- 🔧 Configure a [API](api/endpoints.md)

## 💡 Dicas de Desenvolvimento

### Hot Reload

Para desenvolvimento com hot reload, use [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
air
```

### Debug no VS Code

Adicione ao `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch EasyQueue",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/src/internal/cmd/main.go",
      "env": {},
      "args": []
    }
  ]
}
```

### Logs Estruturados

A aplicação usa Zap para logging estruturado. Para visualizar logs formatados:

```bash
go run src/internal/cmd/main.go | jq
```

---

**Pronto!** Você está preparado para começar a desenvolver com o EasyQueue! 🎉
