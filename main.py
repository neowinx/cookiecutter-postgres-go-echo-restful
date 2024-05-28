#!/usr/bin/env python3

import psycopg2
import inquirer
import json
import os

HOOK_PATH=os.path.dirname(__file__)

def get_postgresql_tables(host, port, db, user, password):
    connection = psycopg2.connect(
        host=host,
        port=port,
        database=db,
        user=user,
        password=password
    )
    cursor = connection.cursor()
    cursor.execute("""
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public'
    """)
    tables = [row[0] for row in cursor.fetchall()]
    cursor.close()
    connection.close()
    return tables

def main():
    # Load cookiecutter.json to get initial values
    COOKIE_CUTTER_JSON = os.path.join(HOOK_PATH, 'cookiecutter.json')

    with open(COOKIE_CUTTER_JSON) as f:
        context = json.load(f)

    questions = [
        inquirer.Text('postgresql_host', default="localhost"),
        inquirer.Text('postgresql_port', default=5432),
        inquirer.Text('postgresql_db', default="super"),
        inquirer.Text('postgresql_user', default="postgres"),
        inquirer.Text('postgresql_password', default="postgres"),
    ]

    answers = inquirer.prompt(questions)

    if not answers:
        print("No answers returned. Exiting.")
        exit(1)

    host = answers["postgresql_host"]
    port = answers["postgresql_port"]
    db = answers["postgresql_db"]
    user = answers["postgresql_user"]
    password = answers["postgresql_password"]

    tables = get_postgresql_tables(host, port, db, user, password)

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
    
    # Update context with selected tables
    context['selected_tables'] = answers['selected_tables'] if answers['selected_tables'] else ["placeholder"]
    
    # Write the updated context back to cookiecutter.json
    with open(COOKIE_CUTTER_JSON, 'w') as f:
        json.dump(context, f, indent=4)

if __name__ == "__main__":
    main()
