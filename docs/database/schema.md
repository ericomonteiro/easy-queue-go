# EasyQueue Database Schema

This document contains the Entity-Relationship Diagram (ERD) for the EasyQueue system.

## Database Diagram

```mermaid
erDiagram
    USERS ||--o{ BUSINESS_OWNERS : "is"
    USERS ||--o{ CUSTOMERS : "is"
    BUSINESS_OWNERS ||--o{ BUSINESSES : "owns"
    BUSINESSES ||--o{ SERVICES : "offers"
    BUSINESSES ||--o{ QUEUES : "manages"
    SERVICES ||--o{ QUEUE_ENTRIES : "has"
    CUSTOMERS ||--o{ QUEUE_ENTRIES : "joins"
    CUSTOMERS ||--o{ APPOINTMENTS : "books"
    SERVICES ||--o{ APPOINTMENTS : "for"
    QUEUE_ENTRIES ||--o| CHECK_INS : "has"
    CUSTOMERS ||--o{ CHECK_INS : "performs"
    CUSTOMERS ||--o{ CUSTOMER_RATINGS : "receives"
    BUSINESSES ||--o{ BUSINESS_RATINGS : "receives"
    QUEUE_ENTRIES ||--o| NOTIFICATIONS : "triggers"
    APPOINTMENTS ||--o| NOTIFICATIONS : "triggers"
    CUSTOMERS ||--o{ NOTIFICATIONS : "receives"

    USERS {
        uuid id PK
        string email UK
        string password_hash
        string phone
        string role "BO or CU"
        timestamp created_at
        timestamp updated_at
        boolean is_active
    }

    BUSINESS_OWNERS {
        uuid id PK
        uuid user_id FK
        string business_name
        string document_number
        timestamp created_at
        timestamp updated_at
    }

    CUSTOMERS {
        uuid id PK
        uuid user_id FK
        string full_name
        decimal reputation_score
        int total_appointments
        int completed_appointments
        int no_shows
        timestamp created_at
        timestamp updated_at
    }

    BUSINESSES {
        uuid id PK
        uuid owner_id FK
        string name
        string description
        string address
        decimal latitude
        decimal longitude
        decimal min_check_in_distance_km "Minimum distance for check-in"
        int check_in_tolerance_minutes "Time tolerance for check-in"
        string timezone
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }

    SERVICES {
        uuid id PK
        uuid business_id FK
        string name
        string description
        int average_duration_minutes
        decimal price
        boolean is_active
        int max_queue_size
        timestamp created_at
        timestamp updated_at
    }

    QUEUES {
        uuid id PK
        uuid business_id FK
        date queue_date
        string status "OPEN, CLOSED, PAUSED"
        int current_position
        timestamp created_at
        timestamp updated_at
    }

    QUEUE_ENTRIES {
        uuid id PK
        uuid queue_id FK
        uuid customer_id FK
        uuid service_id FK
        int position
        string status "WAITING, NOTIFIED, CHECKED_IN, IN_SERVICE, COMPLETED, CANCELLED, NO_SHOW"
        timestamp estimated_service_time
        timestamp actual_service_time
        timestamp notified_at
        timestamp checked_in_at
        timestamp completed_at
        timestamp cancelled_at
        timestamp created_at
        timestamp updated_at
    }

    APPOINTMENTS {
        uuid id PK
        uuid customer_id FK
        uuid service_id FK
        timestamp scheduled_time
        int duration_minutes
        string status "SCHEDULED, CONFIRMED, CHECKED_IN, IN_SERVICE, COMPLETED, CANCELLED, NO_SHOW"
        timestamp checked_in_at
        timestamp completed_at
        timestamp cancelled_at
        string cancellation_reason
        timestamp created_at
        timestamp updated_at
    }

    CHECK_INS {
        uuid id PK
        uuid customer_id FK
        uuid queue_entry_id FK
        decimal latitude
        decimal longitude
        decimal distance_from_business_km
        timestamp check_in_time
        boolean is_valid
        string validation_message
        timestamp created_at
    }

    NOTIFICATIONS {
        uuid id PK
        uuid customer_id FK
        uuid queue_entry_id FK
        uuid appointment_id FK
        string type "QUEUE_POSITION, TIME_TO_LEAVE, CHECK_IN_REMINDER, SERVICE_READY, COMPLETED, CANCELLED"
        string title
        string message
        string channel "PUSH, SMS, EMAIL"
        boolean is_sent
        timestamp sent_at
        boolean is_read
        timestamp read_at
        timestamp created_at
    }

    CUSTOMER_RATINGS {
        uuid id PK
        uuid customer_id FK
        uuid business_id FK
        uuid queue_entry_id FK
        uuid appointment_id FK
        int rating "1-5"
        string comment
        timestamp created_at
    }

    BUSINESS_RATINGS {
        uuid id PK
        uuid business_id FK
        uuid customer_id FK
        int rating "1-5"
        string comment
        timestamp created_at
    }
```

