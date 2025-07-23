import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../bloc/auth/auth_bloc.dart';
import '../../widgets/common/custom_app_bar.dart';
import '../../widgets/common/dashboard_card.dart';
import '../../../core/theme/app_theme.dart';
import 'product_inventory_page.dart';

class CashierDashboardPage extends StatelessWidget {
  const CashierDashboardPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: CustomAppBar(
        title: 'Cashier Dashboard',
        onProfileTap: () => context.push('/cashier/profile'),
        onLogoutTap: () {
          context.read<AuthBloc>().add(AuthLogoutRequested());
        },
      ),
      body: BlocListener<AuthBloc, AuthState>(
        listener: (context, state) {
          if (state is AuthUnauthenticated) {
            context.go('/auth/login');
          }
        },
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Welcome section
              BlocBuilder<AuthBloc, AuthState>(
                builder: (context, state) {
                  String userName = 'Cashier';
                  if (state is AuthAuthenticated) {
                    userName = state.user.fullName;
                  }
                  
                  return Card(
                    child: Padding(
                      padding: const EdgeInsets.all(20.0),
                      child: Row(
                        children: [
                          CircleAvatar(
                            radius: 30,
                            backgroundColor: AppTheme.primaryColor,
                            child: const Icon(
                              Icons.person,
                              size: 30,
                              color: Colors.white,
                            ),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'Welcome back,',
                                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                    color: AppTheme.textSecondary,
                                  ),
                                ),
                                Text(
                                  userName,
                                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                                    fontWeight: FontWeight.bold,
                                    color: AppTheme.textPrimary,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Icon(
                            Icons.point_of_sale,
                            size: 40,
                            color: AppTheme.primaryColor,
                          ),
                        ],
                      ),
                    ),
                  );
                },
              ),
              const SizedBox(height: 24),
              
              // Quick actions
              Text(
                'Quick Actions',
                style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: AppTheme.textPrimary,
                ),
              ),
              const SizedBox(height: 16),
              
              GridView.count(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                crossAxisCount: 2,
                crossAxisSpacing: 16,
                mainAxisSpacing: 16,
                children: [
                  DashboardCard(
                    title: 'Vehicle Sales',
                    subtitle: 'Process vehicle sales',
                    icon: Icons.directions_car,
                    color: AppTheme.primaryColor,
                    onTap: () {
                      // TODO: Navigate to vehicle sales
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Vehicle Sales - Coming Soon')),
                      );
                    },
                  ),
                  DashboardCard(
                    title: 'Customer Info',
                    subtitle: 'View customer details',
                    icon: Icons.people,
                    color: AppTheme.successColor,
                    onTap: () {
                      // TODO: Navigate to customer info
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Customer Info - Coming Soon')),
                      );
                    },
                  ),
                  DashboardCard(
                    title: 'Product Inventory',
                    subtitle: 'Browse available products',
                    icon: Icons.inventory_2,
                    color: AppTheme.infoColor,
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => const ProductInventoryPage(),
                        ),
                      );
                    },
                  ),
                  DashboardCard(
                    title: 'Reports',
                    subtitle: 'Sales reports',
                    icon: Icons.assessment,
                    color: AppTheme.infoColor,
                    onTap: () {
                      // TODO: Navigate to reports
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Reports - Coming Soon')),
                      );
                    },
                  ),
                ],
              ),
              const SizedBox(height: 24),
              
              // Information section
              Text(
                'Information',
                style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: AppTheme.textPrimary,
                ),
              ),
              const SizedBox(height: 16),
              
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    children: [
                      ListTile(
                        leading: Icon(
                          Icons.info_outline,
                          color: AppTheme.infoColor,
                        ),
                        title: const Text('Vehicle Inventory'),
                        subtitle: const Text('Current vehicle stock information'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          // TODO: Show vehicle inventory info
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('Vehicle Inventory Info - Coming Soon')),
                          );
                        },
                      ),
                      const Divider(),
                      ListTile(
                        leading: Icon(
                          Icons.trending_up,
                          color: AppTheme.successColor,
                        ),
                        title: const Text('Sales Summary'),
                        subtitle: const Text('Today\'s sales performance'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          // TODO: Show sales summary
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('Sales Summary - Coming Soon')),
                          );
                        },
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}