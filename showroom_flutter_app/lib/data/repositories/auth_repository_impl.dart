import '../../domain/entities/auth.dart';
import '../../domain/entities/user.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_remote_data_source.dart';
import '../models/auth_model.dart';
import '../models/user_model.dart';
import '../../core/storage/storage_service.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource _remoteDataSource;
  final StorageService _storageService;

  AuthRepositoryImpl(this._remoteDataSource, this._storageService);

  @override
  Future<AuthSession> login(LoginCredentials credentials) async {
    final request = LoginRequest(
      username: credentials.username,
      password: credentials.password,
    );

    final response = await _remoteDataSource.login(request);
    
    // Store tokens and user data
    await _storageService.saveAccessToken(response.token);
    await _storageService.saveUserData(response.user.toJson());
    await _storageService.setLoggedIn(true);

    return AuthSession(
      token: response.token,
      tokenType: response.tokenType,
      expiresIn: response.expiresIn,
      expiresAt: response.expiresAt,
      user: _mapUserInfoModelToEntity(response.user),
      sessionId: response.sessionId,
    );
  }

  @override
  Future<void> logout() async {
    try {
      await _remoteDataSource.logout();
    } finally {
      // Always clear local data even if server logout fails
      await clearSession();
    }
  }

  @override
  Future<UserInfo> getCurrentUser() async {
    final userModel = await _remoteDataSource.getCurrentUser();
    return _mapUserInfoModelToEntity(userModel);
  }

  @override
  Future<UserProfile> getProfile() async {
    final profileResponse = await _remoteDataSource.getProfile();
    
    return UserProfile(
      user: _mapUserModelToEntity(profileResponse.user),
      sessions: profileResponse.sessions.map((session) => _mapUserSessionModelToEntity(session)).toList(),
    );
  }

  @override
  Future<void> changePassword(PasswordChangeRequest request) async {
    final changePasswordRequest = ChangePasswordRequest(
      currentPassword: request.currentPassword,
      newPassword: request.newPassword,
      confirmPassword: request.confirmPassword,
    );

    await _remoteDataSource.changePassword(changePasswordRequest);
  }

  @override
  Future<AuthSession> refreshToken() async {
    final response = await _remoteDataSource.refreshToken();
    
    // Update stored token and user data
    await _storageService.saveAccessToken(response.token);
    await _storageService.saveUserData(response.user.toJson());

    return AuthSession(
      token: response.token,
      tokenType: response.tokenType,
      expiresIn: response.expiresIn,
      expiresAt: response.expiresAt,
      user: _mapUserInfoModelToEntity(response.user),
      sessionId: response.sessionId,
    );
  }

  @override
  Future<bool> isLoggedIn() async {
    return await _storageService.isLoggedIn();
  }

  @override
  Future<void> clearSession() async {
    await _storageService.clearTokens();
  }

  @override
  Future<UserInfo?> getCachedUser() async {
    final userData = await _storageService.getUserData();
    if (userData != null) {
      final userModel = UserInfoModel.fromJson(userData);
      return _mapUserInfoModelToEntity(userModel);
    }
    return null;
  }

  // Helper methods for mapping between models and entities
  UserInfo _mapUserInfoModelToEntity(UserInfoModel model) {
    return UserInfo(
      userId: model.userId,
      username: model.username,
      email: model.email,
      fullName: model.fullName,
      role: UserRole.fromString(model.role),
      isActive: model.isActive,
      profileImage: model.profileImage,
    );
  }

  User _mapUserModelToEntity(UserModel model) {
    return User(
      userId: model.userId,
      username: model.username,
      email: model.email,
      fullName: model.fullName,
      phone: model.phone,
      address: model.address,
      role: UserRole.fromString(model.role),
      salary: model.salary,
      hireDate: model.hireDate,
      createdAt: model.createdAt,
      updatedAt: model.updatedAt,
      createdBy: model.createdBy,
      isActive: model.isActive,
      profileImage: model.profileImage,
      notes: model.notes,
    );
  }

  UserSession _mapUserSessionModelToEntity(UserSessionModel model) {
    return UserSession(
      sessionId: model.sessionId,
      userId: model.userId,
      loginAt: model.loginAt,
      logoutAt: model.logoutAt,
      ipAddress: model.ipAddress,
      userAgent: model.userAgent,
      isActive: model.isActive,
      duration: model.duration,
    );
  }
}