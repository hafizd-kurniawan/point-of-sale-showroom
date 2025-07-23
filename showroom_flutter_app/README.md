# Showroom POS Flutter App

A Flutter application for Point of Sale management specifically designed for Vehicle Showroom operations. This app supports **Cashier** and **Mechanic** roles only, providing a clean and intuitive interface for their daily operations.

## Features

### For Cashiers
- Vehicle sales processing
- Customer information management
- Transaction history
- Sales reports
- Vehicle inventory information

### For Mechanics
- Repair work queue management
- Spare parts requests
- Work history tracking
- Time tracking for repairs
- Parts inventory information

## Architecture

This project follows **Clean Architecture** principles with:

- **Domain Layer**: Business entities, use cases, and repository interfaces
- **Data Layer**: API integration, data models, and repository implementations
- **Presentation Layer**: UI components, BLoC state management, and routing

### Technologies Used

- **Flutter** for cross-platform mobile development
- **BLoC** for state management
- **Dio** for HTTP networking
- **Go Router** for navigation
- **SharedPreferences** for local storage
- **Get It** for dependency injection
- **JSON Annotation** for serialization

## Project Structure

```
lib/
├── core/
│   ├── constants/     # App constants and API endpoints
│   ├── di/           # Dependency injection setup
│   ├── network/      # HTTP client configuration
│   ├── router/       # App routing configuration
│   ├── storage/      # Local storage service
│   ├── theme/        # App theme and styling
│   └── utils/        # Utility functions
├── data/
│   ├── datasources/  # Remote data sources
│   ├── models/       # Data transfer objects
│   └── repositories/ # Repository implementations
├── domain/
│   ├── entities/     # Business entities
│   ├── repositories/ # Repository interfaces
│   └── usecases/     # Business use cases
└── presentation/
    ├── bloc/         # BLoC state management
    ├── pages/        # UI screens
    └── widgets/      # Reusable UI components
```

## Backend Integration

The app integrates with a Go backend API that provides:

- Authentication endpoints (`/auth/*`)
- Product management (`/admin/products/*`)
- Stock management (`/admin/stock-*`)
- Purchase orders (`/admin/purchase-orders/*`)
- Customer management (`/admin/customers/*`)
- Vehicle management (`/admin/vehicle-*`)

## Role-Based Access

- **Admin features are NOT implemented** in the mobile app - only informational displays
- **Cashier role**: Access to sales, customer, and transaction features
- **Mechanic role**: Access to repair, parts, and work order features

## Getting Started

1. Ensure Flutter is installed on your development machine
2. Clone the repository
3. Run `flutter pub get` to install dependencies
4. Configure the backend API URL in `lib/core/constants/app_constants.dart`
5. Run `flutter run` to start the application

## Configuration

Update the API base URL in `lib/core/constants/app_constants.dart`:

```dart
static const String baseUrl = 'http://your-backend-url:8080/api/v1';
```

## Authentication

The app supports JWT-based authentication with automatic token refresh. Users must login with valid credentials from the backend system.

## Note

This is a mobile-specific implementation focusing on cashier and mechanic workflows. Admin functionality is accessible only through the web interface of the backend system.