import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'core/di/dependency_injection.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';
import 'presentation/bloc/auth/auth_bloc.dart';
import 'core/storage/storage_service.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  try {
    print('Starting app initialization...');
    // Initialize dependencies
    await DependencyInjection.init();
    print('Dependencies initialized successfully');
    runApp(const ShowroomPOSApp());
  } catch (e) {
    print('Failed to initialize app: $e');
    // Run a minimal error app if initialization fails
    runApp(MaterialApp(
      title: 'Showroom POS - Error',
      home: Scaffold(
        backgroundColor: Colors.red.shade50,
        body: Center(
          child: Padding(
            padding: const EdgeInsets.all(32.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.error_outline, size: 80, color: Colors.red.shade700),
                const SizedBox(height: 24),
                Text(
                  'App Initialization Failed',
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    color: Colors.red.shade800,
                  ),
                ),
                const SizedBox(height: 16),
                Text(
                  'The app failed to start properly. This could be due to missing dependencies or configuration issues.',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 16,
                    color: Colors.red.shade700,
                  ),
                ),
                const SizedBox(height: 16),
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.red.shade100,
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: Colors.red.shade300),
                  ),
                  child: Text(
                    'Error: $e',
                    style: TextStyle(
                      fontSize: 14,
                      fontFamily: 'monospace',
                      color: Colors.red.shade800,
                    ),
                  ),
                ),
                const SizedBox(height: 32),
                ElevatedButton.icon(
                  onPressed: () {
                    // Try to restart the app
                    main();
                  },
                  icon: const Icon(Icons.refresh),
                  label: const Text('Retry Initialization'),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.red.shade700,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    ));
  }
}

class ShowroomPOSApp extends StatelessWidget {
  const ShowroomPOSApp({super.key});

  @override
  Widget build(BuildContext context) {
    print('Building ShowroomPOSApp');
    
    return MultiBlocProvider(
      providers: [
        BlocProvider<AuthBloc>(
          create: (context) {
            print('Creating AuthBloc');
            try {
              final authBloc = DependencyInjection.getIt<AuthBloc>();
              print('AuthBloc created successfully');
              return authBloc;
            } catch (e) {
              print('Failed to create AuthBloc: $e');
              rethrow;
            }
          },
        ),
      ],
      child: MaterialApp.router(
        title: 'Showroom POS',
        theme: AppTheme.lightTheme,
        darkTheme: AppTheme.darkTheme,
        routerConfig: AppRouter.router,
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}