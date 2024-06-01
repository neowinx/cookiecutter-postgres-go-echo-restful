#!/usr/bin/env python3

import psycopg2
import inquirer
import json
import os
from caseconverter import camelcase, pascalcase, snakecase

from cookiecutter.main import cookiecutter

import argparse

parser = argparse.ArgumentParser(description="Cookiecutter template for a RESTFul app using postgresql and echo")
parser.add_argument('-f', '--overwrite-if-exists', help='Overwrite the contents of the output directory')

TMP_PATH=os.path.dirname(__file__)
IGNORED_FIELDS = [ "project_slug", "author", "_copy_without_render", "selected_tables", "table", "table_pascalcase", "table_snakecase", "columns", "_jinja2_env_vars" ]

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


def to_go_data_type(postgres_data_type: str):
   match postgres_data_type:
       case 'integer': return 'int'
       case 'text': return 'string'
       case _: raise Exception("UNMAPPED DATA_TYPE! {postgres_data_type}")

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
        SELECT c.column_name
         FROM information_schema.table_constraints tc
         JOIN information_schema.constraint_column_usage AS ccu USING (constraint_schema, constraint_name)
         JOIN information_schema.columns AS c ON c.table_schema = tc.constraint_schema
           AND tc.table_name = c.table_name AND ccu.column_name = c.column_name
         WHERE constraint_type = 'PRIMARY KEY' and tc.table_name = '{table}' and tc.table_schema = '{schema}';
    """)

    pks = [ row[0] for row in cursor.fetchall() ]

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

    columns = [{
        "column_name": row[0],
        "column_name_uppercase": row[0].upper(),
        "column_name_pascalcase": pascalcase(row[0]),
        "column_name_snakecase": snakecase(row[0]),
        "data_type": row[1],
        "go_data_type": to_go_data_type(row[1]),
        "not_null": row[2],
        "primary_key": row[0] in pks
    } for row in cursor.fetchall()]

    cursor.close()
    connection.close()
    return columns


def main():
    # Load cookiecutter.json to get initial values
    COOKIE_CUTTER_JSON = os.path.join(TMP_PATH, 'cookiecutter.json')

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

    cookiecutter(template=TMP_PATH, no_input=True, overwrite_if_exists=True)

    print("generated. Processing tables...")

    # repeat the template processing for each table but pass the "skip_if_file_exists" attribute
    # in order to generate the individual files

    for table in context["selected_tables"]["values"]:
        columns = get_columns_info(host, port, db, user, password, schema, table)
        cookiecutter(template=f"{TMP_PATH}/per_table_templates", no_input=True, overwrite_if_exists=True, skip_if_file_exists=True,
                     extra_context={
                         "table": table,
                         "table_camelcase": camelcase(table),
                         "table_pascalcase": pascalcase(table),
                         "table_snakecase": snakecase(table),
                         "columns": { "values": columns }
                     })

    print("done.")

if __name__ == "__main__":
    main()
