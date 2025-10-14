### Iter

`iter` is a zero-based row counter available inside every generator function. It represents the record index currently being generated for a model.

#### Key points:
- `iter` starts at 0 and increases by 1 for each record.
- The same `iter` value is used across all fields of a model for a given record.

#### How to run:
```shell
$ datagen gen ./examples/4_iter -n 3
```

Output:
```shell
users{id:1 name:Eino Mayert}
users{id:2 name:Savannah Casper}
users{id:3 name:Keven Steuber}
```

Notice how the `id` field uses `iter + 1` to IDs starting from 1, while `name` generates random names for each record.
