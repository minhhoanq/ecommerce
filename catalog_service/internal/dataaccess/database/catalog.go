package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/minhhoanq/ecommerce/common/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Interface definition
type CatalogDataAccessor interface {
	CreateProduct(ctx context.Context, arg *CreateProductParams) (*CreateProductResponse, error)
	ListProducts(ctx context.Context, arg *ListProductRequest) (*ListProductResponse, error)
	CreateCart(ctx context.Context, arg *CreateCartRequest) (*CreateCartResponse, error)
	AddToCartItem(ctx context.Context, arg *AddToCartItemRequest) (*AddToCartItemResponse, error)
	GetSKU(ctx context.Context, arg *GetSKURequest) (*GetSKUResponse, error)
	GetInventorySKU(ctx context.Context, arg *GetInventorySKURequest) (*GetInventorySKUResponse, error)
	UpdateInventorySKU(ctx context.Context, arg *UpdateInventorySKURequest) (*UpdateInventorySKUResponse, error)
}

type catalogDataAccessor struct {
	database Database
	l        logger.Interface
}

func NewCatalogDataAccessor(database Database, l logger.Interface) CatalogDataAccessor {
	l.Info("initializing catalog data accessor")
	return &catalogDataAccessor{
		database: database,
		l:        l,
	}
}

func (c *catalogDataAccessor) CreateProduct(ctx context.Context, arg *CreateProductParams) (*CreateProductResponse, error) {
	c.l.Info("starting product creation transaction", zap.String("product_name", arg.Name))

	var response *CreateProductResponse

	// Start transaction
	err := c.database.Transaction(func(tx *gorm.DB) error {
		// Create product
		product := Product{
			Name:        arg.Name,
			Description: arg.Description,
			Image:       arg.Image,
			CategoryID:  arg.CategoryID,
			BrandID:     arg.BrandID,
		}

		if err := tx.Create(&product).Error; err != nil {
			c.l.Error("failed to create product",
				zap.String("product_name", arg.Name),
				zap.Error(err))
			return err
		}

		var skuDetails []SKUDetail
		// Create SKUs and related entities
		for _, skuParam := range arg.SKUs {
			skuDetail, err := c.createSKUWithDetails(tx, &product, &skuParam)
			if err != nil {
				return err
			}
			skuDetails = append(skuDetails, *skuDetail)
		}

		response = &CreateProductResponse{
			ProductDetail: ProductDetail{
				Product: product,
				SKUs:    skuDetails,
			},
		}

		return nil
	})

	if err != nil {
		c.l.Error("transaction failed", zap.Error(err))
		return nil, err
	}

	c.l.Info("product creation completed successfully",
		zap.String("product_id", response.ProductDetail.Product.ID.String()))

	return response, nil
}

func (c *catalogDataAccessor) createSKUWithDetails(
	tx *gorm.DB,
	product *Product,
	skuParam *CreateSKUParams,
) (*SKUDetail, error) {
	// Create SKU
	sku := SKU{
		ProductID: product.ID,
		Name:      skuParam.Name,
		Slug:      skuParam.Slug,
	}

	if err := tx.Create(&sku).Error; err != nil {
		c.l.Error("failed to create SKU",
			zap.String("sku_name", skuParam.Name),
			zap.Error(err))
		return nil, err
	}

	// Create price
	price := Price{
		SkuID:         sku.ID,
		OriginalPrice: skuParam.Price.OriginalPrice,
		EffectiveDate: skuParam.Price.EffectiveDate,
		Active:        true,
	}

	if err := tx.Create(&price).Error; err != nil {
		c.l.Error("failed to create price",
			zap.String("sku_id", sku.ID.String()),
			zap.Error(err))
		return nil, err
	}

	// Create inventory
	inventory := Inventory{
		SkuID: sku.ID,
		Stock: skuParam.Inventory.Stock,
	}

	if err := tx.Create(&inventory).Error; err != nil {
		c.l.Error("failed to create inventory",
			zap.String("sku_id", sku.ID.String()),
			zap.Error(err))
		return nil, err
	}

	// Create SKU attributes
	var skuAttributes []SKUAttribute
	for _, attrParam := range skuParam.Attributes {
		skuAttr := SKUAttribute{
			SkuID:       sku.ID,
			AttributeID: attrParam.AttributeID,
			Value:       attrParam.Value,
		}

		if err := tx.Create(&skuAttr).Error; err != nil {
			c.l.Error("failed to create SKU attribute",
				zap.String("sku_id", sku.ID.String()),
				zap.Int32("attribute_id", attrParam.AttributeID),
				zap.Error(err))
			return nil, err
		}

		skuAttributes = append(skuAttributes, skuAttr)
	}

	return &SKUDetail{
		SKU: sku,
	}, nil
}

