# 🔄 Migrações de Banco de Dados

Esta página documenta as migrações do banco de dados do EasyQueue.

## 📋 Visão Geral

As migrações de banco de dados são usadas para versionar e aplicar mudanças no schema do banco de dados de forma controlada e rastreável.

## 🛠️ Ferramenta de Migração

O EasyQueue utiliza [golang-migrate](https://github.com/golang-migrate/migrate) para gerenciar migrações.

### Instalação

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Windows
scoop install migrate
```

## 📁 Estrutura de Migrações

```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_businesses_table.up.sql
├── 000002_create_businesses_table.down.sql
└── ...
```

## 🚀 Comandos Úteis

### Criar uma Nova Migração

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

### Aplicar Migrações

```bash
# Aplicar todas as migrações pendentes
migrate -path migrations -database "postgresql://easyqueue:easyqueue123@localhost:5432/easyqueue?sslmode=disable" up

# Aplicar N migrações
migrate -path migrations -database "postgresql://..." up 2
```

### Reverter Migrações

```bash
# Reverter última migração
migrate -path migrations -database "postgresql://..." down 1

# Reverter todas as migrações
migrate -path migrations -database "postgresql://..." down
```

### Verificar Versão Atual

```bash
migrate -path migrations -database "postgresql://..." version
```

### Forçar Versão (em caso de erro)

```bash
migrate -path migrations -database "postgresql://..." force 1
```

## 📝 Migrações Planejadas

### Migration 001: Tabelas de Usuários

**Arquivo:** `000001_create_users_tables.up.sql`

```sql
-- Tabela de usuários base
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(10) NOT NULL CHECK (role IN ('BO', 'CU')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Business Owners
CREATE TABLE business_owners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    business_name VARCHAR(255),
    document_number VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Customers
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    reputation_score DECIMAL(3,2) DEFAULT 5.00,
    total_appointments INTEGER DEFAULT 0,
    completed_appointments INTEGER DEFAULT 0,
    no_shows INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_business_owners_user_id ON business_owners(user_id);
CREATE INDEX idx_customers_user_id ON customers(user_id);
```

**Arquivo:** `000001_create_users_tables.down.sql`

```sql
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS business_owners;
DROP TABLE IF EXISTS users;
```

### Migration 002: Tabelas de Negócios e Serviços

**Arquivo:** `000002_create_businesses_services.up.sql`

```sql
-- Tabela de Businesses
CREATE TABLE businesses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES business_owners(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address TEXT NOT NULL,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    min_check_in_distance_km DECIMAL(5,2) DEFAULT 2.00,
    check_in_tolerance_minutes INTEGER DEFAULT 15,
    timezone VARCHAR(50) DEFAULT 'America/Sao_Paulo',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Services
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    average_duration_minutes INTEGER NOT NULL,
    price DECIMAL(10,2),
    is_active BOOLEAN DEFAULT true,
    max_queue_size INTEGER DEFAULT 50,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX idx_businesses_owner_id ON businesses(owner_id);
CREATE INDEX idx_businesses_location ON businesses(latitude, longitude);
CREATE INDEX idx_businesses_is_active ON businesses(is_active);
CREATE INDEX idx_services_business_id ON services(business_id);
CREATE INDEX idx_services_is_active ON services(is_active);
```

### Migration 003: Tabelas de Filas

**Arquivo:** `000003_create_queues.up.sql`

```sql
-- Tabela de Queues
CREATE TABLE queues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    queue_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED', 'PAUSED')),
    current_position INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(business_id, queue_date)
);

-- Tabela de Queue Entries
CREATE TABLE queue_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    queue_id UUID NOT NULL REFERENCES queues(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES services(id),
    position INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'WAITING' CHECK (status IN ('WAITING', 'NOTIFIED', 'CHECKED_IN', 'IN_SERVICE', 'COMPLETED', 'CANCELLED', 'NO_SHOW')),
    estimated_service_time TIMESTAMP WITH TIME ZONE,
    actual_service_time TIMESTAMP WITH TIME ZONE,
    notified_at TIMESTAMP WITH TIME ZONE,
    checked_in_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX idx_queues_business_date ON queues(business_id, queue_date);
CREATE INDEX idx_queue_entries_queue_id ON queue_entries(queue_id);
CREATE INDEX idx_queue_entries_customer_id ON queue_entries(customer_id);
CREATE INDEX idx_queue_entries_status ON queue_entries(status);
CREATE INDEX idx_queue_entries_position ON queue_entries(queue_id, position);
```

## 🔒 Boas Práticas

1. **Sempre crie migrações reversíveis** - Cada `.up.sql` deve ter um `.down.sql` correspondente
2. **Teste as migrações** - Teste tanto `up` quanto `down` em ambiente de desenvolvimento
3. **Nunca modifique migrações aplicadas** - Crie uma nova migração para correções
4. **Use transações** - Envolva mudanças complexas em transações
5. **Documente mudanças** - Adicione comentários explicando o propósito da migração

## 📚 Referências

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [Database Migration Best Practices](https://www.prisma.io/dataguide/types/relational/migration-strategies)
- [Schema do Banco de Dados](schema.md)

---

**Status:** 🚧 Em desenvolvimento - Migrações serão implementadas na próxima fase do projeto.
