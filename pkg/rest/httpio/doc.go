// Package httpio provides utility functions and types for handling HTTP responses.
//
// This package includes types for representing HTTP responses, such as the Response[T] structure,
// which encapsulates the status code, headers, and body of an HTTP response with a generic body type T.
// It also provides functions like WriteJSON, which writes a JSON response to the http.ResponseWriter.
//
// Types:
//   - Response[T]: Represents an HTTP response structure with a generic body T.
//   - Error: Represents an application-specific error structure with an HTTP status code, a key, and a description.
//     It implements the error interface.
//   - Message: Represents a generic message structure used for communication, with a key and content fields.
//
// Functions:
//   - WriteJSON[T]: Writes a JSON response to the provided http.ResponseWriter.
//   - URLParam[T]: Extracts a URL parameter named 'key' from the given HTTP request and converts it to the specified type T.
//
// Example usage:
//
//	// Construct a custom error and respond with a JSON representation of the error.
//	customError := respond.Error{Status: 400, Key: "invalid_input", Desc: "Invalid input provided"}
//	response := respond.Response[Message]{Status: customError.Status, Body: respond.Message{Key: customError.Key}}
//	respond.WriteJSON(w, r, response)
//
// Note: This package assumes the use of the chi router for URL parameter extraction.
package httpio
