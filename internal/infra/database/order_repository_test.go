package database

import (
	"database/sql"
	"testing"

	"github.com/kenesparta/fullcycle-clean-architecture/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("SELECT id, price, tax, final_price FROM orders WHERE id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestGivenOrders_WhenList_ThenShouldGetAllOrders() {
	repo := NewOrderRepository(suite.Db)

	order1, _ := entity.NewOrder("123", 10.0, 2.0)
	_ = order1.CalculateFinalPrice()
	_ = repo.Save(order1)

	order2, _ := entity.NewOrder("456", 1.8, 2.6)
	_ = order2.CalculateFinalPrice()
	_ = repo.Save(order2)

	orderList, listErr := repo.List()
	suite.NoError(listErr)

	var orderResult entity.Order
	err := suite.Db.QueryRow("SELECT * FROM orders").
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.EqualValues(
		orderList,
		[]entity.Order{
			{
				ID:         order1.ID,
				Price:      order1.Price,
				Tax:        order1.Tax,
				FinalPrice: order1.FinalPrice,
			},
			{
				ID:         order2.ID,
				Price:      order2.Price,
				Tax:        order2.Tax,
				FinalPrice: order2.FinalPrice,
			},
		},
	)
}
