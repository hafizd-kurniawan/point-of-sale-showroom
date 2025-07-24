import 'package:dio/dio.dart';
import '../models/product_model.dart';
import '../models/api_response.dart';
import '../../core/network/network_service.dart';
import '../../core/constants/app_constants.dart';

abstract class ProductRemoteDataSource {
  Future<List<ProductModel>> getProducts({
    int page = 1,
    int limit = 20,
    String? search,
  });
  Future<ProductModel> getProduct(int productId);
  Future<List<ProductModel>> getLowStockProducts();
  Future<List<StockMovementModel>> getProductStockMovements(int productId);
  Future<int> getCurrentStock(int productId);
  Future<List<PurchaseOrderModel>> getPurchaseOrders({
    int page = 1,
    int limit = 20,
    String? status,
  });
  Future<PurchaseOrderModel> getPurchaseOrder(int poId);
}

class ProductRemoteDataSourceImpl implements ProductRemoteDataSource {
  final NetworkService _networkService;

  ProductRemoteDataSourceImpl(this._networkService);

  @override
  Future<List<ProductModel>> getProducts({
    int page = 1,
    int limit = 20,
    String? search,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'limit': limit,
      };
      
      if (search != null && search.isNotEmpty) {
        queryParams['search'] = search;
      }

      final response = await _networkService.get(
        ApiConstants.products,
        queryParameters: queryParams,
      );

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          final data = apiResponse.data!['data'] as List<dynamic>;
          return data.map((json) => ProductModel.fromJson(json)).toList();
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get products');
        }
      } else {
        throw Exception('Failed to get products');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get products');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<ProductModel> getProduct(int productId) async {
    try {
      final response = await _networkService.get('${ApiConstants.products}/$productId');

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          return ProductModel.fromJson(apiResponse.data!);
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get product');
        }
      } else {
        throw Exception('Failed to get product');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get product');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<List<ProductModel>> getLowStockProducts() async {
    try {
      final response = await _networkService.get(ApiConstants.lowStockProducts);

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          final data = apiResponse.data!['data'] as List<dynamic>;
          return data.map((json) => ProductModel.fromJson(json)).toList();
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get low stock products');
        }
      } else {
        throw Exception('Failed to get low stock products');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get low stock products');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<List<StockMovementModel>> getProductStockMovements(int productId) async {
    try {
      final response = await _networkService.get(
        '${ApiConstants.products}/$productId/stock-movements',
      );

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          final data = apiResponse.data!['data'] as List<dynamic>;
          return data.map((json) => StockMovementModel.fromJson(json)).toList();
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get stock movements');
        }
      } else {
        throw Exception('Failed to get stock movements');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get stock movements');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<int> getCurrentStock(int productId) async {
    try {
      final response = await _networkService.get(
        '${ApiConstants.products}/$productId/current-stock',
      );

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          return apiResponse.data!['stock_quantity'] as int;
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get current stock');
        }
      } else {
        throw Exception('Failed to get current stock');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get current stock');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<List<PurchaseOrderModel>> getPurchaseOrders({
    int page = 1,
    int limit = 20,
    String? status,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'limit': limit,
      };
      
      if (status != null && status.isNotEmpty) {
        queryParams['status'] = status;
      }

      final response = await _networkService.get(
        ApiConstants.purchaseOrders,
        queryParameters: queryParams,
      );

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          final data = apiResponse.data!['data'] as List<dynamic>;
          return data.map((json) => PurchaseOrderModel.fromJson(json)).toList();
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get purchase orders');
        }
      } else {
        throw Exception('Failed to get purchase orders');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get purchase orders');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<PurchaseOrderModel> getPurchaseOrder(int poId) async {
    try {
      final response = await _networkService.get('${ApiConstants.purchaseOrders}/$poId');

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          return PurchaseOrderModel.fromJson(apiResponse.data!);
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get purchase order');
        }
      } else {
        throw Exception('Failed to get purchase order');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get purchase order');
      }
      throw Exception('Network error occurred');
    }
  }
}