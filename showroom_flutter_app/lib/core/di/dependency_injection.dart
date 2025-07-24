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
    try {
      print('Initializing dependencies...');
      
      // Storage
      final storageService = StorageService();
      await storageService.init();
      getIt.registerSingleton<StorageService>(storageService);
      print('Storage service initialized');

      // Network
      getIt.registerLazySingleton<NetworkService>(
        () => NetworkService(getIt<StorageService>()),
      );
      print('Network service registered');

      // Data sources
      getIt.registerLazySingleton<AuthRemoteDataSource>(
        () => AuthRemoteDataSourceImpl(getIt<NetworkService>()),
      );
      print('Auth data source registered');

      // Repositories
      getIt.registerLazySingleton<AuthRepository>(
        () => AuthRepositoryImpl(
          getIt<AuthRemoteDataSource>(),
          getIt<StorageService>(),
        ),
      );
      print('Auth repository registered');

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
      print('Use cases registered');

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
      print('Auth BLoC registered');
      print('Dependency injection completed successfully');
    } catch (e) {
      print('Dependency injection failed: $e');
      rethrow;
    }
  }
}