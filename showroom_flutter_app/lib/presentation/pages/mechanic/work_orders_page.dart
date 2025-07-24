import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';

class WorkOrdersPage extends StatefulWidget {
  const WorkOrdersPage({super.key});

  @override
  State<WorkOrdersPage> createState() => _WorkOrdersPageState();
}

class _WorkOrdersPageState extends State<WorkOrdersPage> with TickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Work Orders'),
        backgroundColor: AppTheme.accentColor,
        foregroundColor: Colors.white,
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: Colors.white,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.white70,
          tabs: const [
            Tab(text: 'Pending'),
            Tab(text: 'In Progress'),
            Tab(text: 'Review'),
            Tab(text: 'Completed'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildWorkOrdersList('pending'),
          _buildWorkOrdersList('in_progress'),
          _buildWorkOrdersList('review'),
          _buildWorkOrdersList('completed'),
        ],
      ),
    );
  }

  Widget _buildWorkOrdersList(String status) {
    final workOrders = _getWorkOrdersByStatus(status);
    
    if (workOrders.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.assignment_outlined,
              size: 64,
              color: Colors.grey[400],
            ),
            const SizedBox(height: 16),
            Text(
              'No ${status.replaceAll('_', ' ')} work orders',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                color: Colors.grey[600],
              ),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: workOrders.length,
      itemBuilder: (context, index) {
        final workOrder = workOrders[index];
        return _buildWorkOrderCard(workOrder);
      },
    );
  }

  Widget _buildWorkOrderCard(Map<String, dynamic> workOrder) {
    final priority = workOrder['priority'] as String;
    Color priorityColor;
    
    switch (priority) {
      case 'high':
        priorityColor = AppTheme.errorColor;
        break;
      case 'medium':
        priorityColor = AppTheme.warningColor;
        break;
      case 'low':
        priorityColor = AppTheme.successColor;
        break;
      default:
        priorityColor = AppTheme.textSecondary;
    }

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ExpansionTile(
        leading: Container(
          width: 50,
          height: 50,
          decoration: BoxDecoration(
            color: priorityColor.withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Icon(
            Icons.build_circle,
            color: priorityColor,
          ),
        ),
        title: Text(
          workOrder['title'] as String,
          style: const TextStyle(fontWeight: FontWeight.w600),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Vehicle: ${workOrder['vehicle']}'),
            Text('Customer: ${workOrder['customer']}'),
            const SizedBox(height: 4),
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: priorityColor.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    '${priority.toUpperCase()} PRIORITY',
                    style: TextStyle(
                      color: priorityColor,
                      fontSize: 10,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                if (workOrder['estimated_hours'] != null)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: AppTheme.infoColor.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      '${workOrder['estimated_hours']}h',
                      style: const TextStyle(
                        color: AppTheme.infoColor,
                        fontSize: 10,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
              ],
            ),
          ],
        ),
        trailing: _buildActionButton(workOrder),
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildDetailRow('Work Order ID', workOrder['id']),
                _buildDetailRow('Assigned Date', workOrder['assigned_date']),
                _buildDetailRow('Due Date', workOrder['due_date']),
                if (workOrder['description'] != null) ...[
                  const SizedBox(height: 12),
                  Text(
                    'Description:',
                    style: Theme.of(context).textTheme.labelMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(workOrder['description'] as String),
                ],
                if (workOrder['parts_needed'] != null) ...[
                  const SizedBox(height: 12),
                  Text(
                    'Parts Needed:',
                    style: Theme.of(context).textTheme.labelMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 4),
                  ...((workOrder['parts_needed'] as List<String>).map(
                    (part) => Padding(
                      padding: const EdgeInsets.only(left: 8, bottom: 2),
                      child: Row(
                        children: [
                          const Icon(Icons.circle, size: 6),
                          const SizedBox(width: 8),
                          Expanded(child: Text(part)),
                        ],
                      ),
                    ),
                  )),
                ],
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton.icon(
                        onPressed: () => _showWorkOrderDetails(workOrder),
                        icon: const Icon(Icons.visibility),
                        label: const Text('View Details'),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton.icon(
                        onPressed: () => _updateWorkOrderStatus(workOrder),
                        icon: Icon(_getActionIcon(workOrder['status'])),
                        label: Text(_getActionText(workOrder['status'])),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: AppTheme.accentColor,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 100,
            child: Text(
              '$label:',
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: Theme.of(context).textTheme.bodySmall,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActionButton(Map<String, dynamic> workOrder) {
    final status = workOrder['status'] as String;
    Color buttonColor;
    IconData icon;
    
    switch (status) {
      case 'pending':
        buttonColor = AppTheme.primaryColor;
        icon = Icons.play_arrow;
        break;
      case 'in_progress':
        buttonColor = AppTheme.warningColor;
        icon = Icons.pause;
        break;
      case 'review':
        buttonColor = AppTheme.infoColor;
        icon = Icons.check;
        break;
      case 'completed':
        buttonColor = AppTheme.successColor;
        icon = Icons.check_circle;
        break;
      default:
        buttonColor = AppTheme.textSecondary;
        icon = Icons.help;
    }

    return IconButton(
      onPressed: () => _updateWorkOrderStatus(workOrder),
      icon: Icon(icon),
      color: buttonColor,
    );
  }

  IconData _getActionIcon(String status) {
    switch (status) {
      case 'pending':
        return Icons.play_arrow;
      case 'in_progress':
        return Icons.assignment_turned_in;
      case 'review':
        return Icons.check_circle;
      default:
        return Icons.visibility;
    }
  }

  String _getActionText(String status) {
    switch (status) {
      case 'pending':
        return 'Start Work';
      case 'in_progress':
        return 'Submit';
      case 'review':
        return 'Complete';
      default:
        return 'View';
    }
  }

  void _updateWorkOrderStatus(Map<String, dynamic> workOrder) {
    final currentStatus = workOrder['status'] as String;
    String newStatus;
    String action;
    
    switch (currentStatus) {
      case 'pending':
        newStatus = 'in_progress';
        action = 'started';
        break;
      case 'in_progress':
        newStatus = 'review';
        action = 'submitted for review';
        break;
      case 'review':
        newStatus = 'completed';
        action = 'completed';
        break;
      default:
        return;
    }

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Update Work Order'),
        content: Text('Are you sure you want to mark this work order as $action?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.of(context).pop();
              setState(() {
                workOrder['status'] = newStatus;
              });
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Work order $action')),
              );
            },
            child: const Text('Confirm'),
          ),
        ],
      ),
    );
  }

  void _showWorkOrderDetails(Map<String, dynamic> workOrder) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(workOrder['title'] as String),
        content: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              Text('Full work order details would be displayed here.'),
              const SizedBox(height: 16),
              Text('This would include:'),
              const SizedBox(height: 8),
              const Text('• Detailed vehicle information'),
              const Text('• Complete part requirements'),
              const Text('• Service history'),
              const Text('• Customer notes'),
              const Text('• Diagnostic reports'),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  List<Map<String, dynamic>> _getWorkOrdersByStatus(String status) {
    final allWorkOrders = _getDummyWorkOrders();
    return allWorkOrders.where((wo) => wo['status'] == status).toList();
  }

  List<Map<String, dynamic>> _getDummyWorkOrders() {
    return [
      {
        'id': 'WO-2024-001',
        'title': 'Engine Oil Change & Filter Replacement',
        'vehicle': '2020 Honda Civic',
        'customer': 'John Smith',
        'priority': 'low',
        'status': 'pending',
        'assigned_date': '2024-01-15',
        'due_date': '2024-01-16',
        'estimated_hours': '1',
        'description': 'Routine maintenance: oil change with high-quality synthetic oil and replacement of oil filter.',
        'parts_needed': ['Engine Oil 5W-30 (5L)', 'Oil Filter'],
      },
      {
        'id': 'WO-2024-002',
        'title': 'Brake System Inspection & Repair',
        'vehicle': '2018 Toyota Camry',
        'customer': 'Sarah Johnson',
        'priority': 'high',
        'status': 'in_progress',
        'assigned_date': '2024-01-14',
        'due_date': '2024-01-15',
        'estimated_hours': '3',
        'description': 'Customer reports squeaking noise when braking. Inspect brake pads, rotors, and fluid levels.',
        'parts_needed': ['Brake Pads (Front)', 'Brake Fluid'],
      },
      {
        'id': 'WO-2024-003',
        'title': 'Transmission Service',
        'vehicle': '2019 Ford F-150',
        'customer': 'Mike Wilson',
        'priority': 'medium',
        'status': 'review',
        'assigned_date': '2024-01-13',
        'due_date': '2024-01-15',
        'estimated_hours': '2',
        'description': 'Complete transmission service including fluid change and filter replacement.',
        'parts_needed': ['Transmission Fluid', 'Transmission Filter'],
      },
      {
        'id': 'WO-2024-004',
        'title': 'A/C System Diagnostic',
        'vehicle': '2021 BMW X5',
        'customer': 'Lisa Brown',
        'priority': 'medium',
        'status': 'completed',
        'assigned_date': '2024-01-12',
        'due_date': '2024-01-13',
        'estimated_hours': '1.5',
        'description': 'A/C not cooling properly. Diagnostic test completed, found refrigerant leak.',
        'parts_needed': ['R-134a Refrigerant', 'A/C Seal Kit'],
      },
    ];
  }
}