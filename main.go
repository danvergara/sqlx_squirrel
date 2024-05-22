package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

type Rental struct {
	RentalID    int            `db:"rental_id"`
	RentalDate  time.Time      `db:"rental_date"`
	InventoryID int            `db:"inventory_id"`
	CustomerID  int            `db:"customer_id"`
	ReturnDate  *time.Time     `db:"return_date"`
	StaffID     int            `db:"staff_id"`
	LastUpdate  time.Time      `db:"last_update"`
	FirstName   string         `db:"first_name"`
	LastName    string         `db:"last_name"`
	Email       sql.NullString `db:"email"`
}

type RentalQueryParam struct {
	District   *string
	PostalCode *int
	RentalDate *time.Time
	Returned   bool
}

type Film struct {
	FilmID             int            `db:"film_id"`
	Title              string         `db:"title"`
	Description        sql.NullString `db:"description"`
	ReleaseYear        sql.NullInt32  `db:"release_year"`
	LanguageID         int            `db:"language_id"`
	OriginalLanguageID sql.NullInt32  `db:"original_language_id"`
	RentalDuration     int            `db:"rental_duration"`
	RentalRate         float32        `db:"rental_rate"`
	Length             sql.NullInt32  `db:"length"`
	ReplacementCost    float32        `db:"replacement_cost"`
	Rating             sql.NullString `db:"rating"`
	FirstName          string         `db:"first_name"`
	LastName           string         `db:"last_name"`
	LastUpdate         *time.Time     `db:"last_update"`
}

type FilmQueryParam struct {
	ActorLastName   *string
	ActorFirsttName *string
	Rating          []string
}

func ListRentals(db *sqlx.DB, qp RentalQueryParam) ([]Rental, error) { // HL5
	var rentals []Rental
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select(
		"c.email",
		"c.first_name",
		"c.last_name",
		"r.rental_id",
		"r.rental_date",
		"r.return_date",
	).
		From("customer As c").
		Join("rental AS r ON c.customer_id = r.customer_id").
		Join("address AS a ON c.address_id = a.address_id").
		Join("city AS ct ON a.city_id = ct.city_id")

	if qp.District != nil {
		query = query.Where(sq.Eq{"a.district": *qp.District})
	}

	if qp.PostalCode != nil {
		query = query.Where(sq.Eq{"a.postal_code": *qp.PostalCode})
	}

	if qp.RentalDate != nil {
		query = query.Where(sq.Lt{"r.rental_date": *qp.RentalDate})
	}

	if qp.Returned {
		query = query.Where(sq.NotEq{"r.return_date": nil})
	} else {
		query = query.Where(sq.Eq{"r.return_date": nil})
	}

	query = query.Limit(100).OrderBy("r.rental_date desc")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if err := db.Select(&rentals, sql, args...); err != nil {
		return nil, err
	}

	return rentals, nil
}

func ListFilms(db *sqlx.DB, qp FilmQueryParam) ([]Film, error) { // HL6
	var films []Film

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select(
		"f.title",
		"f.rating",
	).
		From("film AS f")

	if qp.ActorFirsttName != nil || qp.ActorLastName != nil {
		query = query.
			Join("film_actor AS fa ON f.film_id = fa.film_id").
			Join("actor AS a ON fa.actor_id = a.actor_id")
	}

	if qp.ActorFirsttName != nil {
		query = query.Where(sq.Eq{"a.first_name": *qp.ActorFirsttName})
	}

	if qp.ActorLastName != nil {
		query = query.Where(sq.Eq{"a.last_name": *qp.ActorLastName})
	}

	if len(qp.Rating) != 0 {
		query = query.Where(sq.Eq{"f.rating": qp.Rating})
	}

	query = query.Limit(100)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if err := db.Select(&films, sql, args...); err != nil {
		return nil, err
	}

	return films, nil
}

func printRentalResultSet(rentals []Rental) {
	w := tabwriter.NewWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "First Name\tLast Name\tEmail\tRental ID\tRental Date\t Return Date\t")

	for _, r := range rentals {
		fmt.Fprintln(
			w,
			r.FirstName,
			"\t",
			r.LastName,
			"\t",
			r.Email,
			"\t",
			r.RentalID,
			"\t",
			r.RentalDate,
			"\t",
			r.ReturnDate,
			"\t",
		)
	}

	w.Flush()
}

func printFilmsResultSet(films []Film) {
	w := tabwriter.NewWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "Title\tRating\t")

	for _, f := range films {
		fmt.Fprintln(
			w,
			f.Title,
			"\t",
			f.Rating,
			"\t",
		)
	}

	w.Flush()
}

func main() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	district := "California"
	postalCode := 52137
	rentalDate := time.Now().UTC()

	rqp := RentalQueryParam{
		District:   &district,
		PostalCode: &postalCode,
		RentalDate: &rentalDate,
		Returned:   false,
	}

	rentals, err := ListRentals(db, rqp)
	if err != nil {
		log.Fatalln(err)
	}

	printRentalResultSet(rentals)

	fmt.Println(
		"---------------------------------------------------------------------------------------------------------",
	)

	firstName := "CATE"
	lastName := "MCQUEEN"
	fqp := FilmQueryParam{
		ActorFirsttName: &firstName,
		ActorLastName:   &lastName,
		Rating:          []string{"G", "PG"},
	}
	films, err := ListFilms(
		db,
		fqp,
	)
	if err != nil {
		log.Fatalln(err)
	}

	printFilmsResultSet(films)
}
