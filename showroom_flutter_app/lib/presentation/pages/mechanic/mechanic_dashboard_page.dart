import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../bloc/auth/auth_bloc.dart';
import '../../widgets/common/custom_app_bar.dart';
import '../../widgets/common/dashboard_card.dart';
import '../../../core/theme/app_theme.dart';
import 'parts_inventory_page.dart';
import 'work_orders_page.dart';

class MechanicDashboardPage extends StatelessWidget {
  const MechanicDashboardPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: CustomAppBar(
        title: 'Mechanic Dashboard',
        onProfileTap: () => context.push('/mechanic/profile'),
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
                  String userName = 'Mechanic';
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
                            backgroundColor: AppTheme.accentColor,
                            child: const Icon(
                              Icons.build,
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
                            Icons.engineering,
                            size: 40,
                            color: AppTheme.accentColor,
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
                'Work Orders',
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
                    title: 'Repair Queue',
                    subtitle: 'Pending repair jobs',
                    icon: Icons.build_circle,
                    color: AppTheme.warningColor,
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => const WorkOrdersPage(),
                        ),
                      );
                    },
                  ),
                  DashboardCard(
                    title: 'Parts Inventory',
                    subtitle: 'Browse spare parts',
                    icon: Icons.inventory,
                    color: AppTheme.primaryColor,
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => const MechanicPartsInventoryPage(),
                        ),
                      );
                    },
                  ),
                  DashboardCard(
                    title: 'Work History',
                    subtitle: 'Completed repairs',
                    icon: Icons.history,
                    color: AppTheme.successColor,
                    onTap: () {
                      // TODO: Navigate to work history
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Work History - Coming Soon')),
                      );
                    },
                  ),
                  DashboardCard(
                    title: 'Time Tracking',
                    subtitle: 'Log work hours',
                    icon: Icons.access_time,
                    color: AppTheme.infoColor,
                    onTap: () {
                      // TODO: Navigate to time tracking
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Time Tracking - Coming Soon')),
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
                          Icons.inventory,
                          color: AppTheme.infoColor,
                        ),
                        title: const Text('Parts Inventory'),
                        subtitle: const Text('Available spare parts'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          // TODO: Show parts inventory info
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('Parts Inventory Info - Coming Soon')),
                          );
                        },
                      ),
                      const Divider(),
                      ListTile(
                        leading: Icon(
                          Icons.schedule,
                          color: AppTheme.warningColor,
                        ),
                        title: const Text('Work Schedule'),
                        subtitle: const Text('Today\'s scheduled repairs'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          // TODO: Show work schedule
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('Work Schedule - Coming Soon')),
                          );
                        },
                      ),
                      const Divider(),
                      ListTile(
                        leading: Icon(
                          Icons.emergency,
                          color: AppTheme.errorColor,
                        ),
                        title: const Text('Priority Repairs'),
                        subtitle: const Text('Urgent repair requests'),
                        trailing: const Icon(Icons.chevron_right),
                        onTap: () {
                          // TODO: Show priority repairs
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('Priority Repairs - Coming Soon')),
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