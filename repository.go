package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/morenovic/go-storage/models"
)

type Repository interface {
	GetAll() ([]models.Product, error)
	GetOne(id int) (models.Product, error)
	Store(models.Product) (models.Product, error)
	GetByName(name string) (models.Product, error)
	UpdateWithContext(ctx context.Context, p models.Product) (models.Product, error)
	Delete(id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const (
	Q_GET_ONE     = "SELECT id, name, type, count, price FROM products WHERE id=?;"
	Q_STORE       = "INSERT INTO products(name, type, count, price) VALUES( ?, ?, ?, ? )"
	Q_GET_BY_NAME = "SELECT id, name, type, count, price FROM products WHERE name=?;"
	Q_GET_ALL     = "SELECT id, name, type, count, price FROM products"
	Q_UPDATE      = "UPDATE products SET name = ?, type = ?, count = ?, price = ? WHERE id = ?"
	Q_DELETE      = "DELETE FROM products WHERE id=?"
)

func (r *repository) GetOne(id int) (models.Product, error) {
	row := r.db.QueryRow(Q_GET_ONE, id)
	product := models.Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price)
	if err != nil {
		return models.Product{}, nil
	}
	fmt.Println(product)
	return product, nil
}

func (r *repository) Store(product models.Product) (models.Product, error) {
	stmt, err := r.db.Prepare(Q_STORE) // se prepara el SQL
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	var result sql.Result
	result, err = stmt.Exec(product.Name, product.Type, product.Count, product.Price) // retorna un sql.Result y un error
	if err != nil {
		return models.Product{}, err
	}
	insertedId, _ := result.LastInsertId() // del sql.Resul devuelto en la ejecución obtenemos el Id insertado
	product.ID = int(insertedId)
	return product, nil
}

func (r *repository) GetByName(name string) (models.Product, error) {
	product := models.Product{}

	row := r.db.QueryRow(Q_GET_BY_NAME, name)
	err := row.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price)

	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (r *repository) GetAll() ([]models.Product, error) {
	var products []models.Product
	rows, err := r.db.Query(Q_GET_ALL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// se recorren todas las filas
	for rows.Next() {
		// por cada fila se obtiene un objeto del tipo Product
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price); err != nil {
			log.Fatal(err)
			return nil, err
		}
		//se añade el objeto obtenido al slice products
		products = append(products, product)
	}
	return products, nil
}

func (r *repository) UpdateWithContext(ctx context.Context, product models.Product) (models.Product, error) { // se inicializa la base
	stmt, err := r.db.Prepare(Q_UPDATE) //
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // se cierra la sentencia al terminar. Si quedan abiertas se genera consumos de memoria
	_, err = stmt.ExecContext(ctx, product.Name, product.Type, product.Count, product.Price, product.ID)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *repository) Delete(id int) error {
	stmt, err := r.db.Prepare(Q_DELETE)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return errors.New("no Encontrado")
	}

	return nil
}
