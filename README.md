## Parmesan 
Parmesan is a CLI tool which can be used to generate requests from valid and well-written Open API Specs.

The idea of the tool is to generate requests which are valid and can be sent without any modification straight out of the box.

## Installation
1. Prerequisites:

- Go v1.24 installed (could work on earlier versions but not yet tested. 1.24 is the safest option).
- Parmesan currently only supports OAS v3.0 

2. Install
- `go install github.com/alexplayer15/parmesan@v0.7.0` (will be updated once there are more stable versions)

## Running
As is, there are two commands available for Parmesan: `generate-request` and `send-request`. `generate-request` will take generate a `.http` file containing all requests defined in the provided OAS. `send-request` will work off the `generate-request` logic to send HTTP requests to the endpoints defined in the OAS.

## Generate Request Command

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

## Flags 

Currently, there are two flags for `generate-request`; `output` and `with-server`. 

`output` allows you to control the directory you want your `.http` files to be outputted to. 

For instance, if I ran `parmesan generate-request hello-oas.yml --output httpRequests`. Then a `httpRequests'
directory would be created unless it already exists and either way the new `.http` file would be generated there.

`with-server` allows you to control which server url you want to generate requests to. Parmesan will look for server urls defined in the OAS and make the choice based off the index value you choose. 0 will choose the first server url and is the default. 

## Send Request Command

To run the `send-request` command, enter the following:

`parmesan send-request <oas-file-location>`

Without flags this will send a request to all endpoints defined in the provided OAS. The requests will only be successful if the example values in the OAS can generate a valid request (or if changes are made through the hooks flag).

All responses will be saved in the current directory in JSON format unless this behaviour is changed with a flag.

## Flags

`output` allows you to control the directory you want your JSON response files to be outputted to.

`with-server` allows you to control which server url you want to send requests to. Parmesan will look for server urls defined in the OAS and make the choice based off the index value you choose. 0 will choose the first server url and is the default. 

`path` allows you to filter which requests you send by path. For instance, if your Spec defines two requests with paths: `health/live` & `health/status` and you enter `health/live` using this path you will only send a request to this path.

`method` allows you to filter which requests you send by method. Allowed methods are GET, POST, UPDATE, PUT, PATH, DELETE (not case-sensitive)

`hooks` allows you to modify specific values in your request body's using a self-made YAML file. A basic hooks file would look like this:

```yaml
- path: /hello
  method: GET
  body:
    name: Theo
    education.degree: Physics
```

You can specify multiple hooks in the same file if you want to modify values in different requests.

You must specify the path and method of the request body you want to modify or Parmesan will not recognise you are trying to modify that request. 

For now, you can only modify string, int, and boolean values including nested values of these types. You can not alter arrays or objects directly. 

This flag currently relies on Go's marshalling rules so if you want to modify a string value which will be interpreted as an int you must use "".

## Roadmap
These are features I plan on working on soon:

- A chain-request command which automates request generation for chains of related requests

- Your ideas? If people find the tool useful. I would love to hear any suggestions on how we can improve it!

## Contributing
Contributions, issues, and feature requests are welcome!

Feel free to open a GitHub Issue or submit a Pull Request.

## License
This project is licensed under the MIT License.



