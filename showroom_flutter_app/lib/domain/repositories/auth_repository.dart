import '../entities/auth.dart';
import '../entities/user.dart';

abstract class AuthRepository {
  Future<AuthSession> login(LoginCredentials credentials);
  Future<void> logout();
  Future<UserInfo> getCurrentUser();
  Future<UserProfile> getProfile();
  Future<void> changePassword(PasswordChangeRequest request);
  Future<AuthSession> refreshToken();
  Future<bool> isLoggedIn();
  Future<void> clearSession();
  Future<UserInfo?> getCachedUser();
}