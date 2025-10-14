### Calls

The `calls` section passes initialization arguments to field wrappers. Each line is a function call where the function name matches a field name; its arguments become parameters available to your generator implementation.

#### Usage highlights:
- Arguments in `calls { ... }` must be valid Go expressions.
- The generated wrappers forward those arguments to `gens` functions.

#### How to run:
```shell
$ datagen gen ./examples/2_calls -f csv -n 3
```

```shell
$ cat session.csv
created_by,created_at
admin,2023-11-24 18:42:13
admin,2023-09-09 16:09:14
admin,2023-11-13 06:20:23
```
