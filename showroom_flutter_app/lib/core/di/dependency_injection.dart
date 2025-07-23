import 'package:get_it/get_it.dart';
import '../network/network_service.dart';
import '../storage/storage_service.dart';
import '../../data/datasources/auth_remote_data_source.dart';
import '../../data/repositories/auth_repository_impl.dart';
import '../../domain/repositories/auth_repository.dart';
import '../../domain/usecases/auth_usecases.dart';
import '../../presentation/bloc/auth/auth_bloc.dart';

class DependencyInjection {
  static final getIt = GetIt.instance;

  static Future<void> init() async {
    // Storage
    final storageService = StorageService();
    await storageService.init();
    getIt.registerSingleton<StorageService>(storageService);

    // Network
    getIt.registerLazySingleton<NetworkService>(
      () => NetworkService(getIt<StorageService>()),
    );

    // Data sources
    getIt.registerLazySingleton<AuthRemoteDataSource>(
      () => AuthRemoteDataSourceImpl(getIt<NetworkService>()),
    );

    // Repositories
    getIt.registerLazySingleton<AuthRepository>(
      () => AuthRepositoryImpl(
        getIt<AuthRemoteDataSource>(),
        getIt<StorageService>(),
      ),
    );

    // Use cases
    getIt.registerLazySingleton<LoginUseCase>(
      () => LoginUseCase(getIt<AuthRepository>()),
    );
    getIt.registerLazySingleton<LogoutUseCase>(
      () => LogoutUseCase(getIt<AuthRepository>()),
    );
    getIt.registerLazySingleton<GetCurrentUserUseCase>(
      () => GetCurrentUserUseCase(getIt<AuthRepository>()),
    );
    getIt.registerLazySingleton<GetProfileUseCase>(
      () => GetProfileUseCase(getIt<AuthRepository>()),
    );
    getIt.registerLazySingleton<ChangePasswordUseCase>(
      () => ChangePasswordUseCase(getIt<AuthRepository>()),
    );
    getIt.registerLazySingleton<CheckAuthStatusUseCase>(
      () => CheckAuthStatusUseCase(getIt<AuthRepository>()),
    );
    getIt.registerLazySingleton<GetCachedUserUseCase>(
      () => GetCachedUserUseCase(getIt<AuthRepository>()),
    );

    // BLoCs
    getIt.registerFactory<AuthBloc>(
      () => AuthBloc(
        loginUseCase: getIt<LoginUseCase>(),
        logoutUseCase: getIt<LogoutUseCase>(),
        getCurrentUserUseCase: getIt<GetCurrentUserUseCase>(),
        getProfileUseCase: getIt<GetProfileUseCase>(),
        changePasswordUseCase: getIt<ChangePasswordUseCase>(),
        checkAuthStatusUseCase: getIt<CheckAuthStatusUseCase>(),
        getCachedUserUseCase: getIt<GetCachedUserUseCase>(),
      ),
    );
  }
}