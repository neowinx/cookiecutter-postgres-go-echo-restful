package handler

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"{{cookiecutter.project_slug}}/pkg/db"
)

// {{cookiecutter.table_pascal_case}} represents the {{cookiecutter.table_pascal_case}} model.
type {{cookiecutter.table_pascal_case}} struct {
{{cookiecutter.columns}}
{% for col in cookiecutter.columns["values"] %}
  {{col['column_name']}}   *{{col['data_type']}}
{% endfor %}
	ID   *int   `json:"id,omitempty"`
	Name string `json:"name"`
}

//	 List{{cookiecutter.table_pascal_case}}Handler handles the listing of all heroes.
//		@Summary		List all heroes
//		@Description	Lists all heroes in the database
//		@Produce		json
//		@Success		200	{array}		handler.{{cookiecutter.table_pascal_case}}
//		@Failure		500	{object}	map[string]string
//		@Router			/heroes   [get]
func List{{cookiecutter.table_pascal_case}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()
		heroes, err := db.New(dbpool).List{{cookiecutter.table_pascal_case}}s(ctx)
    print(heroes)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list heroes"})
		}

		return c.JSON(http.StatusOK, heroes)
	}
}

// Create{{cookiecutter.table_pascal_case}}Handler handles the creation of a new hero.
//
//	@Summary		Create a new hero
//	@Description	Creates a new hero in the database
//	@Param			hero	body	handler.{{cookiecutter.table_pascal_case}}	true	"{{cookiecutter.table_pascal_case}} information"
//	@Produce		json
//	@Success		201	{object}	handler.{{cookiecutter.table_pascal_case}}
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/heroes   [post]
func Create{{cookiecutter.table_pascal_case}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		var hero {{cookiecutter.table_pascal_case}}
		if err := c.Bind(&hero); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		}

		ctx := context.Background()

    // TODO: try to DRY this part
    if hero.ID == nil {
      result, err := db.New(dbpool).Create{{cookiecutter.table_pascal_case}}(ctx, hero.Name)
      if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create hero"})
      }
      return c.JSON(http.StatusCreated, result)
    } else {
      arg := db.Create{{cookiecutter.table_pascal_case}}WithIDParams{
        ID:   int32(*hero.ID),
        Name: hero.Name,
      }
      result, err := db.New(dbpool).Create{{cookiecutter.table_pascal_case}}WithID(ctx, arg)
      if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create hero"})
      }
      return c.JSON(http.StatusCreated, result)
    }
  }
}
// Get{{cookiecutter.table_pascal_case}}Handler retrieves a hero by ID.
//
//	@Summary		Get a hero by ID
//	@Description	Retrieves a hero from the database by its ID
//	@Param			id	path	int	true	"{{cookiecutter.table_pascal_case}} ID"
//	@Produce		json
//	@Success		200	{object}	handler.{{cookiecutter.table_pascal_case}}
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/heroes/{id}   [get]
func Get{{cookiecutter.table_pascal_case}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID in the request"})
		}

		ctx := context.Background()
		result, err := db.New(dbpool).Get{{cookiecutter.table_pascal_case}}(ctx, int32(id))
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "{{cookiecutter.table_pascal_case}} not found"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch hero"})
		}

		return c.JSON(http.StatusOK, result)
	}
}

// Update{{cookiecutter.table_pascal_case}}Handler updates a hero by ID.
//
//	@Summary		Update a hero
//	@Description	Updates a hero in the database by its ID
//	@Param			id		path	int				true	"{{cookiecutter.table_pascal_case}} ID"
//	@Param			hero	body	handler.{{cookiecutter.table_pascal_case}}	true	"{{cookiecutter.table_pascal_case}} information"
//	@Produce		json
//	@Success		200	{object}	map[string]string	//	Success		message	can			be	more		descriptive	here	(e.g., "{{cookiecutter.table_pascal_case}} updated successfully")
//	@Failure		400	{object}	map[string]string	//	Detail		error	messages	for	bad			requests
//	@Failure		404	{object}	map[string]string	//	Specific	error	message		for	not			found
//	@Failure		500	{object}	map[string]string	//	Generic		error	message		for	internal	errors
//	@Router			/heroes/{id}   [put]
func Update{{cookiecutter.table_pascal_case}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID in the request"})
		}

		var updated{{cookiecutter.table_pascal_case}} {{cookiecutter.table_pascal_case}}
		if err := c.Bind(&updated{{cookiecutter.table_pascal_case}}); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		}

		ctx := context.Background()
		arg := db.Update{{cookiecutter.table_pascal_case}}Params{
			ID:   int32(id),
			Name: updated{{cookiecutter.table_pascal_case}}.Name,
		}
		err = db.New(dbpool).Update{{cookiecutter.table_pascal_case}}(ctx, arg)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "{{cookiecutter.table_pascal_case}} not found"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update hero"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "{{cookiecutter.table_pascal_case}} updated successfully"})
	}
}

// Delete{{cookiecutter.table_pascal_case}}Handler deletes a hero by ID.
//
//	@Summary		Delete a hero
//	@Description	Deletes a hero from the database by its ID
//	@Param			id	path	int	true	"{{cookiecutter.table_pascal_case}} ID"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	map[string]string	//	Detail		error	messages	for	bad			requests
//	@Failure		404	{object}	map[string]string	//	Specific	error	message		for	not			found
//	@Failure		500	{object}	map[string]string	//	Generic		error	message		for	internal	errors
//	@Router			/heroes/{id}   [delete]
func Delete{{cookiecutter.table_pascal_case}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID in the request"})
		}

		ctx := context.Background()
		err = db.New(dbpool).Delete{{cookiecutter.table_pascal_case}}(ctx, int32(id))
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "{{cookiecutter.table_pascal_case}} not found"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete hero"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
