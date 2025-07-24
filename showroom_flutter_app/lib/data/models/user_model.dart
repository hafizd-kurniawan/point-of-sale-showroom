import 'package:json_annotation/json_annotation.dart';
import 'package:equatable/equatable.dart';

part 'user_model.g.dart';

@JsonSerializable()
class UserModel extends Equatable {
  @JsonKey(name: 'user_id')
  final int userId;
  final String username;
  final String email;
  @JsonKey(name: 'full_name')
  final String fullName;
  final String phone;
  final String? address;
  final String role;
  final double? salary;
  @JsonKey(name: 'hire_date')
  final DateTime? hireDate;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;
  @JsonKey(name: 'created_by')
  final int createdBy;
  @JsonKey(name: 'is_active')
  final bool isActive;
  @JsonKey(name: 'profile_image')
  final String? profileImage;
  final String? notes;

  const UserModel({
    required this.userId,
    required this.username,
    required this.email,
    required this.fullName,
    required this.phone,
    this.address,
    required this.role,
    this.salary,
    this.hireDate,
    required this.createdAt,
    required this.updatedAt,
    required this.createdBy,
    required this.isActive,
    this.profileImage,
    this.notes,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) => _$UserModelFromJson(json);
  Map<String, dynamic> toJson() => _$UserModelToJson(this);

  @override
  List<Object?> get props => [
        userId,
        username,
        email,
        fullName,
        phone,
        address,
        role,
        salary,
        hireDate,
        createdAt,
        updatedAt,
        createdBy,
        isActive,
        profileImage,
        notes,
      ];
}

@JsonSerializable()
class UserInfoModel extends Equatable {
  @JsonKey(name: 'user_id')
  final int userId;
  final String username;
  final String email;
  @JsonKey(name: 'full_name')
  final String fullName;
  final String role;
  @JsonKey(name: 'is_active')
  final bool isActive;
  @JsonKey(name: 'profile_image')
  final String? profileImage;

  const UserInfoModel({
    required this.userId,
    required this.username,
    required this.email,
    required this.fullName,
    required this.role,
    required this.isActive,
    this.profileImage,
  });

  factory UserInfoModel.fromJson(Map<String, dynamic> json) => _$UserInfoModelFromJson(json);
  Map<String, dynamic> toJson() => _$UserInfoModelToJson(this);

  @override
  List<Object?> get props => [
        userId,
        username,
        email,
        fullName,
        role,
        isActive,
        profileImage,
      ];
}

@JsonSerializable()
class UserSessionModel extends Equatable {
  @JsonKey(name: 'session_id')
  final int sessionId;
  @JsonKey(name: 'user_id')
  final int userId;
  @JsonKey(name: 'login_at')
  final DateTime loginAt;
  @JsonKey(name: 'logout_at')
  final DateTime? logoutAt;
  @JsonKey(name: 'ip_address')
  final String? ipAddress;
  @JsonKey(name: 'user_agent')
  final String? userAgent;
  @JsonKey(name: 'is_active')
  final bool isActive;
  final String? duration;

  const UserSessionModel({
    required this.sessionId,
    required this.userId,
    required this.loginAt,
    this.logoutAt,
    this.ipAddress,
    this.userAgent,
    required this.isActive,
    this.duration,
  });

  factory UserSessionModel.fromJson(Map<String, dynamic> json) => _$UserSessionModelFromJson(json);
  Map<String, dynamic> toJson() => _$UserSessionModelToJson(this);

  @override
  List<Object?> get props => [
        sessionId,
        userId,
        loginAt,
        logoutAt,
        ipAddress,
        userAgent,
        isActive,
        duration,
      ];
}