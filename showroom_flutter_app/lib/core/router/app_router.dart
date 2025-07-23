import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../di/dependency_injection.dart';
import '../storage/storage_service.dart';
import '../../domain/entities/user.dart';
import '../../presentation/pages/auth/login_page.dart';
import '../../presentation/pages/common/splash_page.dart';
import '../../presentation/pages/cashier/cashier_dashboard_page.dart';
import '../../presentation/pages/mechanic/mechanic_dashboard_page.dart';

class AppRouter {
  static final GoRouter router = GoRouter(
    initialLocation: '/splash',
    redirect: (BuildContext context, GoRouterState state) async {
      // Get auth status from storage
      final storageService = DependencyInjection.getIt<StorageService>();
      final isLoggedIn = await storageService.isLoggedIn();
      
      // Check if we're on auth-related pages
      final currentLocation = state.uri.toString();
      final isAuthPage = currentLocation.startsWith('/auth');
      final isSplashPage = currentLocation == '/splash';
      
      // If not logged in and not on auth/splash page, redirect to login
      if (!isLoggedIn && !isAuthPage && !isSplashPage) {
        return '/auth/login';
      }
      
      // If logged in and on auth page, redirect to appropriate dashboard
      if (isLoggedIn && isAuthPage) {
        final userData = await storageService.getUserData();
        if (userData != null) {
          final userRole = userData['role'] as String;
          switch (userRole) {
            case 'cashier':
              return '/cashier';
            case 'mechanic':
              return '/mechanic';
            default:
              return '/auth/login'; // Unsupported role
          }
        }
      }
      
      return null;
    },
    routes: [
      GoRoute(
        path: '/splash',
        builder: (context, state) => const SplashPage(),
      ),
      GoRoute(
        path: '/auth/login',
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: '/cashier',
        builder: (context, state) => const CashierDashboardPage(),
        routes: [
          GoRoute(
            path: '/profile',
            builder: (context, state) => const ProfilePage(),
          ),
          // Add more cashier-specific routes here
        ],
      ),
      GoRoute(
        path: '/mechanic',
        builder: (context, state) => const MechanicDashboardPage(),
        routes: [
          GoRoute(
            path: '/profile',
            builder: (context, state) => const ProfilePage(),
          ),
          // Add more mechanic-specific routes here
        ],
      ),
    ],
    errorBuilder: (context, state) => Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(
              Icons.error,
              size: 64,
              color: Colors.red,
            ),
            const SizedBox(height: 16),
            Text(
              'Page not found: ${state.uri}',
              style: const TextStyle(fontSize: 16),
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => context.go('/auth/login'),
              child: const Text('Go to Login'),
            ),
          ],
        ),
      ),
    ),
  );
}

// Placeholder for ProfilePage - will be implemented later
class ProfilePage extends StatelessWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Profile')),
      body: const Center(
        child: Text('Profile Page - To be implemented'),
      ),
    );
  }
}