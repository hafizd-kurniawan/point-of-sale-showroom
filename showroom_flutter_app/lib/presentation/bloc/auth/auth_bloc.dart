import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../../domain/entities/auth.dart';
import '../../../domain/entities/user.dart';
import '../../../domain/usecases/auth_usecases.dart';

// Events
abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class AuthCheckRequested extends AuthEvent {}

class AuthLoginRequested extends AuthEvent {
  final String username;
  final String password;

  const AuthLoginRequested({
    required this.username,
    required this.password,
  });

  @override
  List<Object> get props => [username, password];
}

class AuthLogoutRequested extends AuthEvent {}

class AuthPasswordChangeRequested extends AuthEvent {
  final String currentPassword;
  final String newPassword;
  final String confirmPassword;

  const AuthPasswordChangeRequested({
    required this.currentPassword,
    required this.newPassword,
    required this.confirmPassword,
  });

  @override
  List<Object> get props => [currentPassword, newPassword, confirmPassword];
}

class AuthProfileRequested extends AuthEvent {}

class AuthTokenRefreshRequested extends AuthEvent {}

// States
abstract class AuthState extends Equatable {
  const AuthState();

  @override
  List<Object?> get props => [];
}

class AuthInitial extends AuthState {}

class AuthLoading extends AuthState {}

class AuthAuthenticated extends AuthState {
  final UserInfo user;

  const AuthAuthenticated(this.user);

  @override
  List<Object> get props => [user];
}

class AuthUnauthenticated extends AuthState {}

class AuthError extends AuthState {
  final String message;

  const AuthError(this.message);

  @override
  List<Object> get props => [message];
}

class AuthPasswordChangeSuccess extends AuthState {}

class AuthProfileLoaded extends AuthState {
  final UserProfile profile;

  const AuthProfileLoaded(this.profile);

  @override
  List<Object> get props => [profile];
}

// BLoC
class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final LoginUseCase _loginUseCase;
  final LogoutUseCase _logoutUseCase;
  final GetCurrentUserUseCase _getCurrentUserUseCase;
  final GetProfileUseCase _getProfileUseCase;
  final ChangePasswordUseCase _changePasswordUseCase;
  final CheckAuthStatusUseCase _checkAuthStatusUseCase;
  final GetCachedUserUseCase _getCachedUserUseCase;

  AuthBloc({
    required LoginUseCase loginUseCase,
    required LogoutUseCase logoutUseCase,
    required GetCurrentUserUseCase getCurrentUserUseCase,
    required GetProfileUseCase getProfileUseCase,
    required ChangePasswordUseCase changePasswordUseCase,
    required CheckAuthStatusUseCase checkAuthStatusUseCase,
    required GetCachedUserUseCase getCachedUserUseCase,
  })  : _loginUseCase = loginUseCase,
        _logoutUseCase = logoutUseCase,
        _getCurrentUserUseCase = getCurrentUserUseCase,
        _getProfileUseCase = getProfileUseCase,
        _changePasswordUseCase = changePasswordUseCase,
        _checkAuthStatusUseCase = checkAuthStatusUseCase,
        _getCachedUserUseCase = getCachedUserUseCase,
        super(AuthInitial()) {
    on<AuthCheckRequested>(_onAuthCheckRequested);
    on<AuthLoginRequested>(_onAuthLoginRequested);
    on<AuthLogoutRequested>(_onAuthLogoutRequested);
    on<AuthPasswordChangeRequested>(_onAuthPasswordChangeRequested);
    on<AuthProfileRequested>(_onAuthProfileRequested);
  }

  Future<void> _onAuthCheckRequested(
    AuthCheckRequested event,
    Emitter<AuthState> emit,
  ) async {
    try {
      final isLoggedIn = await _checkAuthStatusUseCase();
      if (isLoggedIn) {
        final cachedUser = await _getCachedUserUseCase();
        if (cachedUser != null) {
          emit(AuthAuthenticated(cachedUser));
        } else {
          // Try to get current user from server
          try {
            final user = await _getCurrentUserUseCase();
            emit(AuthAuthenticated(user));
          } catch (e) {
            // If failed, user needs to login again
            emit(AuthUnauthenticated());
          }
        }
      } else {
        emit(AuthUnauthenticated());
      }
    } catch (e) {
      emit(AuthUnauthenticated());
    }
  }

  Future<void> _onAuthLoginRequested(
    AuthLoginRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      final credentials = LoginCredentials(
        username: event.username,
        password: event.password,
      );
      
      final session = await _loginUseCase(credentials);
      emit(AuthAuthenticated(session.user));
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }

  Future<void> _onAuthLogoutRequested(
    AuthLogoutRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      await _logoutUseCase();
      emit(AuthUnauthenticated());
    } catch (e) {
      // Even if logout fails, clear local state
      emit(AuthUnauthenticated());
    }
  }

  Future<void> _onAuthPasswordChangeRequested(
    AuthPasswordChangeRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      final request = PasswordChangeRequest(
        currentPassword: event.currentPassword,
        newPassword: event.newPassword,
        confirmPassword: event.confirmPassword,
      );
      
      await _changePasswordUseCase(request);
      emit(AuthPasswordChangeSuccess());
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }

  Future<void> _onAuthProfileRequested(
    AuthProfileRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      final profile = await _getProfileUseCase();
      emit(AuthProfileLoaded(profile));
    } catch (e) {
      emit(AuthError(e.toString()));
    }
  }
}