import 'package:equatable/equatable.dart';
import 'user.dart';

class LoginCredentials extends Equatable {
  final String username;
  final String password;

  const LoginCredentials({
    required this.username,
    required this.password,
  });

  @override
  List<Object> get props => [username, password];
}

class AuthSession extends Equatable {
  final String token;
  final String tokenType;
  final int expiresIn;
  final DateTime expiresAt;
  final UserInfo user;
  final int sessionId;

  const AuthSession({
    required this.token,
    required this.tokenType,
    required this.expiresIn,
    required this.expiresAt,
    required this.user,
    required this.sessionId,
  });

  bool get isExpired => DateTime.now().isAfter(expiresAt);

  @override
  List<Object> get props => [
        token,
        tokenType,
        expiresIn,
        expiresAt,
        user,
        sessionId,
      ];
}

class PasswordChangeRequest extends Equatable {
  final String currentPassword;
  final String newPassword;
  final String confirmPassword;

  const PasswordChangeRequest({
    required this.currentPassword,
    required this.newPassword,
    required this.confirmPassword,
  });

  bool get isValid => newPassword == confirmPassword && newPassword.length >= 6;

  @override
  List<Object> get props => [currentPassword, newPassword, confirmPassword];
}

class UserProfile extends Equatable {
  final User user;
  final List<UserSession> sessions;

  const UserProfile({
    required this.user,
    required this.sessions,
  });

  @override
  List<Object> get props => [user, sessions];
}