import 'package:dio/dio.dart';
import '../models/auth_model.dart';
import '../models/user_model.dart';
import '../models/api_response.dart';
import '../../core/network/network_service.dart';
import '../../core/constants/app_constants.dart';

abstract class AuthRemoteDataSource {
  Future<LoginResponse> login(LoginRequest request);
  Future<void> logout();
  Future<UserInfoModel> getCurrentUser();
  Future<ProfileResponse> getProfile();
  Future<void> changePassword(ChangePasswordRequest request);
  Future<LoginResponse> refreshToken();
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final NetworkService _networkService;

  AuthRemoteDataSourceImpl(this._networkService);

  @override
  Future<LoginResponse> login(LoginRequest request) async {
    try {
      final response = await _networkService.post(
        ApiConstants.login,
        data: request.toJson(),
      );

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          return LoginResponse.fromJson(apiResponse.data!);
        } else {
          throw Exception(apiResponse.error ?? 'Login failed');
        }
      } else {
        throw Exception('Login failed');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Login failed');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<void> logout() async {
    try {
      await _networkService.post(ApiConstants.logout);
    } on DioException catch (e) {
      // Even if logout fails on server, we should clear local data
      if (e.response?.statusCode != 401) {
        throw Exception('Logout failed');
      }
    }
  }

  @override
  Future<UserInfoModel> getCurrentUser() async {
    try {
      final response = await _networkService.get(ApiConstants.me);

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          return UserInfoModel.fromJson(apiResponse.data!);
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get user info');
        }
      } else {
        throw Exception('Failed to get user info');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get user info');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<ProfileResponse> getProfile() async {
    try {
      final response = await _networkService.get(ApiConstants.profile);

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          return ProfileResponse.fromJson(apiResponse.data!);
        } else {
          throw Exception(apiResponse.error ?? 'Failed to get profile');
        }
      } else {
        throw Exception('Failed to get profile');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to get profile');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<void> changePassword(ChangePasswordRequest request) async {
    try {
      final response = await _networkService.post(
        ApiConstants.changePassword,
        data: request.toJson(),
      );

      if (response.statusCode != 200) {
        throw Exception('Failed to change password');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Failed to change password');
      }
      throw Exception('Network error occurred');
    }
  }

  @override
  Future<LoginResponse> refreshToken() async {
    try {
      final response = await _networkService.post(ApiConstants.refreshToken);

      if (response.statusCode == 200) {
        final apiResponse = ApiResponse<Map<String, dynamic>>.fromJson(
          response.data,
          (json) => json as Map<String, dynamic>,
        );
        
        if (apiResponse.success && apiResponse.data != null) {
          return LoginResponse.fromJson(apiResponse.data!);
        } else {
          throw Exception(apiResponse.error ?? 'Token refresh failed');
        }
      } else {
        throw Exception('Token refresh failed');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorData = e.response!.data;
        throw Exception(errorData['message'] ?? 'Token refresh failed');
      }
      throw Exception('Network error occurred');
    }
  }
}