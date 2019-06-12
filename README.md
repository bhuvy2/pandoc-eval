# pandoc-eval

Pandoc Filter for evaluating code snippets.

    # test.md
    
    ```evallua
    print("Hello World")
    ```

```bash
$ pandoc test.md --filter pandoc-eval
# test.md

Hello World
```

## Supported Environments

|   name   | language  |
|----------|-----------|
| evallua  |    lua    |

## Upcoming Environments

* Python
* Ruby
* Javascript