func (c *catalogDataAccessor) ListProducts(ctx context.Context, arg *ListProductRequest) (*ListProductResponse, error) {
	var products []ProductDetail
	var rawResult []byte

	offset := (arg.Page - 1) * arg.PageSize
	fmt.Println("page size", arg.PageSize)

	query := `
		WITH limited_products AS (
			SELECT * FROM products
			ORDER BY created_at DESC
			LIMIT ?
			OFFSET ?
		)
		SELECT json_agg(
			json_build_object(
				'product', json_build_object(
					'id', p.id,
					'name', p.name,
					'description', p.description,
					'image', p.image,
					'category_id', p.category_id,
					'brand_id', p.brand_id,
					'created_at', to_char(p.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),
					'updated_at', to_char(p.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')
				),
				'skus', (
					SELECT json_agg(
						json_build_object(
							'sku', json_build_object(
								'id', s.id,
								'product_id', s.product_id,
								'name', s.name,
								'slug', s.slug,
								'price', json_build_object(
									'id', pr.id,
									'sku_id', pr.sku_id,
									'original_price', pr.original_price,
									'effective_date', to_char(pr.effective_date, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),
									'active', pr.active
								),
								'inventory', json_build_object(
									'id', i.id,
									'sku_id', i.sku_id,
									'stock', i.stock,
									'reservations', i.reservations
								),
								'attributes', (
									SELECT json_agg(
										json_build_object(
											'attribute_id', sa.attribute_id,
											'value', sa.value
										)
									)
									FROM sku_attributes sa
									WHERE sa.sku_id = s.id
								)
							)
						)
					)
					FROM skus s
					INNER JOIN prices pr ON s.id = pr.sku_id
					INNER JOIN inventories i ON s.id = i.sku_id
					WHERE s.product_id = p.id
				)
			)
		) AS result
		FROM limited_products p;`

	if err := c.database.Raw(query, arg.PageSize, offset).Row().Scan(&rawResult); err != nil {
		c.l.Error("failed to list products", zap.Error(err))
		return &ListProductResponse{
			Products: []ProductDetail{},
		}, err
	}

	if rawResult == nil {
		return &ListProductResponse{
			Products: []ProductDetail{},
		}, nil
	}

	if err := json.Unmarshal(rawResult, &products); err != nil {
		c.l.Error("failed to unmarshal products", zap.Error(err))
		return nil, err
	}

	return &ListProductResponse{
		Products: products,
	}, nil
}

func (c *catalogDataAccessor) CreateCart(ctx context.Context, arg *CreateCartRequest) (*CreateCartResponse, error) {
	c.l.Info("create cart request", zap.String("user_id", arg.UserID.String()))

	cart := &Cart{
		UserID: arg.UserID,
	}

	if err := c.database.WithContext(ctx).Create(&cart).Error; err != nil {
		c.l.Error("failed to create cart", zap.Error(err))
		return nil, err
	}

	return &CreateCartResponse{
		CartID: cart.ID,
		UserID: arg.UserID,
	}, nil
}

func (c *catalogDataAccessor) AddToCartItem(ctx context.Context, arg *AddToCartItemRequest) (*AddToCartItemResponse, error) {
	c.l.Info("add to cart item request",
		zap.String("user_id", arg.CartID.String()),
		zap.String("sku_id", arg.SkuID.String()),
		zap.Int32("quantity", arg.Quantity))

	cart := &Cart{}

	if err := c.database.WithContext(ctx).Where("id = ?", arg.CartID).First(&cart).Error; err != nil {
		c.l.Error("cart not found", zap.Error(err))
		return nil, err
	}

	sku := &SKU{}
	if err := c.database.WithContext(ctx).Select("id").Where("id = ?", arg.SkuID).First(&sku).Error; err != nil {
		c.l.Error("sku not found", zap.Error(err))
		return nil, err
	}

	// Check if the SKU is already in the cart
	var cartItems []CartItem
	if err := c.database.WithContext(ctx).Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.l.Error("failed to check if SKU is already in cart", zap.Error(err))
			return nil, err
		}
	}

	if len(cartItems) > 0 {
		for _, cartItem := range cartItems {
			if cartItem.SkuID == arg.SkuID {
				c.l.Info("SKU already in cart", zap.String("sku_id", arg.SkuID.String()))
				return nil, fmt.Errorf("SKU already in cart")
			}
		}
	}

	cartItem := &CartItem{
		CartID:   cart.ID,
		SkuID:    arg.SkuID,
		Quantity: arg.Quantity,
	}

	if err := c.database.WithContext(ctx).Create(&cartItem).Error; err != nil {
		c.l.Error("failed to add to cart item", zap.Error(err))
		return nil, err
	}

	cartItems = append(cartItems, *cartItem)

	return &AddToCartItemResponse{
		Cart:     Cart{ID: cartItem.ID},
		CartItem: cartItems,
	}, nil
}

func (c *catalogDataAccessor) GetSKU(ctx context.Context, arg *GetSKURequest) (*GetSKUResponse, error) {
	var sku SKU
	fmt.Println("sku id: ", arg.SkuID)
	if err := c.database.Raw(`select * from skus where id = ?`, arg.SkuID).Scan(&sku).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("sku not found")
		}
		return nil, err
	}

	// Check if the SKU is not found and return an empty response
	if sku.ID == uuid.Nil {
		return nil, fmt.Errorf("sku not found")
	}

	return &GetSKUResponse{
		SKU: sku,
	}, nil
}

func (c *catalogDataAccessor) GetInventorySKU(ctx context.Context, arg *GetInventorySKURequest) (*GetInventorySKUResponse, error) {
	var inventory Inventory
	fmt.Println("sku id: ", arg.SkuID)
	if err := c.database.WithContext(ctx).Raw(`select * from inventories where sku_id = ?`, arg.SkuID).Scan(&inventory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("inventory not found")
		}
		return nil, err
	}

	return &GetInventorySKUResponse{
		Inventory: inventory,
	}, nil
}

func (c *catalogDataAccessor) UpdateInventorySKU(ctx context.Context, arg *UpdateInventorySKURequest) (*UpdateInventorySKUResponse, error) {
	query := `
        UPDATE inventories 
        SET stock = COALESCE(stock - COALESCE($1, 0), stock)
        WHERE sku_id = $2
          AND stock >= COALESCE($1, 0)
        RETURNING *
    `

	var inventory Inventory
	if err := c.database.WithContext(ctx).Raw(query, arg.Quantity, arg.SkuID).Scan(&inventory).Error; err != nil {
		return nil, fmt.Errorf("failed to update inventory")
	}

	return &UpdateInventorySKUResponse{
		Inventory: inventory,
	}, nil
}
