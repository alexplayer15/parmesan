## Parmesan 
Parmesan is a CLI tool which can be used to generate requests from valid and well-written Open API Specs.

The idea of the tool is to generate requests which are valid and can be sent without any modification straight out of the box.

## Installation
1. Prerequisites:

- Go v1.24 installed (could work on earlier versions but not yet tested. 1.24 is the safest option).
- Parmesan currently only supports OAS v3.0 

2. Install
- `go install github.com/alexplayer15/parmesan@v0.3.0` (will be updated once there are more stable versions)

## Running
As is stand there is one command available for Parmesan which is the `generate-request` command. This will take in an OAS and generate a `.http` file with actionable requests to the API defined in the spec.

To run the `generate-request` command, enter the following:

`parmesan generate-request <oas-file-location>`

## Example 

If I had a basic OAS like so named hello-oas.yml:

```yaml
openapi: 3.0.3
info:
  title: Parmesan Test API
  description: "A very basic Open API Spec to test Parmesan"
  version: 1.0.0
servers:
  - url: http://localhost:8080
paths:
  /hello:
    get:
      tags:
      - Hello World
      summary: Says hello.
      description: Returns a greeting.
      parameters:
      - name: X-Hello-Header
        in: header 
        description: identifies how many people to say Hello to
        required: true
        schema:
          type: string
        example: "1"
      requestBody:
        content: 
          application/json:
            schema:
              $ref: "#/components/schemas/HelloWorldRequest"
      responses:
        '200':
          description: A successful greeting response
          content:
            application/json:
              example:
                message: "Hello, world!"
  /goodbye:
    get:
      tags:
      - Goodbye
      summary: Says goodbye.
      description: Returns a farewell.
      parameters:
      - name: X-Goodbye-Header
        in: header 
        description: identifies how many people to say Goodbye to
        required: true
        schema:
          type: string
        example: "1"
      requestBody:
        content: 
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                  description: "The name of the user"
                  example: "Mia"
      responses:
        '200':
          description: A successful farewell response
          content:
            application/json:
              example:
                message: "Goodbye!"
components:
  schemas:
    Education:
      type: object
      properties:
        university:
          type: string
          description: "University the user studied at"
          example: "University of Manchester"
        degree:
          type: string 
          description: "Degree user studied"
          example: "Chemical Engineering"
        grade:
          type: string 
          description: "Grade user achieved"
          example: "2:1"
    HelloWorldRequest:
      type: object
      properties:
        name:
          type: string
          description: "The name of the user"
          example: "Alex"
        age:
          type: integer
          description: "The age of the user"
          example: 25
        hobbies:
          type: array
          items:
            type: string
          description: "Hobbies of the user"
          example: 
            [ "Boxing", "Video Games", "Football" ]
        favouriteNumbers:
          type: array
          items:
            type: integer
          description: "Users favourite numbers"
          example: 
            [ 15, 15, 15 ]
        favouriteColours:
          type: array
          items:
            type: string
          description: "Users favourite numbers"
          example: 
            ["blue", "purple", "red"] 
        education:
          type: array
          items:
            $ref: "#/components/schemas/Education"
```

If you ran:
`parmesan generate-request hello-oas.yml` then it will create a `.http` in your current directory and generate a file which looks like this:

```http
#### Summary: Says hello.
GET http://localhost:8080/hello
X-Hello-Header: 1
Content-Type: application/json

{
  "name": "Alex",
  "age": 25,
  "hobbies": ["Boxing", "Video Games", "Football"],
  "favouriteNumbers": [15, 15, 15],
  "favouriteColours": ["blue", "purple", "red"],
  "education": [
    {
      "university": "University of Manchester",
      "degree": "Chemical Engineering",
      "grade": "2:1"
    }
  ]
}

#### Summary: Says goodbye.
GET http://localhost:8080/goodbye
X-Goodbye-Header: 1
Content-Type: application/json

{
  "name": "Mia"
}
```

## Importance of Example values 
You get the most value out of Parmesan if you have example values in your Spec. Otherwise, it is most likely Parmesan won't be able to generate a request which can be sent without modification.

## Roadmap
These are features I plan on working on soon:

- Customise output (e.g., file path, naming conventions)

- Handling example generation for anyOf

- A send-request command which automates request generation and execution

- Your ideas? If people find the tool useful. I would love to hear any suggestions on how we can improve it!

## Contributing
Contributions, issues, and feature requests are welcome!

Feel free to open a GitHub Issue or submit a Pull Request.

## License
This project is licensed under the MIT License.



