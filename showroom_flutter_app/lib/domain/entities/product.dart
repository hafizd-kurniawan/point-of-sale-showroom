import 'package:equatable/equatable.dart';

class Product extends Equatable {
  final int productId;
  final String productCode;
  final String productName;
  final String? description;
  final int brandId;
  final int categoryId;
  final String unitMeasure;
  final double costPrice;
  final double sellingPrice;
  final double markupPercentage;
  final int stockQuantity;
  final int minStockLevel;
  final int maxStockLevel;
  final String? locationRack;
  final String? barcode;
  final double? weight;
  final String? dimensions;
  final DateTime createdAt;
  final DateTime updatedAt;
  final int createdBy;
  final bool isActive;
  final String? productImage;
  final String? notes;

  const Product({
    required this.productId,
    required this.productCode,
    required this.productName,
    this.description,
    required this.brandId,
    required this.categoryId,
    required this.unitMeasure,
    required this.costPrice,
    required this.sellingPrice,
    required this.markupPercentage,
    required this.stockQuantity,
    required this.minStockLevel,
    required this.maxStockLevel,
    this.locationRack,
    this.barcode,
    this.weight,
    this.dimensions,
    required this.createdAt,
    required this.updatedAt,
    required this.createdBy,
    required this.isActive,
    this.productImage,
    this.notes,
  });

  bool get isLowStock => stockQuantity <= minStockLevel;
  bool get isOutOfStock => stockQuantity <= 0;
  
  StockStatus get stockStatus {
    if (isOutOfStock) return StockStatus.outOfStock;
    if (isLowStock) return StockStatus.lowStock;
    return StockStatus.inStock;
  }

  @override
  List<Object?> get props => [
        productId,
        productCode,
        productName,
        description,
        brandId,
        categoryId,
        unitMeasure,
        costPrice,
        sellingPrice,
        markupPercentage,
        stockQuantity,
        minStockLevel,
        maxStockLevel,
        locationRack,
        barcode,
        weight,
        dimensions,
        createdAt,
        updatedAt,
        createdBy,
        isActive,
        productImage,
        notes,
      ];
}

enum StockStatus {
  inStock,
  lowStock,
  outOfStock;

  String get displayName {
    switch (this) {
      case StockStatus.inStock:
        return 'In Stock';
      case StockStatus.lowStock:
        return 'Low Stock';
      case StockStatus.outOfStock:
        return 'Out of Stock';
    }
  }
}

enum MovementType {
  stockIn('stock_in'),
  stockOut('stock_out'),
  adjustment('adjustment'),
  transfer('transfer'),
  purchase('purchase'),
  sale('sale'),
  return_('return');

  const MovementType(this.value);
  final String value;

  static MovementType fromString(String value) {
    switch (value.toLowerCase()) {
      case 'stock_in':
        return MovementType.stockIn;
      case 'stock_out':
        return MovementType.stockOut;
      case 'adjustment':
        return MovementType.adjustment;
      case 'transfer':
        return MovementType.transfer;
      case 'purchase':
        return MovementType.purchase;
      case 'sale':
        return MovementType.sale;
      case 'return':
        return MovementType.return_;
      default:
        throw ArgumentError('Invalid movement type: $value');
    }
  }

  String get displayName {
    switch (this) {
      case MovementType.stockIn:
        return 'Stock In';
      case MovementType.stockOut:
        return 'Stock Out';
      case MovementType.adjustment:
        return 'Adjustment';
      case MovementType.transfer:
        return 'Transfer';
      case MovementType.purchase:
        return 'Purchase';
      case MovementType.sale:
        return 'Sale';
      case MovementType.return_:
        return 'Return';
    }
  }
}

class StockMovement extends Equatable {
  final int movementId;
  final int productId;
  final MovementType movementType;
  final int quantity;
  final double? unitCost;
  final String? referenceType;
  final int? referenceId;
  final String? notes;
  final DateTime createdAt;
  final int createdBy;

  const StockMovement({
    required this.movementId,
    required this.productId,
    required this.movementType,
    required this.quantity,
    this.unitCost,
    this.referenceType,
    this.referenceId,
    this.notes,
    required this.createdAt,
    required this.createdBy,
  });

  @override
  List<Object?> get props => [
        movementId,
        productId,
        movementType,
        quantity,
        unitCost,
        referenceType,
        referenceId,
        notes,
        createdAt,
        createdBy,
      ];
}

enum PurchaseOrderStatus {
  draft('draft'),
  pending('pending'),
  approved('approved'),
  rejected('rejected'),
  received('received'),
  completed('completed'),
  cancelled('cancelled');

  const PurchaseOrderStatus(this.value);
  final String value;

  static PurchaseOrderStatus fromString(String value) {
    switch (value.toLowerCase()) {
      case 'draft':
        return PurchaseOrderStatus.draft;
      case 'pending':
        return PurchaseOrderStatus.pending;
      case 'approved':
        return PurchaseOrderStatus.approved;
      case 'rejected':
        return PurchaseOrderStatus.rejected;
      case 'received':
        return PurchaseOrderStatus.received;
      case 'completed':
        return PurchaseOrderStatus.completed;
      case 'cancelled':
        return PurchaseOrderStatus.cancelled;
      default:
        throw ArgumentError('Invalid purchase order status: $value');
    }
  }

  String get displayName {
    switch (this) {
      case PurchaseOrderStatus.draft:
        return 'Draft';
      case PurchaseOrderStatus.pending:
        return 'Pending';
      case PurchaseOrderStatus.approved:
        return 'Approved';
      case PurchaseOrderStatus.rejected:
        return 'Rejected';
      case PurchaseOrderStatus.received:
        return 'Received';
      case PurchaseOrderStatus.completed:
        return 'Completed';
      case PurchaseOrderStatus.cancelled:
        return 'Cancelled';
    }
  }
}

class PurchaseOrder extends Equatable {
  final int poId;
  final String poNumber;
  final int supplierId;
  final DateTime orderDate;
  final DateTime? expectedDate;
  final PurchaseOrderStatus status;
  final double totalAmount;
  final String? notes;
  final DateTime createdAt;
  final DateTime updatedAt;
  final int createdBy;
  final int? approvedBy;
  final DateTime? approvedAt;

  const PurchaseOrder({
    required this.poId,
    required this.poNumber,
    required this.supplierId,
    required this.orderDate,
    this.expectedDate,
    required this.status,
    required this.totalAmount,
    this.notes,
    required this.createdAt,
    required this.updatedAt,
    required this.createdBy,
    this.approvedBy,
    this.approvedAt,
  });

  @override
  List<Object?> get props => [
        poId,
        poNumber,
        supplierId,
        orderDate,
        expectedDate,
        status,
        totalAmount,
        notes,
        createdAt,
        updatedAt,
        createdBy,
        approvedBy,
        approvedAt,
      ];
}