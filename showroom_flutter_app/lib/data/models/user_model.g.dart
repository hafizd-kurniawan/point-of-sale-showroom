// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserModel _$UserModelFromJson(Map<String, dynamic> json) => UserModel(
      userId: json['user_id'] as int,
      username: json['username'] as String,
      email: json['email'] as String,
      fullName: json['full_name'] as String,
      phone: json['phone'] as String,
      address: json['address'] as String?,
      role: json['role'] as String,
      salary: (json['salary'] as num?)?.toDouble(),
      hireDate: json['hire_date'] == null
          ? null
          : DateTime.parse(json['hire_date'] as String),
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
      createdBy: json['created_by'] as int,
      isActive: json['is_active'] as bool,
      profileImage: json['profile_image'] as String?,
      notes: json['notes'] as String?,
    );

Map<String, dynamic> _$UserModelToJson(UserModel instance) => <String, dynamic>{
      'user_id': instance.userId,
      'username': instance.username,
      'email': instance.email,
      'full_name': instance.fullName,
      'phone': instance.phone,
      'address': instance.address,
      'role': instance.role,
      'salary': instance.salary,
      'hire_date': instance.hireDate?.toIso8601String(),
      'created_at': instance.createdAt.toIso8601String(),
      'updated_at': instance.updatedAt.toIso8601String(),
      'created_by': instance.createdBy,
      'is_active': instance.isActive,
      'profile_image': instance.profileImage,
      'notes': instance.notes,
    };

UserInfoModel _$UserInfoModelFromJson(Map<String, dynamic> json) => UserInfoModel(
      userId: json['user_id'] as int,
      username: json['username'] as String,
      email: json['email'] as String,
      fullName: json['full_name'] as String,
      role: json['role'] as String,
      isActive: json['is_active'] as bool,
      profileImage: json['profile_image'] as String?,
    );

Map<String, dynamic> _$UserInfoModelToJson(UserInfoModel instance) =>
    <String, dynamic>{
      'user_id': instance.userId,
      'username': instance.username,
      'email': instance.email,
      'full_name': instance.fullName,
      'role': instance.role,
      'is_active': instance.isActive,
      'profile_image': instance.profileImage,
    };

UserSessionModel _$UserSessionModelFromJson(Map<String, dynamic> json) =>
    UserSessionModel(
      sessionId: json['session_id'] as int,
      userId: json['user_id'] as int,
      loginAt: DateTime.parse(json['login_at'] as String),
      logoutAt: json['logout_at'] == null
          ? null
          : DateTime.parse(json['logout_at'] as String),
      ipAddress: json['ip_address'] as String?,
      userAgent: json['user_agent'] as String?,
      isActive: json['is_active'] as bool,
      duration: json['duration'] as String?,
    );

Map<String, dynamic> _$UserSessionModelToJson(UserSessionModel instance) =>
    <String, dynamic>{
      'session_id': instance.sessionId,
      'user_id': instance.userId,
      'login_at': instance.loginAt.toIso8601String(),
      'logout_at': instance.logoutAt?.toIso8601String(),
      'ip_address': instance.ipAddress,
      'user_agent': instance.userAgent,
      'is_active': instance.isActive,
      'duration': instance.duration,
    };