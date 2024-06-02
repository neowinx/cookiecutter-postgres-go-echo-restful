{% for table in cookiecutter.tables["values"] %}

-- name: Get{{table["table_pascalcase"]}} :one
SELECT * FROM "{{table["table_snakecase"]}}"
WHERE "id" = $1 LIMIT 1;

-- name: List{{table["table_pascalcase"]}}s :many
SELECT * FROM "{{table["table_snakecase"]}}"
ORDER BY "{{table["columns"][0]["column_name_snakecase"]}}";

-- name: Create{{table["table_pascalcase"]}}WithID :one
INSERT INTO "{{table["table_snakecase"]}}" (
{% for col in table["columns"] %}
  {% if loop.last %}
  "{{col["column_name_snakecase"]}}"
  {% else %}
  "{{col["column_name_snakecase"]}}",
  {% endif %}
{% endfor %}
) VALUES (
{% for col in table["columns"] %}
  {% if loop.last %}
  ${{ loop.index }}
  {% else %}
  ${{ loop.index }},
  {% endif %}
{% endfor %}
)
RETURNING *;

-- name: Create{{table["table_pascalcase"]}} :one
INSERT INTO "{{table["table_snakecase"]}}" (
{% for col in table["columns"] %}
  {% if not col["primary_key"] %}
    {% if loop.last %}
  "{{col["column_name_snakecase"]}}"
    {% else %}
  "{{col["column_name_snakecase"]}}",
    {% endif %}
  {% endif %}
{% endfor %}
) VALUES (
{% for col in table["columns"] %}
  {% if not col["primary_key"] %}
    {% if loop.last %}
  ${{ loop.index0 }}
    {% else %}
  ${{ loop.index0 }},
    {% endif %}
  {% endif %}
{% endfor %}
)
RETURNING *;

-- name: Update{{table["table_pascalcase"]}} :exec
UPDATE "{{table["table_snakecase"]}}"
  SET
{% for col in table["columns"] %}
  {% if not col["primary_key"] %}
    {% if loop.last %}
  "{{col["column_name_snakecase"]}}" = ${{loop.index}}
    {% else %}
  "{{col["column_name_snakecase"]}}" = ${{loop.index}},
    {% endif %}
  {% endif %}
{% endfor %}
{% for col in table["columns"] %}
  {% if col["primary_key"] %}
WHERE {{col["column_name_snakecase"]}} = ${{loop.index}};
  {% endif %}
{% endfor %}

-- name: Delete{{table["table_pascalcase"]}} :exec
DELETE FROM "{{table["table_snakecase"]}}"
{% for col in table["columns"] %}
  {% if col["primary_key"] %}
WHERE "{{col["column_name_snakecase"]}}" = ${{loop.index}};
  {% endif %}
{% endfor %}

{% endfor %}
