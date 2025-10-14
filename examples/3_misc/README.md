### Misc

The `misc` section injects arbitrary Go code into the generated package. You can use it to define helpers, constants, and types that your generators can reuse.

#### Usage highlights:
- Define domain helpers and constants in `misc { ... }`.
- Call those helpers inside `gens { ... }` when producing field values.

#### How to run:
```shell
$ datagen gen ./examples/3_misc -f csv -n 3
```

Output:
```shell
$ cat order.csv
order_id,total_amount,shipping_address
4a15d66f-b1fa-4902-9b57-b133b192285a,291.1,"{""line1"":""Ut rerum temporibus."",""city"":""Et.""}"
30f68a9b-b183-4b86-a592-b30b9fa538c4,353.34,"{""line1"":""Laudantium ducimus harum."",""city"":""Vel.""}"
908583bb-15d2-44b3-aa56-e5476a7a8f81,105.47,"{""line1"":""Enim dolores accusamus."",""city"":""Sed.""}"
```