## Entity Descriptions

### Core Entities

- **USERS**: Base authentication and user management table for both Business Owners and Customers
- **BUSINESS_OWNERS**: Business owner profile extending the Users table
- **CUSTOMERS**: Customer profile with reputation tracking
- **BUSINESSES**: Physical business locations with geolocation data
- **SERVICES**: Services offered by businesses with duration and pricing

### Queue Management

- **QUEUES**: Daily queue instances for each business
- **QUEUE_ENTRIES**: Individual customer positions in queues with status tracking
- **APPOINTMENTS**: Pre-scheduled appointments as an alternative to walk-in queues

### Location & Validation

- **CHECK_INS**: Geolocation-based attendance validation with distance verification

### Communication

- **NOTIFICATIONS**: Multi-channel notification system for queue updates and reminders

### Reputation System

- **CUSTOMER_RATINGS**: Business ratings of customer reliability
- **BUSINESS_RATINGS**: Customer ratings of business service quality

## Key Features Supported

1. **Geolocation Check-in**: `CHECK_INS` table validates customer location against business coordinates
2. **Distance & Tolerance Rules**: `BUSINESSES` table stores configurable `min_check_in_distance_km` and `check_in_tolerance_minutes`
3. **Real-time Queue Tracking**: `QUEUE_ENTRIES` with position and status management
4. **Smart Notifications**: `NOTIFICATIONS` table supports multiple channels and types
5. **Reputation System**: Both `CUSTOMER_RATINGS` and `BUSINESS_RATINGS` for accountability
6. **Hybrid Scheduling**: Supports both walk-in queues and pre-scheduled appointments

## Indexes Recommendations

```sql
-- User lookups
CREATE INDEX idx_users_email ON USERS(email);
CREATE INDEX idx_users_role ON USERS(role);

-- Business queries
CREATE INDEX idx_businesses_owner_id ON BUSINESSES(owner_id);
CREATE INDEX idx_businesses_location ON BUSINESSES(latitude, longitude);
CREATE INDEX idx_businesses_is_active ON BUSINESSES(is_active);

-- Service queries
CREATE INDEX idx_services_business_id ON SERVICES(business_id);
CREATE INDEX idx_services_is_active ON SERVICES(is_active);

-- Queue management
CREATE INDEX idx_queues_business_date ON QUEUES(business_id, queue_date);
CREATE INDEX idx_queue_entries_queue_id ON QUEUE_ENTRIES(queue_id);
CREATE INDEX idx_queue_entries_customer_id ON QUEUE_ENTRIES(customer_id);
CREATE INDEX idx_queue_entries_status ON QUEUE_ENTRIES(status);
CREATE INDEX idx_queue_entries_position ON QUEUE_ENTRIES(queue_id, position);

-- Appointments
CREATE INDEX idx_appointments_customer_id ON APPOINTMENTS(customer_id);
CREATE INDEX idx_appointments_service_id ON APPOINTMENTS(service_id);
CREATE INDEX idx_appointments_scheduled_time ON APPOINTMENTS(scheduled_time);
CREATE INDEX idx_appointments_status ON APPOINTMENTS(status);

-- Check-ins
CREATE INDEX idx_check_ins_customer_id ON CHECK_INS(customer_id);
CREATE INDEX idx_check_ins_queue_entry_id ON CHECK_INS(queue_entry_id);

-- Notifications
CREATE INDEX idx_notifications_customer_id ON NOTIFICATIONS(customer_id);
CREATE INDEX idx_notifications_is_sent ON NOTIFICATIONS(is_sent);
CREATE INDEX idx_notifications_created_at ON NOTIFICATIONS(created_at);

-- Ratings
CREATE INDEX idx_customer_ratings_customer_id ON CUSTOMER_RATINGS(customer_id);
CREATE INDEX idx_customer_ratings_business_id ON CUSTOMER_RATINGS(business_id);
CREATE INDEX idx_business_ratings_business_id ON BUSINESS_RATINGS(business_id);
```
