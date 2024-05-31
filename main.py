#!/usr/bin/env python3

import psycopg2
import inquirer
import json
import os
from caseconverter import camelcase, pascalcase

from cookiecutter.main import cookiecutter

HOOK_PATH=os.path.dirname(__file__)

def get_postgresql_tables(host, port, db, user, password, schema):
    connection = psycopg2.connect(
        host=host,
        port=port,
        database=db,
        user=user,
        password=password
    )
    cursor = connection.cursor()
    cursor.execute(f"""
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = '{schema}'
    """)
    tables = [row[0] for row in cursor.fetchall()]
    cursor.close()
    connection.close()
    return tables

def get_columns_info(host, port, db, user, password, schema, table):
    connection = psycopg2.connect(
        host=host,
        port=port,
        database=db,
        user=user,
        password=password
    )
    cursor = connection.cursor()
    cursor.execute(f"""
        SELECT
          a.attname                                       as "column_name",
          pg_catalog.format_type(a.atttypid, a.atttypmod) as "data_type",
          a.attnotnull                                    as "not_null"
        FROM
          pg_catalog.pg_attribute a
        WHERE
          a.attnum > 0
          AND NOT a.attisdropped
          AND a.attrelid = (
            SELECT c.oid
            FROM pg_catalog.pg_class c
              LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
            WHERE c.relname ~ '^({table})$' AND n.nspname = '{schema}');
    """)
    columns = [{ "column_name": row[0], "data_type": row[1], "not_null": row[2] } for row in cursor.fetchall()]
    cursor.close()
    connection.close()
    return columns


IGNORED_FIELDS = [ "project_slug", "author", "_copy_without_render", "selected_tables", "table" ] 


def main():
    # Load cookiecutter.json to get initial values
    COOKIE_CUTTER_JSON = os.path.join(HOOK_PATH, 'cookiecutter.json')

    with open(COOKIE_CUTTER_JSON) as f:
        context: dict = json.load(f)

    questions = []

    for key in context.keys():
        if key not in IGNORED_FIELDS:
            questions.append(inquirer.Text(f'{key}', default=context.get(key), message=f'Ingrese el {key}'))

    answers = inquirer.prompt(questions)

    if not answers:
        print("No answers returned. Exiting.")
        exit(1)

    host = answers["postgresql_host"]
    port = answers["postgresql_port"]
    db = answers["postgresql_db"]
    user = answers["postgresql_user"]
    password = answers["postgresql_password"]
    schema = answers["postgresql_schema"]

    tables = get_postgresql_tables(host, port, db, user, password, schema)

    questions = [
        inquirer.Checkbox(
            'selected_tables',
            message="Select the tables to include in your project",
            choices=tables
        ),
    ]

    answers = inquirer.prompt(questions)

    if not answers:
        print("No answers returned. Exiting.")
        exit(2)
    
    if answers['selected_tables']:
        context["selected_tables"]["values"] = answers['selected_tables']

    # Write the updated context back to cookiecutter.json
    with open(COOKIE_CUTTER_JSON, 'w') as f:
        json.dump(context, f, indent=4)

    print("Processing template...")

    cookiecutter(template=HOOK_PATH, no_input=True, overwrite_if_exists=True)

    print("generated. Processing tables...")

    # repeat the template processing for each table but pass the "skip_if_file_exists" attribute
    # in order to generate the individual files

    for table in context["selected_tables"]["values"]:
        columns = get_columns_info(host, port, db, user, password, schema, table)
        cookiecutter(template=HOOK_PATH, no_input=True, overwrite_if_exists=True, skip_if_file_exists=True, extra_context={
            "table": table,
            "table_camel_case": camelcase(table),
            "table_pascal_case": pascalcase(table),
            "columns": columns
        })

    print("done.")

if __name__ == "__main__":
    main()
