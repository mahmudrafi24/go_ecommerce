# Implementation Plan

- [x] 1. Set up project structure and dependencies





  - Initialize Go module with `go mod init`
  - Install dependencies: chi, pgx, golang-jwt/jwt, bcrypt
  - Create directory structure: cmd/server, internal/{config,models,repository,service,handler,middleware,database}
  - Create main.go entry point with basic server setup
  - _Requirements: 7.1_

- [x] 2. Implement database layer and migrations





  - [x] 2.1 Create PostgreSQL connection manager


    - Implement database/postgres.go with connection pool
    - Add configuration for database URL from environment
    - _Requirements: 6.1, 6.2_


  - [x] 2.2 Create database migration file





    - Write migrations/001_initial.sql with all tables (users, products, cart_items, orders, order_items)
    - Include indexes and foreign key constraints
    - _Requirements: 6.1_

- [x] 3. Implement models and validation





  - [x] 3.1 Create User model with validation


    - Implement models/user.go with User struct and validation methods
    - Add email format validation, password length check
    - _Requirements: 1.1, 1.2_
  - [x] 3.2 Write property test for user validation






    - **Property 2: Invalid registration data is rejected**
    - **Validates: Requirements 1.2**



  - [x] 3.3 Create Product model with validation




    - Implement models/product.go with Product struct


    - Add validation for name, price >= 0, stock >= 0
    - _Requirements: 2.1, 3.1_
  - [x] 3.4 Create Cart and CartItem models





    - Implement models/cart.go with Cart, CartItem structs
    - Add total calculation method


    - _Requirements: 4.4_
  - [x] 3.5 Write property test for cart total calculation






    - **Property 10: Cart total equals sum of item prices**
    - **Validates: Requirements 4.4**
  - [x] 3.6 Create Order and OrderItem models





    - Implement models/order.go with Order, OrderItem structs
    - _Requirements: 5.1, 5.3_

- [x] 4. Implement repositories





  - [x] 4.1 Create User repository


    - Implement repository/user_repo.go with Create, GetByID, GetByEmail
    - _Requirements: 1.1, 1.3_
  - [ ]* 4.2 Write property test for user persistence round-trip
    - **Property 14: Data persistence round-trip (User)**
    - **Validates: Requirements 6.1, 6.2**

  - [x] 4.3 Create Product repository

    - Implement repository/product_repo.go with CRUD, List, Search, UpdateStock
    - _Requirements: 2.1, 2.2, 3.1, 3.2, 3.3_
  - [ ]* 4.4 Write property test for product persistence round-trip
    - **Property 14: Data persistence round-trip (Product)**
    - **Validates: Requirements 6.1, 6.2**
  - [ ]* 4.5 Write property test for product search
    - **Property 6: Search results match search term**
    - **Validates: Requirements 2.2**

  - [x] 4.6 Create Cart repository

    - Implement repository/cart_repo.go with GetByUserID, AddItem, UpdateItemQuantity, RemoveItem, Clear
    - _Requirements: 4.1, 4.2, 4.3, 4.4_

  - [x] 4.7 Create Order repository

    - Implement repository/order_repo.go with Create, GetByID, ListByUserID
    - Use transactions for order creation
    - _Requirements: 5.1, 5.3, 5.4, 6.4_
  - [-] 4.8 Write property test for order persistence round-trip




    - **Property 14: Data persistence round-trip (Order)**
    - **Validates: Requirements 6.1, 6.2**

- [ ] 5. Checkpoint
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 6. Implement authentication service
  - [ ] 6.1 Create auth service with JWT handling
    - Implement service/auth_service.go with Register, Login, ValidateToken
    - Use bcrypt for password hashing
    - Generate and validate JWT tokens
    - _Requirements: 1.1, 1.3, 1.5_
  - [ ]* 6.2 Write property test for registration creates retrievable user
    - **Property 1: User registration creates retrievable user**
    - **Validates: Requirements 1.1**
  - [ ]* 6.3 Write property test for valid credentials produce valid JWT
    - **Property 3: Valid credentials produce valid JWT**
    - **Validates: Requirements 1.3**
  - [ ]* 6.4 Write property test for invalid credentials rejection
    - **Property 4: Invalid credentials are rejected**
    - **Validates: Requirements 1.4**

