import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';
import '../../widgets/common/custom_text_field.dart';

class MechanicPartsInventoryPage extends StatefulWidget {
  const MechanicPartsInventoryPage({super.key});

  @override
  State<MechanicPartsInventoryPage> createState() => _MechanicPartsInventoryPageState();
}

class _MechanicPartsInventoryPageState extends State<MechanicPartsInventoryPage> {
  final _searchController = TextEditingController();
  String _selectedCategory = 'all';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Parts Inventory'),
        backgroundColor: AppTheme.accentColor,
        foregroundColor: Colors.white,
        actions: [
          IconButton(
            icon: const Icon(Icons.add_shopping_cart),
            onPressed: () {
              // TODO: Navigate to parts request
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Parts Request - Coming Soon')),
              );
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Search and Filter Section
          Container(
            padding: const EdgeInsets.all(16),
            color: Colors.grey[50],
            child: Column(
              children: [
                CustomTextField(
                  controller: _searchController,
                  label: 'Search Parts',
                  hint: 'Enter part name or code',
                  prefixIcon: Icons.search,
                  onChanged: (value) {
                    // TODO: Implement search
                  },
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Text(
                      'Category:',
                      style: Theme.of(context).textTheme.labelLarge,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: DropdownButtonFormField<String>(
                        value: _selectedCategory,
                        decoration: InputDecoration(
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(8),
                          ),
                          contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                        ),
                        items: const [
                          DropdownMenuItem(value: 'all', child: Text('All Categories')),
                          DropdownMenuItem(value: 'engine', child: Text('Engine Parts')),
                          DropdownMenuItem(value: 'brake', child: Text('Brake System')),
                          DropdownMenuItem(value: 'electrical', child: Text('Electrical')),
                          DropdownMenuItem(value: 'suspension', child: Text('Suspension')),
                          DropdownMenuItem(value: 'transmission', child: Text('Transmission')),
                        ],
                        onChanged: (value) {
                          setState(() {
                            _selectedCategory = value!;
                          });
                          // TODO: Implement filter
                        },
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          
          // Quick Actions
          Container(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: () {
                      // TODO: Quick request common parts
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Quick Request - Coming Soon')),
                      );
                    },
                    icon: const Icon(Icons.flash_on),
                    label: const Text('Quick Request'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppTheme.primaryColor,
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: () {
                      // TODO: View pending requests
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Pending Requests - Coming Soon')),
                      );
                    },
                    icon: const Icon(Icons.pending_actions),
                    label: const Text('Pending'),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppTheme.warningColor,
                      side: const BorderSide(color: AppTheme.warningColor),
                    ),
                  ),
                ),
              ],
            ),
          ),
          
          // Parts List
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _getDummyParts().length,
              itemBuilder: (context, index) {
                final part = _getDummyParts()[index];
                return _buildPartCard(part);
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // TODO: Create new parts request
          _showPartsRequestDialog();
        },
        backgroundColor: AppTheme.accentColor,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildPartCard(Map<String, dynamic> part) {
    final availability = part['availability'] as String;
    Color availabilityColor;
    IconData availabilityIcon;
    
    switch (availability) {
      case 'available':
        availabilityColor = AppTheme.successColor;
        availabilityIcon = Icons.check_circle;
        break;
      case 'low_stock':
        availabilityColor = AppTheme.warningColor;
        availabilityIcon = Icons.warning;
        break;
      case 'out_of_stock':
        availabilityColor = AppTheme.errorColor;
        availabilityIcon = Icons.error;
        break;
      default:
        availabilityColor = AppTheme.textSecondary;
        availabilityIcon = Icons.help;
    }

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ExpansionTile(
        leading: Container(
          width: 50,
          height: 50,
          decoration: BoxDecoration(
            color: AppTheme.accentColor.withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Icon(
            _getCategoryIcon(part['category'] as String),
            color: AppTheme.accentColor,
          ),
        ),
        title: Text(
          part['name'] as String,
          style: const TextStyle(fontWeight: FontWeight.w600),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Code: ${part['code']}'),
            Text('Category: ${part['category']}'),
            Row(
              children: [
                Icon(availabilityIcon, color: availabilityColor, size: 16),
                const SizedBox(width: 4),
                Text(
                  '${part['stock']} available',
                  style: TextStyle(color: availabilityColor, fontWeight: FontWeight.w500),
                ),
              ],
            ),
          ],
        ),
        trailing: ElevatedButton(
          onPressed: availability == 'out_of_stock' ? null : () {
            _requestPart(part);
          },
          style: ElevatedButton.styleFrom(
            backgroundColor: AppTheme.accentColor,
            minimumSize: const Size(80, 32),
          ),
          child: const Text(
            'Request',
            style: TextStyle(fontSize: 12),
          ),
        ),
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.location_on, size: 16, color: Colors.grey[600]),
                    const SizedBox(width: 4),
                    Text('Location: ${part['location']}'),
                  ],
                ),
                const SizedBox(height: 8),
                Row(
                  children: [
                    Icon(Icons.attach_money, size: 16, color: Colors.grey[600]),
                    const SizedBox(width: 4),
                    Text('Unit Price: \$${part['price']}'),
                  ],
                ),
                const SizedBox(height: 8),
                if (part['last_used'] != null) ...[
                  Row(
                    children: [
                      Icon(Icons.history, size: 16, color: Colors.grey[600]),
                      const SizedBox(width: 4),
                      Text('Last Used: ${part['last_used']}'),
                    ],
                  ),
                  const SizedBox(height: 8),
                ],
                if (part['description'] != null) ...[
                  Text(
                    'Description:',
                    style: Theme.of(context).textTheme.labelMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(part['description'] as String),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  IconData _getCategoryIcon(String category) {
    switch (category.toLowerCase()) {
      case 'engine':
        return Icons.speed;
      case 'brake':
        return Icons.disc_full;
      case 'electrical':
        return Icons.electrical_services;
      case 'suspension':
        return Icons.height;
      case 'transmission':
        return Icons.settings;
      default:
        return Icons.build;
    }
  }

  void _requestPart(Map<String, dynamic> part) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Request ${part['name']}'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('How many units do you need?'),
            const SizedBox(height: 16),
            TextFormField(
              keyboardType: TextInputType.number,
              decoration: InputDecoration(
                labelText: 'Quantity',
                border: OutlineInputBorder(),
                hintText: '1',
              ),
              initialValue: '1',
            ),
            const SizedBox(height: 16),
            TextFormField(
              decoration: InputDecoration(
                labelText: 'Notes (optional)',
                border: OutlineInputBorder(),
                hintText: 'Enter any additional notes',
              ),
              maxLines: 3,
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.of(context).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Requested ${part['name']}')),
              );
            },
            child: const Text('Request'),
          ),
        ],
      ),
    );
  }

  void _showPartsRequestDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('New Parts Request'),
        content: const Text('This feature will allow you to create a new parts request form.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Close'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.of(context).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Parts Request Form - Coming Soon')),
              );
            },
            child: const Text('Create Request'),
          ),
        ],
      ),
    );
  }

  List<Map<String, dynamic>> _getDummyParts() {
    return [
      {
        'name': 'Engine Oil Filter',
        'code': 'EOF001',
        'category': 'Engine',
        'price': '15.99',
        'stock': '12',
        'location': 'A1-01',
        'availability': 'available',
        'last_used': '2 days ago',
        'description': 'High-quality oil filter compatible with most engine types.',
      },
      {
        'name': 'Brake Pads (Front)',
        'code': 'BPF001',
        'category': 'Brake',
        'price': '89.50',
        'stock': '3',
        'location': 'B2-03',
        'availability': 'low_stock',
        'last_used': '1 week ago',
        'description': 'Premium ceramic brake pads for enhanced stopping power.',
      },
      {
        'name': 'Spark Plug Set',
        'code': 'SPS001',
        'category': 'Engine',
        'price': '45.99',
        'stock': '0',
        'location': 'A2-02',
        'availability': 'out_of_stock',
        'last_used': '3 days ago',
        'description': 'High-performance spark plugs for improved engine efficiency.',
      },
      {
        'name': 'Alternator Belt',
        'code': 'AB001',
        'category': 'Electrical',
        'price': '25.99',
        'stock': '8',
        'location': 'C1-05',
        'availability': 'available',
        'last_used': '5 days ago',
        'description': 'Durable alternator belt for reliable electrical system operation.',
      },
      {
        'name': 'Shock Absorber',
        'code': 'SA001',
        'category': 'Suspension',
        'price': '125.00',
        'stock': '2',
        'location': 'D1-01',
        'availability': 'low_stock',
        'last_used': '1 day ago',
        'description': 'High-performance shock absorber for smooth ride quality.',
      },
    ];
  }
}