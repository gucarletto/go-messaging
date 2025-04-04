package usecase

import "github.com/gucarletto/go-messaging/internal/entity"

type CreateProductInputDTO struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateProductOutputDTO struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateProductUseCase struct {
	ProductRepository entity.ProductRepository
}

func NewCreateProductUseCase(productRepository entity.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{
		ProductRepository: productRepository,
	}
}

func (u *CreateProductUseCase) Execute(input CreateProductInputDTO) (*CreateProductOutputDTO, error) {
	product := entity.NewProduct(input.Name, input.Price)
	err := u.ProductRepository.Create(product)
	if err != nil {
		return nil, err
	}

	output := &CreateProductOutputDTO{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}

	return output, nil
}
