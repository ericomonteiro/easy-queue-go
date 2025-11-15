# Tracing with OpenTelemetry

This project uses OpenTelemetry for distributed tracing, with Jaeger as the visualization backend.

## 📋 Table of Contents

- [Architecture](#architecture)
- [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
- [How to Use](#how-to-use)
  - [1. Start Services](#1-start-services)
  - [2. Access Jaeger UI](#2-access-jaeger-ui)
  - [3. Run the Application](#3-run-the-application)
  - [4. Generate Traces](#4-generate-traces)
  - [5. View Traces in Jaeger](#5-view-traces-in-jaeger)
- [What is Automatically Traced](#what-is-automatically-traced)
- [Adding Custom Spans](#adding-custom-spans)
  - [Basic Example](#basic-example)
  - [Example with Error Handling](#example-with-error-handling)
  - [Example with Nested Spans](#example-with-nested-spans)
  - [Example in HTTP Handler](#example-in-http-handler)
  - [Supported Attribute Types](#supported-attribute-types)
- [Disable Tracing](#disable-tracing)
- [Ports Used](#ports-used)
- [Troubleshooting](#troubleshooting)
  - [Traces don't appear in Jaeger](#traces-dont-appear-in-jaeger)
  - [Error connecting to OTLP endpoint](#error-connecting-to-otlp-endpoint)
- [Next Steps](#next-steps)

---

## Architecture

- **Instrumentation**: OpenTelemetry SDK
- **Protocol**: OTLP (OpenTelemetry Protocol) via HTTP
- **Backend**: Jaeger All-in-One
- **Middleware**: otelgin for automatic Gin instrumentation

## Configuration

### Environment Variables

```bash
# Enable/disable tracing
TRACING_ENABLED=true

# Service name (appears in Jaeger)
SERVICE_NAME=easy-queue-go

# Service version
SERVICE_VERSION=1.0.0

# Environment (development, staging, production)
ENVIRONMENT=development

# OTLP collector endpoint (without http://)
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
```

## How to Use

### 1. Start Services

```bash
# Start Jaeger and PostgreSQL
docker-compose up -d

# Check if services are running
docker-compose ps
```

### 2. Access Jaeger UI

Open your browser at: http://localhost:16686

### 3. Run the Application

```bash
# Make sure you have a .env file with configurations
cp .env.example .env

# Run the application
go run src/internal/cmd/main.go
```

### 4. Generate Traces

Make requests to the application:

```bash
# Health check
curl http://localhost:8080/health

# Other routes...
```

### 5. View Traces in Jaeger

1. Access http://localhost:16686
2. Select the "easy-queue-go" service in the dropdown
3. Click "Find Traces"
4. Explore the generated traces

## What is Automatically Traced

The `otelgin` middleware automatically traces:

- **HTTP Requests**: Method, path, status code
- **Timing**: Total request duration
- **Errors**: Errors and stack traces
- **Context Propagation**: Trace context propagation between services

## Adding Custom Spans

### Basic Example

To add custom spans in your code:

```go
import (
    "context"
    "easy-queue-go/src/internal/tracing"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
)

func MyFunction(ctx context.Context) error {
    // Create a new span
    tracer := tracing.Tracer("my-component")
    ctx, span := tracer.Start(ctx, "MyFunction")
    defer span.End()

    // Add attributes to the span
    span.SetAttributes(
        attribute.String("user.id", "123"),
        attribute.Int("items.count", 42),
        attribute.Bool("is.success", true),
    )

    // Your code here...
    
    return nil
}
```

### Example with Error Handling

```go
func ProcessOrder(ctx context.Context, orderID string) error {
    tracer := tracing.Tracer("orders")
    ctx, span := tracer.Start(ctx, "ProcessOrder")
    defer span.End()

    span.SetAttributes(
        attribute.String("order.id", orderID),
        attribute.String("operation.type", "process"),
    )

    // Simulate processing
    err := performProcessing(orderID)
    if err != nil {
        // Record error in span
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return err
    }

    // Mark span as successful
    span.SetStatus(codes.Ok, "Order processed successfully")
    return nil
}
```

### Example with Nested Spans

```go
func CompleteOperation(ctx context.Context) error {
    tracer := tracing.Tracer("example")

    // Parent span
    ctx, parentSpan := tracer.Start(ctx, "CompleteOperation")
    defer parentSpan.End()

    // First operation (child span)
    ctx, span1 := tracer.Start(ctx, "FetchData")
    span1.SetAttributes(attribute.String("source", "database"))
    // ... fetch data ...
    span1.End()

    // Second operation (child span)
    ctx, span2 := tracer.Start(ctx, "ProcessData")
    span2.SetAttributes(attribute.Int("records.count", 100))
    // ... process data ...
    span2.End()

    // Third operation (child span)
    ctx, span3 := tracer.Start(ctx, "SaveResult")
    span3.SetAttributes(attribute.String("destination", "cache"))
    // ... save result ...
    span3.End()

    parentSpan.SetStatus(codes.Ok, "Complete operation finished")
    return nil
}
```

### Example in HTTP Handler

```go
func MyHandler(c *gin.Context) {
    // Context already comes with span from otelgin middleware
    ctx := c.Request.Context()
    
    tracer := tracing.Tracer("handlers")
    ctx, span := tracer.Start(ctx, "ProcessRequest")
    defer span.End()

    // Add request information
    span.SetAttributes(
        attribute.String("user.agent", c.GetHeader("User-Agent")),
        attribute.String("request.id", c.GetString("request_id")),
    )

    // Process the request
    result, err := processLogic(ctx)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "Error processing")
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    span.SetStatus(codes.Ok, "Request processed")
    c.JSON(200, result)
}
```

### Supported Attribute Types

```go
span.SetAttributes(
    // Strings
    attribute.String("key", "value"),
    
    // Integers
    attribute.Int("count", 42),
    attribute.Int64("timestamp", time.Now().Unix()),
    
    // Decimals
    attribute.Float64("price", 19.99),
    
    // Booleans
    attribute.Bool("is_active", true),
    
    // Arrays
    attribute.StringSlice("tags", []string{"tag1", "tag2"}),
    attribute.IntSlice("ids", []int{1, 2, 3}),
)
```

## Disable Tracing

To disable tracing (useful in tests or development):

```bash
TRACING_ENABLED=false
```

## Ports Used

- **16686**: Jaeger UI
- **4318**: OTLP HTTP receiver
- **4317**: OTLP gRPC receiver

## Troubleshooting

### Traces don't appear in Jaeger

1. Check if Jaeger is running: `docker-compose ps`
2. Check Jaeger logs: `docker-compose logs jaeger`
3. Confirm that `TRACING_ENABLED=true`
4. Verify the endpoint is correct: `OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318`

### Error connecting to OTLP endpoint

- If the application is in a container, use the service name: `jaeger:4318`
- If running locally, use: `localhost:4318`

## Next Steps

- Add tracing for database calls
- Add tracing for external HTTP calls
- Configure sampling for production
- Integrate with other backends (Tempo, Zipkin, etc.)
