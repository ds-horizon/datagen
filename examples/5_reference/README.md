### Reference

Shows how you can use `iter` to refer to any arbitrary row of any arbitrary column of any arbitrary model, via `self.datagen.<model>().<field>(iter)`.

#### How to run:
```shell
$ datagen gen ./examples/5_reference -f csv -n 3
```

Output:

```shell
$ cat orders.csv
user_id,user_name,order_count
1,Jordy Bashirian,1
2,Luisa Adams,4
3,Willow Marks,10
```

```shell
$ cat users.csv
id,name
1,Jordy Bashirian
2,Luisa Adams
3,Willow Marks
```
