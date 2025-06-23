# Coding Best Practices

## Organization Structure

This project follows a layered (hexagonal like) architecture with clear separation of concerns:

- **internal/app** - Contains core application logic and business rules
  - Defines domain models and interfaces
  - Implements business processes and workflows
  - Orchestrates interactions between different parts of the system
  - Should be independent of infrastructure details

- **internal/services** - Contains infrastructure-related implementations
  - Provides concrete implementations of interfaces defined in the app layer
  - Handles external system interactions (HTTP clients, databases, etc.)
  - Contains adapters that translate between domain and external models
  - Deals with technical concerns (logging, caching, etc.)

The layered structure helps maintain separation of concerns while still allowing practical code organization. While we strive for clean boundaries between layers, we allow data structures to be shared when it makes the codebase more maintainable.

## Core Principles

### Pragmatic Architecture
- Focus on solving business problems over architectural purity
- Choose simplicity over rigid adherence to patterns
- Favor readability and maintainability over complex abstractions
- Apply architectural boundaries where they provide clear value

### Use Interfaces Strategically
- **Follow "accept interfaces, return concrete types"** - Classic Go principle for flexible consumption
- **Avoid interface overuse** - Don't create interfaces for every service or component
- **Interface when needed** - Add interfaces when you need multiple implementations or for testing

### Efficient Data Flow
- **Share data structures when sensible** - Data models can cross architectural boundaries
- **Keep business logic in application layer** - Core decisions belong in domain/application code
- **Separate behavior from data** - Methods that change behavior belong in specific layers
- **Translation at boundaries** - Convert between formats at integration points when necessary

### Go Idiomatic Style
- **Composition over inheritance** - Use struct embedding rather than class hierarchies
- **Error handling as values** - Follow Go's error handling patterns consistently
- **Flat package structure** - Keep package structure flat and focused on domain concepts. Split into subpackages only when necessary.

## Patterns to Follow

### Common principles

- Keep package structure flat and focused on domain concepts
- Prefer small interfaces with clear purposes
- Use dependency injection for testability and flexibility
- Balance between business domain modeling and technical architecture

### Layer Data Models Appropriately

Choose the right approach for data sharing based on complexity:

```go
// ACCEPTABLE: Direct use of infrastructure level service models in simple applications
import "github.com/myapp/internal/services/payments"

func (orders *OrdersService) ProcessOrder(order Order) error {
    // Using payments.Invoice directly when it's just a data container
    invoice := payments.Invoice{
        CustomerID: order.CustomerID,
        Items:      mapOrderItemsToInvoiceItems(order.Items),
    }

    return orders.paymentsClient.CreateInvoice(ctx, invoice)
}

// BETTER FOR COMPLEX DOMAINS: Define your own models with translation
import "github.com/myapp/internal/services/payments"

// App layer owns its models
type Invoice struct {
    Customer string
    Items    []InvoiceItem
    Total    Money
}

// Translation happens at boundary
func convertToServiceInvoice(inv Invoice) payments.Invoice {
    // mapping logic
}
```

### Use Dependency Injection

Structure your components for clean dependency management:

```go
// Explicitly declare dependencies
type OrderService struct {
    repository   OrderRepository
    payments      PaymentsService
    notification NotificationSender
    logger       *slog.Logger
}

// Group dependencies with struct
type OrderServiceDeps struct {
    dig.In
    
    Repository   OrderRepository
    Payments      PaymentsService
    Notification NotificationSender
    RootLogger   *slog.Logger
}

// Constructor that clearly shows dependencies
func NewOrderService(deps OrderServiceDeps) *OrderService {
    return &OrderService{
        repository:   deps.Repository,
        payments:      deps.Payments,
        notification: deps.Notification,
        logger:       deps.RootLogger.WithGroup("order-service"),
    }
}
```

## Common Anti-patterns to Avoid

- **Interface Explosion** - Creating interfaces for everything "just in case"
- **Premature Optimization** - Overcomplicating code for theoretical performance or architectural gains

## Remember

> "Clear is better than clever." - Rob Pike

Prioritize code that is easy to understand, maintain, and extend.