- [ ] 7. Implement product service
  - [ ] 7.1 Create product service
    - Implement service/product_service.go with Create, GetByID, List, Search, Update, Delete
    - Add authorization checks for admin operations
    - _Requirements: 2.1, 2.2, 2.3, 3.1, 3.2, 3.3, 3.4_
  - [ ]* 7.2 Write property test for product CRUD round-trip
    - **Property 7: Product CRUD round-trip**
    - **Validates: Requirements 3.1, 3.2**
  - [ ]* 7.3 Write property test for non-admin authorization
    - **Property 8: Non-admin users cannot manage products**
    - **Validates: Requirements 3.4**
  - [ ]* 7.4 Write property test for product list fields
    - **Property 5: Product list contains required fields**
    - **Validates: Requirements 2.1**

- [ ] 8. Implement cart service
  - [ ] 8.1 Create cart service
    - Implement service/cart_service.go with GetCart, AddItem, UpdateQuantity, RemoveItem
    - Add stock validation before adding items
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_
  - [ ]* 8.2 Write property test for adding to cart increases size
    - **Property 9: Adding to cart increases cart size**
    - **Validates: Requirements 4.1**
  - [ ]* 8.3 Write property test for insufficient stock rejection
    - **Property 11: Insufficient stock prevents cart addition**
    - **Validates: Requirements 4.5**

- [ ] 9. Implement order service
  - [ ] 9.1 Create order service
    - Implement service/order_service.go with CreateOrder, GetOrder, ListOrders
    - Use transactions for order creation with stock reduction and cart clearing
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_
  - [ ]* 9.2 Write property test for order creation effects
    - **Property 12: Order creation reduces stock and clears cart**
    - **Validates: Requirements 5.1**
  - [ ]* 9.3 Write property test for order history ownership
    - **Property 13: Order history returns user's orders only**
    - **Validates: Requirements 5.3**
  - [ ]* 9.4 Write property test for transaction rollback
    - **Property 15: Transaction rollback on failure**
    - **Validates: Requirements 6.4**

- [ ] 10. Checkpoint
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 11. Implement HTTP handlers and middleware
  - [ ] 11.1 Create JWT authentication middleware
    - Implement middleware/auth.go with JWT validation
    - Extract user ID and role from token
    - _Requirements: 1.5_
  - [ ] 11.2 Create error response helpers
    - Implement handler/response.go with JSON response helpers
    - Map AppError to HTTP status codes
    - _Requirements: 7.1, 7.2, 7.3_
  - [ ]* 11.3 Write property test for API response format
    - **Property 16: API responses are valid JSON**
    - **Validates: Requirements 7.1, 7.2, 7.3**
  - [ ] 11.4 Create auth handlers
    - Implement handler/auth_handler.go with Register, Login endpoints
    - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - [ ] 11.5 Create product handlers
    - Implement handler/product_handler.go with CRUD endpoints
    - Add admin middleware for management endpoints
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 3.1, 3.2, 3.3, 3.4_
  - [ ] 11.6 Create cart handlers
    - Implement handler/cart_handler.go with cart operation endpoints
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_
  - [ ] 11.7 Create order handlers
    - Implement handler/order_handler.go with order endpoints
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [ ] 12. Wire up router and complete server
  - [ ] 12.1 Create router configuration
    - Implement cmd/server/routes.go with all route definitions
    - Apply middleware to protected routes
    - _Requirements: 1.5, 3.4_
  - [ ] 12.2 Complete main.go with full initialization
    - Initialize database connection
    - Create all repositories, services, handlers
    - Start HTTP server
    - _Requirements: 6.1_

- [ ] 13. Final Checkpoint
  - Ensure all tests pass, ask the user if questions arise.
