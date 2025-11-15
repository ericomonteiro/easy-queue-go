# 🔍 Tracing - Guia Rápido

## Início Rápido

### 1. Configure o ambiente

```bash
cp .env.example .env
```

### 2. Inicie o Jaeger

```bash
docker-compose up -d jaeger
```

### 3. Execute a aplicação

```bash
go run src/internal/cmd/main.go
```

### 4. Acesse o Jaeger UI

Abra: http://localhost:16686

### 5. Gere alguns traces

```bash
# Faça algumas requisições
curl http://localhost:8080/health
```

### 6. Visualize no Jaeger

1. Selecione "easy-queue-go" no dropdown de serviços
2. Clique em "Find Traces"
3. Explore os traces!

## Adicionar Tracing Customizado

```go
import (
    "easy-queue-go/src/internal/tracing"
    "go.opentelemetry.io/otel/attribute"
)

func MinhaFuncao(ctx context.Context) error {
    tracer := tracing.Tracer("meu-componente")
    ctx, span := tracer.Start(ctx, "MinhaOperacao")
    defer span.End()

    // Adicione atributos
    span.SetAttributes(
        attribute.String("user.id", userID),
        attribute.Int("items.count", len(items)),
    )

    // Seu código aqui...
    
    return nil
}
```

## Desabilitar Tracing

No arquivo `.env`:

```bash
TRACING_ENABLED=false
```

## Documentação Completa

Veja: [docs/tracing.md](docs/tracing.md)
