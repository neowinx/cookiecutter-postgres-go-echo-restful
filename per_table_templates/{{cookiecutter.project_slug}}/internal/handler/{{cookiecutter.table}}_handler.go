package handler

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"{{cookiecutter.project_slug}}/pkg/db"
	"{{cookiecutter.project_slug}}/internal/models"
)

// {{cookiecutter.table_pascalcase}} represents the {{cookiecutter.table_pascalcase}} model.
type {{cookiecutter.table_pascalcase}} struct {
{% for col in cookiecutter.columns["values"] %}
  {% if col["primary_key"]: %}
    {{col['column_name_uppercase']}} *{{col['go_data_type']}}  `json:"{{col['column_name_snakecase']}},omitempty"`
  {% else %}
    {{col['column_name_pascalcase']}} {{col['go_data_type']}}  `json:"{{col["column_name_snakecase"]}}"`
  {% endif %}
{% endfor %}
}

// List{{cookiecutter.table_pascalcase}}Handler handles the listing of all {{cookiecutter.table_snakecase}}s.
//		@Summary		List all {{cookiecutter.table_snakecase}}s
//		@Description	Lists all {{cookiecutter.table_snakecase}}s in the database
//		@Produce		json
//		@Success		200	{array}		handler.{{cookiecutter.table_pascalcase}}
//		@Failure		500	{object}	models.ErrorResponse
//		@Router			/{{cookiecutter.table_snakecase}}s   [get]
func List{{cookiecutter.table_pascalcase}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()
    {{cookiecutter.table_snakecase}}s, err := db.New(dbpool).List{{cookiecutter.table_pascalcase}}s(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{ Error: "Failed to list {{cookiecutter.table_snakecase}}s"})
		}

		return c.JSON(http.StatusOK, {{cookiecutter.table_snakecase}}s)
	}
}

// Create{{cookiecutter.table_pascalcase}}Handler handles the creation of a new {{cookiecutter.table_snakecase}}.
//	@Summary		Create a new {{cookiecutter.table_snakecase}}
//	@Description	Creates a new {{cookiecutter.table_snakecase}} in the database
//	@Param			{{cookiecutter.table_snakecase}}	body	handler.{{cookiecutter.table_pascalcase}}	true	"{{cookiecutter.table_pascalcase}} information"
//	@Produce		json
//	@Success		201	{object}	handler.{{cookiecutter.table_pascalcase}}
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/{{cookiecutter.table_snakecase}}s   [post]
func Create{{cookiecutter.table_pascalcase}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		var {{cookiecutter.table_snakecase}} {{cookiecutter.table_pascalcase}}

		if err := c.Bind(&{{cookiecutter.table_snakecase}}); err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{ Error: "Invalid request payload"})
		}

		ctx := context.Background()

	{% if cookiecutter.columns["values"]|length == 2 %}
		{% set col = cookiecutter.columns["values"][1] %}
		result, err := db.New(dbpool).Create{{cookiecutter.table_pascalcase}}(ctx, {{cookiecutter.table_snakecase}}.{{ col["column_name_pascalcase"] }})
	{% else %}
		arg := db.Create{{cookiecutter.table_pascalcase}}Params {
			{% for col in cookiecutter.columns["values"] %}
				{% if not col["primary_key"] %}
			{{ col["column_name_sqlc_go"] }}:  {{ col["go_data_type"] }}(*&{{cookiecutter.table_snakecase}}.{{col["column_name_pascalcase"]}}),
				{% endif %}
			{% endfor %}
		}

		result, err := db.New(dbpool).Create{{cookiecutter.table_pascalcase}}(ctx, arg)
	{% endif %}

		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{ Error: "Failed to create {{cookiecutter.table_snakecase}}"})
		}

		return c.JSON(http.StatusCreated, result)
  }
}

