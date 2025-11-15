# Tracing com OpenTelemetry

Este projeto utiliza OpenTelemetry para tracing distribuído, com Jaeger como backend de visualização.

## 📑 Índice

- [Arquitetura](#arquitetura)
- [Configuração](#configuração)
  - [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Como Usar](#como-usar)
  - [1. Iniciar os Serviços](#1-iniciar-os-serviços)
  - [2. Acessar a UI do Jaeger](#2-acessar-a-ui-do-jaeger)
  - [3. Executar a Aplicação](#3-executar-a-aplicação)
  - [4. Gerar Traces](#4-gerar-traces)
  - [5. Visualizar Traces no Jaeger](#5-visualizar-traces-no-jaeger)
- [O que é Rastreado Automaticamente](#o-que-é-rastreado-automaticamente)
- [Adicionando Spans Customizados](#adicionando-spans-customizados)
  - [Exemplo Básico](#exemplo-básico)
  - [Exemplo com Tratamento de Erros](#exemplo-com-tratamento-de-erros)
  - [Exemplo com Spans Aninhados](#exemplo-com-spans-aninhados)
  - [Exemplo em Handler HTTP](#exemplo-em-handler-http)
  - [Tipos de Atributos Suportados](#tipos-de-atributos-suportados)
- [Desabilitar Tracing](#desabilitar-tracing)
- [Portas Utilizadas](#portas-utilizadas)
- [Troubleshooting](#troubleshooting)
  - [Traces não aparecem no Jaeger](#traces-não-aparecem-no-jaeger)
  - [Erro ao conectar no OTLP endpoint](#erro-ao-conectar-no-otlp-endpoint)
- [Próximos Passos](#próximos-passos)

---

## Arquitetura

- **Instrumentação**: OpenTelemetry SDK
- **Protocolo**: OTLP (OpenTelemetry Protocol) via HTTP
- **Backend**: Jaeger All-in-One
- **Middleware**: otelgin para instrumentação automática do Gin

## Configuração

### Variáveis de Ambiente

```bash
# Habilitar/desabilitar tracing
TRACING_ENABLED=true

# Nome do serviço (aparece no Jaeger)
SERVICE_NAME=easy-queue-go

# Versão do serviço
SERVICE_VERSION=1.0.0

# Ambiente (development, staging, production)
ENVIRONMENT=development

# Endpoint do coletor OTLP (sem http://)
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
```

## Como Usar

### 1. Iniciar os Serviços

```bash
# Iniciar Jaeger e PostgreSQL
docker-compose up -d

# Verificar se os serviços estão rodando
docker-compose ps
```

### 2. Acessar a UI do Jaeger

Abra o navegador em: http://localhost:16686

### 3. Executar a Aplicação

```bash
# Certifique-se de ter um arquivo .env com as configurações
cp .env.example .env

# Executar a aplicação
go run src/internal/cmd/main.go
```

### 4. Gerar Traces

Faça requisições para a aplicação:

```bash
# Health check
curl http://localhost:8080/health

# Outras rotas...
```

### 5. Visualizar Traces no Jaeger

1. Acesse http://localhost:16686
2. Selecione o serviço "easy-queue-go" no dropdown
3. Clique em "Find Traces"
4. Explore os traces gerados

## O que é Rastreado Automaticamente

O middleware `otelgin` rastreia automaticamente:

- **HTTP Requests**: Método, path, status code
- **Timing**: Duração total da requisição
- **Errors**: Erros e stack traces
- **Context Propagation**: Propagação de trace context entre serviços

## Adicionando Spans Customizados

### Exemplo Básico

Para adicionar spans customizados no seu código:

```go
import (
    "context"
    "easy-queue-go/src/internal/tracing"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
)

func MinhaFuncao(ctx context.Context) error {
    // Criar um novo span
    tracer := tracing.Tracer("meu-componente")
    ctx, span := tracer.Start(ctx, "MinhaFuncao")
    defer span.End()

    // Adicionar atributos ao span
    span.SetAttributes(
        attribute.String("user.id", "123"),
        attribute.Int("items.count", 42),
        attribute.Bool("is.success", true),
    )

    // Seu código aqui...
    
    return nil
}
```

### Exemplo com Tratamento de Erros

```go
func ProcessarPedido(ctx context.Context, pedidoID string) error {
    tracer := tracing.Tracer("pedidos")
    ctx, span := tracer.Start(ctx, "ProcessarPedido")
    defer span.End()

    span.SetAttributes(
        attribute.String("pedido.id", pedidoID),
        attribute.String("operation.type", "processar"),
    )

    // Simular processamento
    err := realizarProcessamento(pedidoID)
    if err != nil {
        // Registrar erro no span
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return err
    }

    // Marcar span como bem-sucedido
    span.SetStatus(codes.Ok, "Pedido processado com sucesso")
    return nil
}
```

### Exemplo com Spans Aninhados

```go
func OperacaoCompleta(ctx context.Context) error {
    tracer := tracing.Tracer("exemplo")

    // Span pai
    ctx, parentSpan := tracer.Start(ctx, "OperacaoCompleta")
    defer parentSpan.End()

    // Primeira operação (span filho)
    ctx, span1 := tracer.Start(ctx, "BuscarDados")
    span1.SetAttributes(attribute.String("source", "database"))
    // ... buscar dados ...
    span1.End()

    // Segunda operação (span filho)
    ctx, span2 := tracer.Start(ctx, "ProcessarDados")
    span2.SetAttributes(attribute.Int("records.count", 100))
    // ... processar dados ...
    span2.End()

    // Terceira operação (span filho)
    ctx, span3 := tracer.Start(ctx, "SalvarResultado")
    span3.SetAttributes(attribute.String("destination", "cache"))
    // ... salvar resultado ...
    span3.End()

    parentSpan.SetStatus(codes.Ok, "Operação completa finalizada")
    return nil
}
```

### Exemplo em Handler HTTP

```go
func MeuHandler(c *gin.Context) {
    // O contexto já vem com o span do middleware otelgin
    ctx := c.Request.Context()
    
    tracer := tracing.Tracer("handlers")
    ctx, span := tracer.Start(ctx, "ProcessarRequisicao")
    defer span.End()

    // Adicionar informações da requisição
    span.SetAttributes(
        attribute.String("user.agent", c.GetHeader("User-Agent")),
        attribute.String("request.id", c.GetString("request_id")),
    )

    // Processar a requisição
    resultado, err := processarLogica(ctx)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "Erro ao processar")
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    span.SetStatus(codes.Ok, "Requisição processada")
    c.JSON(200, resultado)
}
```

### Tipos de Atributos Suportados

```go
span.SetAttributes(
    // Strings
    attribute.String("key", "value"),
    
    // Números inteiros
    attribute.Int("count", 42),
    attribute.Int64("timestamp", time.Now().Unix()),
    
    // Números decimais
    attribute.Float64("price", 19.99),
    
    // Booleanos
    attribute.Bool("is_active", true),
    
    // Arrays
    attribute.StringSlice("tags", []string{"tag1", "tag2"}),
    attribute.IntSlice("ids", []int{1, 2, 3}),
)
```

## Desabilitar Tracing

Para desabilitar o tracing (útil em testes ou desenvolvimento):

```bash
TRACING_ENABLED=false
```

## Portas Utilizadas

- **16686**: Jaeger UI
- **4318**: OTLP HTTP receiver
- **4317**: OTLP gRPC receiver

## Troubleshooting

### Traces não aparecem no Jaeger

1. Verifique se o Jaeger está rodando: `docker-compose ps`
2. Verifique os logs do Jaeger: `docker-compose logs jaeger`
3. Confirme que `TRACING_ENABLED=true`
4. Verifique se o endpoint está correto: `OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318`

### Erro ao conectar no OTLP endpoint

- Se a aplicação estiver em container, use o nome do serviço: `jaeger:4318`
- Se estiver rodando localmente, use: `localhost:4318`

## Próximos Passos

- Adicionar tracing para chamadas de banco de dados
- Adicionar tracing para chamadas HTTP externas
- Configurar sampling para produção
- Integrar com outros backends (Tempo, Zipkin, etc.)
