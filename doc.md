# Reglas del proxy

El usuario debe proporcionar un fichero de reglas desde la línea de comandos.
Para permitirle hacerlo, usamos el paquete `flag`.

```go
func main() {
    rulesFilename := flag.String("rules", "rules.csv", "path to the file containing the rules to process")
    flag.Parse()
}
```

A continuación, intentamos abrir el fichero de reglas.

> Si abrimos el fichero, usamos `defer file.Close()` para cerrarlo al finalizar la ejecución.

```go
func main() {
    ...
    file, err := os.Open(*rulesFilename)
    if err != nil {
        log.Printf("Unable to open rules file %q", *rulesFilename)
        os.Exit(1)
    }
    defer file.Close()
}
```

Una vez hemos abierto con éxito el fichero, lo leemos.

Como es un fichero CSV, usamos la función `Read()` del paquete [`encoding/csv`](https://pkg.go.dev/encoding/csv).

Usamos `Read()` y no `ReadAll()` porque el fichero de reglas puede ser arbitrariamente grande y no queremos tener problemas de memoria.

Esto nos obliga a leer línea a línea el fichero CSV hasta encontrar el final, marcado por el error `io.EOF`.

Si encontramos algún error durante la lectura del fichero, por ejemplo, porque contiene un número no consistente de campos, registramos el error pero seguimos adelante leyendo el fichero.

Finalizamos el bucle infinito de lectura del fichero sólo si hemos llegado al final del fichero:

```go
func main() {
    ...
    r := csv.NewReader(file)
    for {
        line, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("error processing line %s: %s. Ignoring line.\n", line, err.Error())
            continue
        }
        fmt.Println(line)
    }
    ...
}
```

Creamos un tipo `Rule` para almacenar las reglas.

```go
type Rule struct {
    Protocol string
    Fqdn     string
    Port     string
    Action   string
}
```

Antes de inicial el bucle en `main`, definimos e inicializamos (vacío) un *slice* de `[]*Rule{}`.

Cada una de las líneas válidas del fichero de entrada lo *parsearemos* como un elemento en este *slice* mediante la función `NewRule`:

```go
func main() {
    ...
    r := csv.NewReader(file)
    rules := []*Rule{}
    for {
        line, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("error processing line %s: %s. Ignoring line.\n", line, err.Error())
            continue
        }

        rules = append(rules, NewRule(line))
    }

    for _, rule := range rules {
        fmt.Println(rule.String())
    }
}

func NewRule(fields []string) *Rule {
    rule := Rule{}
    for range fields {
        rule.Protocol = fields[0]
        rule.Fqdn = fields[1]
        rule.Port = fields[2]
        rule.Action = fields[3]
    }
    return &rule
}
```

Cada uno de los campos de una `Rule` deben cumplir unas determinadas condiciones, pero lo dejamos para más adelante.

Una vez que tenemos todas las reglas (*válidas*) en el *slice* de `Rule`, las convertimos en un fichero JSON de salida.

```go
func main() {
    ...
    jsonOutput, err := json.Marshal(rules)
    if err != nil {
        log.Printf("unable to convert to JSON: %s\n", err.Error())
        os.Exit(1)
    }
    fmt.Println(string(jsonOutput))
}
```

La salida tiene los nombres de los campos de cada regla en mayúculas:

```bash
[{"Protocol":"tcp","Fqdn":" www.ubuntu.com","Port":" 443","Action":" allow"},{"Protocol":"tcp","Fqdn":" www.badsite.com","Port":" 80","Action":" deny"}]
```

Lo habitual es que los nombres de los campos en un fichero JSON estén en minúsculas; para conseguirlo, usamos *JSON tags*:

```go
type Rule struct {
    Protocol string `json:"protocol"`
    Fqdn     string `json:"fqdn"`
    Port     string `json:"port"`
    Action   string `json:"action"`
}
```

Ahora la salida es:

```bash
[{"protocol":"tcp","fqdn":" www.ubuntu.com","port":" 443","action":" allow"},{"protocol":"tcp","fqdn":" www.badsite.com","port":" 80","action":" deny"}]
```

El siguiente paso es guardar el resultado en un fichero.

Añadimos un *flag* para que el usuario pueda especificar el nombre del fichero de salida.

```go
func main() {
    jsonFilename := flag.String("out", "rules.json", "path to the output JSON file containing the processed rules")
    ...
    jsonOutput, err := json.Marshal(rules)
    if err != nil {
        log.Printf("unable to convert to JSON: %s\n", err.Error())
        os.Exit(1)
    }
    os.WriteFile(*jsonFilename, jsonOutput, 0644)
}
```

Una vez tenemos la funcionalidad básica, añadiremos las funciones de validación de los diferentes campos de cada `Rule`.

Como el objetivo es únicamente incluir en el *slice* de `Rule` reglas básicas, debemos realizar la validación de cada campo en el fichero de entrada **antes** de crear la `Rule`.

El lugar lógico para insertar las validaciones de los campos de las reglas expresadas en el fichero CSV es **después** de realizar la validación de que cada línea del fichero CSV es válida, pero **antes** de llamar a la función `NewRule()`.

```go
func main() {
    ...
    r := csv.NewReader(file)
    rules := []*Rule{}
    for {
        line, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("error processing line %s: %s. Ignoring line.\n", line, err.Error())
            continue
        }
        // --> Validate rules here <--

        rules = append(rules, NewRule(line))
    }
    ...
}
```

Como no hemos usado TDD, no hemos visto hasta ahora que el fichero CSV contiene espacios (el paquete CSV los considera parte del valor de cada campo)... Por tanto `..., 80, ...` significa que el campo `Port` es `80`, no `80` (contiene un espacio antes del `8`).

Para evitar problemas, una vez hemos pasado las validaciones de sintaxis del fichero CSV, eliminamos todos los espacios para cada línea antes de validar sus valores:

```go
func main() {
    ...
    for i := range line {
        line[i] = strings.TrimSpace(line[i])
    }
    // --> Validate rules here <--
    ...
}
```

Todos los campos de una regla deben ser válidos para poder crear la regla.

En una misma regla puede, del mismo modo, haber uno o más campos inválidos, por lo que la función de validación debe devolver un *slice* de errores.

```go
func validateFields(fields []string) []error {
    var validationErrors = make([]error, 0)
    if !(fields[0] == "tcp" || fields[0] == "udp") {
        validationErrors = append(validationErrors, fmt.Errorf("%w. %q must be one of 'tcp' or 'udp'", ErrInvalidProtocol, fields[0]))
    }
    if !(fields[3] == "allow" || fields[3] == "deny") {
        validationErrors = append(validationErrors, fmt.Errorf("%w. %q must be one of 'allow' or 'deny'", ErrInvalidAction, fields[3]))
    }
    return validationErrors
}
```

Si no se produce ningún error de validación, el *slice* está vacío; si no, se ha producido algún error de validación.

```go
func main() {
    ...
    // Validate fields
    fieldErrors := validateFields(line)
    if len(fieldErrors) > 0 {
        for i := range fieldErrors {
            log.Printf("error processing line %s: %s.", line, fieldErrors[i])
        }
        continue
    }
    ...
}
```
