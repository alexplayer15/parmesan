Parmesan 
Parmesan is a lightweight CLI tool for automatically generating valid example HTTP requests from an OpenAPI Specification (OAS).
It helps you quickly understand and interact with unfamiliar APIs by producing .http files based on the API spec — ready to use and test.

Features
Generate .http example requests from OpenAPI specs.

Supports application/json request bodies.

Resolves $ref, allOf, oneOf, and example/default values intelligently.

Fast and easy to use — just one command.

Minimal setup, no boilerplate.

Installation
Prerequisite:

Go installed

Then install Parmesan using:
go install github.com/alexplayer15/parmesan@v0.3.0

Quick Start
To generate HTTP request examples from an OpenAPI file run:

parmesan generate-request <oas-file-location>


How It Works
Parses your OpenAPI spec (YAML or JSON).

Walks through each path and HTTP method.

Generates requests based on your OAS.

Populates headers and bodies using examples, defaults, or intelligent fallbacks.

Handles complex schemas with $ref, allOf, and oneOf constructs.

📝 Example Output
#### Summary: Create a new user
POST https://api.example.com/users
Content-Type: application/json

{
  "name": "example value",
  "email": "example@example.com",
  "age": 0
}
 Roadmap
 Initial .http request generation

 Handle schema references and compositions

 Add support for more content-types (e.g., multipart/form-data)

 Customise output (e.g., file path, naming conventions)

 Handling example generation for anyOf

Contributing
Contributions, issues, and feature requests are welcome!

Feel free to open a GitHub Issue or submit a Pull Request.

📄 License
This project is licensed under the MIT License.



