{% for table in cookiecutter.tables["values"] %}
 CREATE TABLE "{{table["table_snakecase"]}}" (
  {% for col in table["columns"] %}
{# ################################################################################ #}
{# TODO: Inferir el tipo de dato en vez de asumir que esd SERIAL (integer + seq)    #}
{# ################################################################################ #}
    {% if col["primary_key"] %}
  "{{col["column_name_snakecase"]}}" SERIAL PRIMARY KEY{% if not loop.last %},{% endif %}

    {% else %}
  "{{col["column_name_snakecase"]}}" {{ col["data_type"] }}{% if col["not_null"] %} NOT NULL{% endif %}{% if not loop.last %},{% endif %}

    {% endif  %}
  {% endfor %}
 );
{% endfor %}
