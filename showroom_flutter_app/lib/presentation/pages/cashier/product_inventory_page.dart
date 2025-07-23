import 'package:flutter/material.dart';
import '../../../core/theme/app_theme.dart';
import '../../widgets/common/custom_text_field.dart';

class ProductInventoryPage extends StatefulWidget {
  const ProductInventoryPage({super.key});

  @override
  State<ProductInventoryPage> createState() => _ProductInventoryPageState();
}

class _ProductInventoryPageState extends State<ProductInventoryPage> {
  final _searchController = TextEditingController();
  String _selectedFilter = 'all';

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Product Inventory'),
        backgroundColor: AppTheme.primaryColor,
        foregroundColor: Colors.white,
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
                  label: 'Search Products',
                  hint: 'Enter product name or code',
                  prefixIcon: Icons.search,
                  onChanged: (value) {
                    // TODO: Implement search
                  },
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Text(
                      'Filter:',
                      style: Theme.of(context).textTheme.labelLarge,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: DropdownButtonFormField<String>(
                        value: _selectedFilter,
                        decoration: InputDecoration(
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(8),
                          ),
                          contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                        ),
                        items: const [
                          DropdownMenuItem(value: 'all', child: Text('All Products')),
                          DropdownMenuItem(value: 'in_stock', child: Text('In Stock')),
                          DropdownMenuItem(value: 'low_stock', child: Text('Low Stock')),
                          DropdownMenuItem(value: 'out_of_stock', child: Text('Out of Stock')),
                        ],
                        onChanged: (value) {
                          setState(() {
                            _selectedFilter = value!;
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
          
          // Quick Stats
          Container(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Expanded(
                  child: _buildStatCard(
                    'Total Products',
                    '248',
                    Icons.inventory,
                    AppTheme.primaryColor,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildStatCard(
                    'Low Stock',
                    '15',
                    Icons.warning,
                    AppTheme.warningColor,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _buildStatCard(
                    'Out of Stock',
                    '3',
                    Icons.error,
                    AppTheme.errorColor,
                  ),
                ),
              ],
            ),
          ),
          
          // Product List
          Expanded(
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: _getDummyProducts().length,
              itemBuilder: (context, index) {
                final product = _getDummyProducts()[index];
                return _buildProductCard(product);
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatCard(String title, String value, IconData icon, Color color) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            Icon(icon, color: color, size: 24),
            const SizedBox(height: 8),
            Text(
              value,
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
                color: color,
              ),
            ),
            Text(
              title,
              style: Theme.of(context).textTheme.bodySmall,
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildProductCard(Map<String, dynamic> product) {
    final stockStatus = product['stock_status'] as String;
    Color statusColor;
    IconData statusIcon;
    
    switch (stockStatus) {
      case 'in_stock':
        statusColor = AppTheme.successColor;
        statusIcon = Icons.check_circle;
        break;
      case 'low_stock':
        statusColor = AppTheme.warningColor;
        statusIcon = Icons.warning;
        break;
      case 'out_of_stock':
        statusColor = AppTheme.errorColor;
        statusIcon = Icons.error;
        break;
      default:
        statusColor = AppTheme.textSecondary;
        statusIcon = Icons.help;
    }

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ListTile(
        leading: Container(
          width: 50,
          height: 50,
          decoration: BoxDecoration(
            color: AppTheme.primaryColor.withOpacity(0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: const Icon(
            Icons.inventory_2,
            color: AppTheme.primaryColor,
          ),
        ),
        title: Text(
          product['name'] as String,
          style: const TextStyle(fontWeight: FontWeight.w600),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Code: ${product['code']}'),
            Text('Price: \$${product['price']}'),
            Row(
              children: [
                Icon(statusIcon, color: statusColor, size: 16),
                const SizedBox(width: 4),
                Text(
                  'Stock: ${product['stock']}',
                  style: TextStyle(color: statusColor, fontWeight: FontWeight.w500),
                ),
              ],
            ),
          ],
        ),
        trailing: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              product['location'] as String,
              style: Theme.of(context).textTheme.bodySmall,
            ),
            const SizedBox(height: 4),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(
                color: statusColor.withOpacity(0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Text(
                stockStatus.replaceAll('_', ' ').toUpperCase(),
                style: TextStyle(
                  color: statusColor,
                  fontSize: 10,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
          ],
        ),
        onTap: () {
          // TODO: Navigate to product details
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Product details for ${product['name']}')),
          );
        },
      ),
    );
  }

  List<Map<String, dynamic>> _getDummyProducts() {
    return [
      {
        'name': 'Engine Oil 5W-30',
        'code': 'EO5W30',
        'price': '45.99',
        'stock': '25',
        'location': 'A1-01',
        'stock_status': 'in_stock',
      },
      {
        'name': 'Brake Pads Front',
        'code': 'BPF001',
        'price': '89.50',
        'stock': '5',
        'location': 'B2-03',
        'stock_status': 'low_stock',
      },
      {
        'name': 'Air Filter',
        'code': 'AF001',
        'price': '25.99',
        'stock': '0',
        'location': 'C1-05',
        'stock_status': 'out_of_stock',
      },
      {
        'name': 'Spark Plugs Set',
        'code': 'SP001',
        'price': '65.99',
        'stock': '18',
        'location': 'A2-02',
        'stock_status': 'in_stock',
      },
      {
        'name': 'Transmission Fluid',
        'code': 'TF001',
        'price': '35.99',
        'stock': '8',
        'location': 'A1-03',
        'stock_status': 'low_stock',
      },
    ];
  }
}