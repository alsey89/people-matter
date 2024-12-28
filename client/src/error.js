class APIError extends Error {
  constructor(message, status) {
    super(message); // Pass the message to the Error constructor
    this.status = status; // Add a status property
    this.name = this.constructor.name; // Set the error name to the class name
  }
}

export { APIError };
