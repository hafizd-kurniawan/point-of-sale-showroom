import 'package:json_annotation/json_annotation.dart';
import 'package:equatable/equatable.dart';
import 'user_model.dart';

part 'auth_model.g.dart';

@JsonSerializable()
class LoginRequest extends Equatable {
  final String username;
  final String password;

  const LoginRequest({
    required this.username,
    required this.password,
  });

  factory LoginRequest.fromJson(Map<String, dynamic> json) => _$LoginRequestFromJson(json);
  Map<String, dynamic> toJson() => _$LoginRequestToJson(this);

  @override
  List<Object> get props => [username, password];
}

@JsonSerializable()
class LoginResponse extends Equatable {
  final String token;
  @JsonKey(name: 'token_type')
  final String tokenType;
  @JsonKey(name: 'expires_in')
  final int expiresIn;
  @JsonKey(name: 'expires_at')
  final DateTime expiresAt;
  final UserInfoModel user;
  final String message;
  @JsonKey(name: 'session_id')
  final int sessionId;

  const LoginResponse({
    required this.token,
    required this.tokenType,
    required this.expiresIn,
    required this.expiresAt,
    required this.user,
    required this.message,
    required this.sessionId,
  });

  factory LoginResponse.fromJson(Map<String, dynamic> json) => _$LoginResponseFromJson(json);
  Map<String, dynamic> toJson() => _$LoginResponseToJson(this);

  @override
  List<Object> get props => [
        token,
        tokenType,
        expiresIn,
        expiresAt,
        user,
        message,
        sessionId,
      ];
}

@JsonSerializable()
class ChangePasswordRequest extends Equatable {
  @JsonKey(name: 'current_password')
  final String currentPassword;
  @JsonKey(name: 'new_password')
  final String newPassword;
  @JsonKey(name: 'confirm_password')
  final String confirmPassword;

  const ChangePasswordRequest({
    required this.currentPassword,
    required this.newPassword,
    required this.confirmPassword,
  });

  factory ChangePasswordRequest.fromJson(Map<String, dynamic> json) => _$ChangePasswordRequestFromJson(json);
  Map<String, dynamic> toJson() => _$ChangePasswordRequestToJson(this);

  @override
  List<Object> get props => [currentPassword, newPassword, confirmPassword];
}

@JsonSerializable()
class ProfileResponse extends Equatable {
  final UserModel user;
  final List<UserSessionModel> sessions;

  const ProfileResponse({
    required this.user,
    required this.sessions,
  });

  factory ProfileResponse.fromJson(Map<String, dynamic> json) => _$ProfileResponseFromJson(json);
  Map<String, dynamic> toJson() => _$ProfileResponseToJson(this);

  @override
  List<Object> get props => [user, sessions];
}