import 'package:json_annotation/json_annotation.dart';
import 'package:equatable/equatable.dart';

part 'api_response.g.dart';

@JsonSerializable(genericArgumentFactories: true)
class ApiResponse<T> extends Equatable {
  final bool success;
  final String message;
  final T? data;
  final String? error;
  final int? code;
  final DateTime timestamp;

  const ApiResponse({
    required this.success,
    required this.message,
    this.data,
    this.error,
    this.code,
    required this.timestamp,
  });

  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) => _$ApiResponseFromJson(json, fromJsonT);

  Map<String, dynamic> toJson(Object Function(T value) toJsonT) => _$ApiResponseToJson(this, toJsonT);

  @override
  List<Object?> get props => [success, message, data, error, code, timestamp];
}

@JsonSerializable()
class PaginationMeta extends Equatable {
  final int total;
  final int page;
  final int limit;
  @JsonKey(name: 'total_pages')
  final int totalPages;
  @JsonKey(name: 'has_more')
  final bool hasMore;

  const PaginationMeta({
    required this.total,
    required this.page,
    required this.limit,
    required this.totalPages,
    required this.hasMore,
  });

  factory PaginationMeta.fromJson(Map<String, dynamic> json) => _$PaginationMetaFromJson(json);
  Map<String, dynamic> toJson() => _$PaginationMetaToJson(this);

  @override
  List<Object> get props => [total, page, limit, totalPages, hasMore];
}

@JsonSerializable(genericArgumentFactories: true)
class PaginatedResponse<T> extends Equatable {
  final List<T> data;
  final PaginationMeta meta;

  const PaginatedResponse({
    required this.data,
    required this.meta,
  });

  factory PaginatedResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) => _$PaginatedResponseFromJson(json, fromJsonT);

  Map<String, dynamic> toJson(Object Function(T value) toJsonT) => _$PaginatedResponseToJson(this, toJsonT);

  @override
  List<Object> get props => [data, meta];
}

@JsonSerializable()
class ErrorResponse extends Equatable {
  final String message;
  final String? error;
  final String? details;
  final int? code;
  final DateTime timestamp;

  const ErrorResponse({
    required this.message,
    this.error,
    this.details,
    this.code,
    required this.timestamp,
  });

  factory ErrorResponse.fromJson(Map<String, dynamic> json) => _$ErrorResponseFromJson(json);
  Map<String, dynamic> toJson() => _$ErrorResponseToJson(this);

  @override
  List<Object?> get props => [message, error, details, code, timestamp];
}