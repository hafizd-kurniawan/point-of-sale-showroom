class ApiConstants {
  static const String baseUrl = 'http://localhost:8080/api/v1';
  
  // Auth endpoints
  static const String login = '/auth/login';
  static const String logout = '/auth/logout';
  static const String me = '/auth/me';
  static const String profile = '/auth/profile';
  static const String changePassword = '/auth/change-password';
  static const String refreshToken = '/auth/refresh';
  
  // Product endpoints
  static const String products = '/admin/products';
  static const String productCategories = '/admin/product-categories';
  static const String lowStockProducts = '/admin/products/low-stock';
  
  // Purchase Order endpoints
  static const String purchaseOrders = '/admin/purchase-orders';
  static const String poDetails = '/admin/purchase-order-details';
  static const String pendingApproval = '/admin/purchase-orders/pending-approval';
  
  // Stock Management endpoints
  static const String stockMovements = '/admin/stock-movements';
  static const String stockAdjustments = '/admin/stock-adjustments';
  static const String goodsReceipts = '/admin/goods-receipts';
  
  // Supplier endpoints
  static const String suppliers = '/admin/suppliers';
  static const String supplierPayments = '/admin/supplier-payments';
  
  // Customer endpoints
  static const String customers = '/admin/customers';
  
  // Vehicle endpoints
  static const String vehicleBrands = '/admin/vehicle-brands';
  static const String vehicleCategories = '/admin/vehicle-categories';
  static const String vehicleModels = '/admin/vehicle-models';
}

class AppConstants {
  static const String appName = 'Showroom POS';
  static const String appVersion = '1.0.0';
  
  // Storage keys
  static const String accessTokenKey = 'access_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String userDataKey = 'user_data';
  static const String isLoggedInKey = 'is_logged_in';
  
  // User roles
  static const String adminRole = 'admin';
  static const String cashierRole = 'cashier';
  static const String mechanicRole = 'mechanic';
  static const String salesRole = 'sales';
  
  // Pagination
  static const int defaultPageSize = 20;
  static const int maxPageSize = 100;
  
  // Timeouts
  static const int connectTimeout = 30000; // 30 seconds
  static const int receiveTimeout = 30000; // 30 seconds
}

class UIConstants {
  // Spacing
  static const double spacing4 = 4.0;
  static const double spacing8 = 8.0;
  static const double spacing12 = 12.0;
  static const double spacing16 = 16.0;
  static const double spacing20 = 20.0;
  static const double spacing24 = 24.0;
  static const double spacing32 = 32.0;
  
  // Border radius
  static const double radiusSmall = 4.0;
  static const double radiusMedium = 8.0;
  static const double radiusLarge = 12.0;
  static const double radiusXLarge = 16.0;
  
  // Icon sizes
  static const double iconSmall = 16.0;
  static const double iconMedium = 24.0;
  static const double iconLarge = 32.0;
  
  // Button heights
  static const double buttonHeightSmall = 32.0;
  static const double buttonHeightMedium = 40.0;
  static const double buttonHeightLarge = 48.0;
}