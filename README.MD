# AI-Powered Mock Interview Tool

This project is a mock interviewing tool built with AI. It consists of a Go server and a React TypeScript client with Vite.

## Server

The server is built with Go 1.22 and is located in the `server` directory. It requires the following environment variables:

- `OPENAI_API_KEY`: Your OpenAI API key
- `DB_URI`: MongoDB connection string

These configurations are optional:
- `PORT`: The port number for the server to run
- `CORS_ALLOWED_ORIGINS`: Allowed origin to call the APIs
- `CORS_ALLOWED_METHODS`: Allowed methods of the APIs call
- `CORS_ALLOWED_HEADERS`: Allowed headers of the APIs call

## Client

The client is built using React TypeScript with Vite and Node.js 20. It is located in the `client` directory. It requires one environment variable:

- `VITE_BACKEND_URL`: URL of the backend server

## Running with Docker

To run the project using Docker and docker-compose, follow these steps:

1. Make sure you have Docker and docker-compose installed on your system.

2. Clone the repository and navigate to the project root directory.

3. Create a `.env` file in the root directory with the following content:

   ```
   OPENAI_API_KEY=your_openai_api_key_here
   MONGO_USERNAME=your_mongo_username
   MONGO_PASSWORD=your_mongo_password
   MONGO_DATABASE=your_mongo_database_name
   ```

4. Set the `VITE_BACKEND_URL` in client Dockerfile accordingly if you change the server's port mapping

5. Run the following command to build and start the containers:

   ```
   docker-compose up --build
   ```

   This will start the MongoDB database, the Go server, and the React client.

6. Access the client application by opening a web browser and navigating to `http://localhost`.

7. To stop the containers, use the following command:

   ```
   docker-compose down
   ```

### Notes

When running the server using the provided `docker-compose.yml`, the `DB_URI` environment variable is infered from Mongo's credential.

If there the server's port is changed or if client and server are deployed in separate domain, please update the `BACKEND_URL` arguments on the client's Dockerfile.

## Development

For development purposes, you can run the server and client separately:

### Server

1. Navigate to the `server` directory.
2. Run `go run main.go` to start the server.

### Client

1. Navigate to the `client` directory.
2. Run `npm install` to install dependencies.
3. Run `npm run dev` to start the development server.

Remember to set the required environment variables when running the applications outside of Docker.

## Project Structure

- `server`: Contains the Go server code
- `client`: Contains the React TypeScript client code
- `docker-compose.yml`: Defines the multi-container Docker application
- `Dockerfile` in each directory: Defines how to build the Docker images for server and client