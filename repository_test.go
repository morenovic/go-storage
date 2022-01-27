package product

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/morenovic/go-storage/internal/product/util"
	"github.com/morenovic/go-storage/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateOK(t *testing.T) {
	db, _ := sql.Open("mysql", "meli_sprint_user:Meli_Sprint#123@/storage")

	newProduct := models.Product{
		Name:  "crema",
		Type:  "manos",
		Count: 10,
		Price: 40,
	}

	repo := NewRepo(db)
	res, _ := repo.Store(newProduct)

	assert.Equal(t, newProduct.Name, res.Name)
	assert.Equal(t, newProduct.Type, res.Type)
	assert.Equal(t, newProduct.Count, res.Count)
	assert.Equal(t, newProduct.Price, res.Price)
}

func TestGetAllOK(t *testing.T) {
	db, _ := sql.Open("mysql", "meli_sprint_user:Meli_Sprint#123@/storage")

	repo := NewRepo(db)
	res, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, res, 4)
}

func TestUpdateWithContext(t *testing.T) {
	db, _ := sql.Open("mysql", "meli_sprint_user:Meli_Sprint#123@/storage")
	// definimos un Product cuyo nombre sea igual al registro de la DB
	product := models.Product{
		ID:    1,
		Name:  "yogurt",
		Type:  "entero",
		Count: 20,
		Price: 20.5,
	}
	myRepo := NewRepo(db)
	// se define un context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	productResult, _ := myRepo.UpdateWithContext(ctx, product)
	assert.Equal(t, product.Name, productResult.Name)
}

func TestGetByName(t *testing.T) {
	db, _ := sql.Open("mysql", "meli_sprint_user:Meli_Sprint#123@/storage")
	name := "crema"

	repo := NewRepo(db)
	res, _ := repo.GetByName(name)

	assert.Equal(t, name, res.Name)
	assert.IsType(t, models.Product{}, res)
}

func TestGetOne(t *testing.T) {
	db, _ := sql.Open("mysql", "meli_sprint_user:Meli_Sprint#123@/storage")

	repo := NewRepo(db)
	res, _ := repo.GetOne(1)

	assert.Equal(t, "yogurt", res.Name)

}

func Test_sqlRepository_Store(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)

	newProduct := models.Product{
		Name:  "crema",
		Type:  "manos",
		Count: 10,
		Price: 40,
	}
	repo := NewRepo(db)

	result, err := repo.Store(newProduct)
	assert.NoError(t, err)
	getResult, err := repo.GetOne(100)
	assert.NoError(t, err)
	fmt.Println(result.ID)
	getResult, err = repo.GetOne(int(result.ID))
	assert.NoError(t, err)
	assert.NotNil(t, getResult)
	assert.Equal(t, "manos", getResult.Type)
}

func Test_sqlRepository_Update(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)

	expectedResult := models.Product{
		ID:    1,
		Name:  "crema",
		Type:  "manos",
		Count: 30,
		Price: 40,
	}
	repo := NewRepo(db)

	result, err := repo.UpdateWithContext(context.Background(), expectedResult)
	assert.NoError(t, err)
	assert.NotEqual(t, expectedResult, result)
}

func Test_sqlRepository_Delete(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)

	IdToFind := 14

	repo := NewRepo(db)

	oldGetallRes, _ := repo.GetAll()

	err = repo.Delete(IdToFind)

	assert.NoError(t, err)

	getOneResult, err := repo.GetOne(IdToFind)
	assert.NoError(t, err)
	assert.Zero(t, getOneResult.ID)
	getAllResult, _ := repo.GetAll()
	assert.NotEqual(t, len(oldGetallRes), len(getAllResult))

}

func Test_sqlRepository_Store_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectPrepare("^INSERT INTO products")
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products")).WillReturnResult(sqlmock.NewResult(1, 1))
	columns := []string{"id", "name", "type", "count", "price"}
	rows := sqlmock.NewRows(columns)

	rows.AddRow(1, "32332", "44", 23, 12)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, type, count, price FROM products")).WillReturnRows(rows)
	repo := NewRepo(db)
	newProduct := models.Product{
		Name:  "crema",
		Type:  "manos",
		Count: 10,
		Price: 40,
	}

	getResult, err := repo.GetOne(15)
	assert.NoError(t, err)
	_, err = repo.Store(newProduct)
	assert.NoError(t, err)
	getResult2, err := repo.GetAll()
	fmt.Println(getResult2, getResult)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