// Get{{cookiecutter.table_pascalcase}}Handler retrieves a {{cookiecutter.table_snakecase}} by ID.
//
//	@Summary		Get a {{cookiecutter.table_snakecase}} by ID
//	@Description	Retrieves a {{cookiecutter.table_snakecase}} from the database by its ID
//	@Param			id	path	int	true	"{{cookiecutter.table_pascalcase}} ID"
//	@Produce		json
//	@Success		200	{object}	handler.{{cookiecutter.table_pascalcase}}
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/{{cookiecutter.table_snakecase}}s/{id}   [get]
func Get{{cookiecutter.table_pascalcase}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{ Error: "Invalid ID in the request"})
		}

		ctx := context.Background()
		result, err := db.New(dbpool).Get{{cookiecutter.table_pascalcase}}(ctx, int32(id))
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{ Error: "{{cookiecutter.table_pascalcase}} not found"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{ Error: "Failed to fetch {{cookiecutter.table_snakecase}}"})
		}

		return c.JSON(http.StatusOK, result)
	}
}

// Update{{cookiecutter.table_pascalcase}}Handler updates a {{cookiecutter.table_snakecase}} by ID.
//
//	@Summary		Update a {{cookiecutter.table_snakecase}}
//	@Description	Updates a {{cookiecutter.table_snakecase}} in the database by its ID
//	@Param			id		path	int				true	"{{cookiecutter.table_pascalcase}} ID"
//	@Param			{{cookiecutter.table_snakecase}}	body	handler.{{cookiecutter.table_pascalcase}}	true	"{{cookiecutter.table_pascalcase}} information"
//	@Produce		json
//	@Success		200	{object}	models.ErrorResponse	//	Success		message	can			be	more		descriptive	here	(e.g., "{{cookiecutter.table_pascalcase}} updated successfully")
//	@Failure		400	{object}	models.ErrorResponse	//	Detail		error	messages	for	bad			requests
//	@Failure		404	{object}	models.ErrorResponse	//	Specific	error	message		for	not			found
//	@Failure		500	{object}	models.ErrorResponse	//	Generic		error	message		for	internal	errors
//	@Router			/{{cookiecutter.table_snakecase}}s/{id}   [put]
func Update{{cookiecutter.table_pascalcase}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{ Error: "Invalid ID in the request"})
		}

		var updated{{cookiecutter.table_pascalcase}} {{cookiecutter.table_pascalcase}}
		if err := c.Bind(&updated{{cookiecutter.table_pascalcase}}); err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{ Error: "Invalid request payload"})
		}

		ctx := context.Background()
		arg := db.Update{{cookiecutter.table_pascalcase}}Params{
			{% for col in cookiecutter.columns["values"] %}
				{% if col["primary_key"]: %}
			{{col['column_name_uppercase']}}: {{col['go_data_type']}}(id),
				{% else %}
			{{col['column_name_pascalcase']}}: updated{{cookiecutter.table_pascalcase}}.{{col['column_name_pascalcase']}},
				{% endif %}
			{% endfor %}
		}
		err = db.New(dbpool).Update{{cookiecutter.table_pascalcase}}(ctx, arg)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{ Error: "{{cookiecutter.table_pascalcase}} not found"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{ Error: "Failed to update {{cookiecutter.table_snakecase}}"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "{{cookiecutter.table_pascalcase}} updated successfully"})
	}
}

// Delete{{cookiecutter.table_pascalcase}}Handler deletes a {{cookiecutter.table_snakecase}} by ID.
//
//	@Summary		Delete a {{cookiecutter.table_snakecase}}
//	@Description	Deletes a {{cookiecutter.table_snakecase}} from the database by its ID
//	@Param			id	path	int	true	"{{cookiecutter.table_pascalcase}} ID"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	models.ErrorResponse	//	Detail		error	messages	for	bad			requests
//	@Failure		404	{object}	models.ErrorResponse	//	Specific	error	message		for	not			found
//	@Failure		500	{object}	models.ErrorResponse	//	Generic		error	message		for	internal	errors
//	@Router			/{{cookiecutter.table_snakecase}}s/{id}   [delete]
func Delete{{cookiecutter.table_pascalcase}}Handler(dbpool *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{ Error: "Invalid ID in the request"})
		}

		ctx := context.Background()
		err = db.New(dbpool).Delete{{cookiecutter.table_pascalcase}}(ctx, int32(id))
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{ Error: "{{cookiecutter.table_pascalcase}} not found"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{ Error: "Failed to delete {{cookiecutter.table_snakecase}}"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
