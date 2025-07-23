import '../entities/auth.dart';
import '../entities/user.dart';
import '../repositories/auth_repository.dart';

class LoginUseCase {
  final AuthRepository _repository;

  LoginUseCase(this._repository);

  Future<AuthSession> call(LoginCredentials credentials) async {
    if (credentials.username.isEmpty || credentials.password.isEmpty) {
      throw Exception('Username and password are required');
    }
    
    return await _repository.login(credentials);
  }
}

class LogoutUseCase {
  final AuthRepository _repository;

  LogoutUseCase(this._repository);

  Future<void> call() async {
    await _repository.logout();
    await _repository.clearSession();
  }
}

class GetCurrentUserUseCase {
  final AuthRepository _repository;

  GetCurrentUserUseCase(this._repository);

  Future<UserInfo> call() async {
    return await _repository.getCurrentUser();
  }
}

class GetProfileUseCase {
  final AuthRepository _repository;

  GetProfileUseCase(this._repository);

  Future<UserProfile> call() async {
    return await _repository.getProfile();
  }
}

class ChangePasswordUseCase {
  final AuthRepository _repository;

  ChangePasswordUseCase(this._repository);

  Future<void> call(PasswordChangeRequest request) async {
    if (!request.isValid) {
      throw Exception('Invalid password change request');
    }
    
    return await _repository.changePassword(request);
  }
}

class CheckAuthStatusUseCase {
  final AuthRepository _repository;

  CheckAuthStatusUseCase(this._repository);

  Future<bool> call() async {
    return await _repository.isLoggedIn();
  }
}

class GetCachedUserUseCase {
  final AuthRepository _repository;

  GetCachedUserUseCase(this._repository);

  Future<UserInfo?> call() async {
    return await _repository.getCachedUser();
  }
}