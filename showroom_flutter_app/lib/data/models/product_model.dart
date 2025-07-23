import 'package:json_annotation/json_annotation.dart';
import 'package:equatable/equatable.dart';

part 'product_model.g.dart';

@JsonSerializable()
class ProductModel extends Equatable {
  @JsonKey(name: 'product_id')
  final int productId;
  @JsonKey(name: 'product_code')
  final String productCode;
  @JsonKey(name: 'product_name')
  final String productName;
  final String? description;
  @JsonKey(name: 'brand_id')
  final int brandId;
  @JsonKey(name: 'category_id')
  final int categoryId;
  @JsonKey(name: 'unit_measure')
  final String unitMeasure;
  @JsonKey(name: 'cost_price')
  final double costPrice;
  @JsonKey(name: 'selling_price')
  final double sellingPrice;
  @JsonKey(name: 'markup_percentage')
  final double markupPercentage;
  @JsonKey(name: 'stock_quantity')
  final int stockQuantity;
  @JsonKey(name: 'min_stock_level')
  final int minStockLevel;
  @JsonKey(name: 'max_stock_level')
  final int maxStockLevel;
  @JsonKey(name: 'location_rack')
  final String? locationRack;
  final String? barcode;
  final double? weight;
  final String? dimensions;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;
  @JsonKey(name: 'created_by')
  final int createdBy;
  @JsonKey(name: 'is_active')
  final bool isActive;
  @JsonKey(name: 'product_image')
  final String? productImage;
  final String? notes;

  const ProductModel({
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

  factory ProductModel.fromJson(Map<String, dynamic> json) => _$ProductModelFromJson(json);
  Map<String, dynamic> toJson() => _$ProductModelToJson(this);

  bool get isLowStock => stockQuantity <= minStockLevel;

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

@JsonSerializable()
class StockMovementModel extends Equatable {
  @JsonKey(name: 'movement_id')
  final int movementId;
  @JsonKey(name: 'product_id')
  final int productId;
  @JsonKey(name: 'movement_type')
  final String movementType;
  final int quantity;
  @JsonKey(name: 'unit_cost')
  final double? unitCost;
  @JsonKey(name: 'reference_type')
  final String? referenceType;
  @JsonKey(name: 'reference_id')
  final int? referenceId;
  final String? notes;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'created_by')
  final int createdBy;

  const StockMovementModel({
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

  factory StockMovementModel.fromJson(Map<String, dynamic> json) => _$StockMovementModelFromJson(json);
  Map<String, dynamic> toJson() => _$StockMovementModelToJson(this);

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

@JsonSerializable()
class PurchaseOrderModel extends Equatable {
  @JsonKey(name: 'po_id')
  final int poId;
  @JsonKey(name: 'po_number')
  final String poNumber;
  @JsonKey(name: 'supplier_id')
  final int supplierId;
  @JsonKey(name: 'order_date')
  final DateTime orderDate;
  @JsonKey(name: 'expected_date')
  final DateTime? expectedDate;
  final String status;
  @JsonKey(name: 'total_amount')
  final double totalAmount;
  final String? notes;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;
  @JsonKey(name: 'created_by')
  final int createdBy;
  @JsonKey(name: 'approved_by')
  final int? approvedBy;
  @JsonKey(name: 'approved_at')
  final DateTime? approvedAt;

  const PurchaseOrderModel({
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

  factory PurchaseOrderModel.fromJson(Map<String, dynamic> json) => _$PurchaseOrderModelFromJson(json);
  Map<String, dynamic> toJson() => _$PurchaseOrderModelToJson(this);

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