### `tags`

1. Filter by single tag (selects orders model) 
```shell
$ datagen gen ./examples/6_metadata/tags -t service=delivery
```
Output:
```shell
$ orders{order_id:89847c16-cb80-40f6-9b44-2170de1bd765 rider_id:387 item:sushi}
```

2. Filter by single tag (selects both models with team=Slytherin)
```shell
$ datagen gen ./examples/6_metadata/tags -t team=Slytherin
```
Output:
```shell
orders{order_id:d19ced78-c0fc-4032-925a-32201232e87a rider_id:356 item:sushi}
users{id:1 name:Laverne Carter email:ashlyhartmann@kris.info}
```

3. Filter by multiple tags (selects only users model)
```shell
$ datagen gen ./examples/6_metadata/tags -t service=user,team=Slytherin
```
Output:
```shell
users{id:1 name:Catherine Kiehn email:lonniekautzer@kautzer.com}
```
4. Filter by multiple tags (selects no models)
```shell
$ datagen gen ./examples/6_metadata/tags -t service=delivery,team=user
```
Output: 
```shell
No output - no models found matching the specified tag
```

