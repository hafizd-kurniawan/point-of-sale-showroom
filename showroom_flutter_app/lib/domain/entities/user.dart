import 'package:equatable/equatable.dart';

enum UserRole {
  admin('admin'),
  sales('sales'),
  cashier('cashier'),
  mechanic('mechanic'),
  manager('manager');

  const UserRole(this.value);
  final String value;

  static UserRole fromString(String value) {
    switch (value.toLowerCase()) {
      case 'admin':
        return UserRole.admin;
      case 'sales':
        return UserRole.sales;
      case 'cashier':
        return UserRole.cashier;
      case 'mechanic':
        return UserRole.mechanic;
      case 'manager':
        return UserRole.manager;
      default:
        throw ArgumentError('Invalid user role: $value');
    }
  }

  bool get isAdmin => this == UserRole.admin;
  bool get isCashier => this == UserRole.cashier;
  bool get isMechanic => this == UserRole.mechanic;
  bool get isSales => this == UserRole.sales;
  bool get isManager => this == UserRole.manager;
}

class User extends Equatable {
  final int userId;
  final String username;
  final String email;
  final String fullName;
  final String phone;
  final String? address;
  final UserRole role;
  final double? salary;
  final DateTime? hireDate;
  final DateTime createdAt;
  final DateTime updatedAt;
  final int createdBy;
  final bool isActive;
  final String? profileImage;
  final String? notes;

  const User({
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

class UserInfo extends Equatable {
  final int userId;
  final String username;
  final String email;
  final String fullName;
  final UserRole role;
  final bool isActive;
  final String? profileImage;

  const UserInfo({
    required this.userId,
    required this.username,
    required this.email,
    required this.fullName,
    required this.role,
    required this.isActive,
    this.profileImage,
  });

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

class UserSession extends Equatable {
  final int sessionId;
  final int userId;
  final DateTime loginAt;
  final DateTime? logoutAt;
  final String? ipAddress;
  final String? userAgent;
  final bool isActive;
  final String? duration;

  const UserSession({
    required this.sessionId,
    required this.userId,
    required this.loginAt,
    this.logoutAt,
    this.ipAddress,
    this.userAgent,
    required this.isActive,
    this.duration,
  });

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