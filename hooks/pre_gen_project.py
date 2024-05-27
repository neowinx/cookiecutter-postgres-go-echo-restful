import psycopg2
import inquirer
from cookiecutter.main import cookiecutter
import os

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
    project_slug = "{{cookiecutter.project_slug}}"
    host = "{{cookiecutter.postgresql_host}}"
    port = "{{cookiecutter.postgresql_port}}"
    db = "{{cookiecutter.postgresql_db}}"
    user = "{{cookiecutter.postgresql_user}}"
    password = "{{cookiecutter.postgresql_password}}"

    tables = get_postgresql_tables(host, port, db, user, password)
    questions = [
        inquirer.Checkbox(
            'selected_tables',
            message="Select the tables to include in your project",
            choices=tables
        ),
    ]
    answers = inquirer.prompt(questions)
    
    selected_tables = answers['selected_tables']
    
    # Update cookiecutter context with selected tables
    context = {
        "cookiecutter": {
            "selected_tables": selected_tables
        }
    }
    
    # Generate the project with the updated context
    cookiecutter(
        os.path.abspath(os.path.join(os.path.dirname(__file__), '..')),
        extra_context=context,
        no_input=True
    )

if __name__ == "__main__":
    main()
