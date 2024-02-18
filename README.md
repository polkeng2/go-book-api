# Go-book-api

This project is a simple Go application that provides a RESTful API for managing a book database. It includes functionality for creating, reading, updating, and deleting book records.

## Features

- **Create Book Records**: This feature allows you to add new books to the database.
- **Read Book Records**: This feature allows you to view the details of a specific book in the database.
- **Update Book Records**: This feature allows you to modify the details of an existing book in the database.
- **Delete Book Records**: This feature allows you to remove a book from the database.

## Dockerization

This project is Dockerized, which means it can be easily set up and run in any environment with Docker. Docker takes care of the dependencies and environment variables, so you don't have to worry about setting up your local environment exactly like the production one.

## Getting Started with Docker

To get started with Docker, follow these steps:

1. Install Docker on your machine. Instructions can be found on the [official Docker website](https://docs.docker.com/get-docker/).
2. Clone this repository to your local machine.
3. Navigate to the project directory in your terminal.
4. Run `docker build -t my-app .` to build the Docker image for the project.
5. Run `docker run -p 8000:8000 my-app` to start the project.

You should now be able to access the project at `localhost:8000`.

## Usage

The project provides the following endpoints:

- `GET /book/{id}`: Retrieves a book record. Replace `{id}` with the ID of the book you want to retrieve.

- `POST /book`: Creates a new book record. The request body should be a JSON object representing the book. The book object should include the following fields: `Titol`, `Autor`, `Prestatge`, `Posicio`, `Habitacio`, `Tipus`, `Editorial`, `Idioma`, `Notes`.

- `PUT /book`: Updates a book record. The request body should be a JSON object representing the book. The book object should include the following fields: `Titol`, `Autor`, `Prestatge`, `Posicio`, `Habitacio`, `Tipus`, `Editorial`, `Idioma`, `Notes`, and `ID`.

- `DELETE /book`: Deletes a book record.The request body should be a JSON object with the ID of the book to be deleted.
