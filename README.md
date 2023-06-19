### mitum-dao

*mitum-dao* is a [mitum](https://github.com/ProtoconNet/mitum2)-based contract model and is a service that provides dao function.

#### Installation

```sh
$ git clone https://github.com/ProtoconNet/mitum-dao

$ cd mitum-dao

$ go build -o ./md ./main.go
```

#### Run

```sh
$ ./md init --design=<config file> <genesis config file>

$ ./md run --design=<config file>
```

[standalong.yml](standalone.yml) is a sample of `config file`.
[genesis-design.yml](genesis-design.yml) is a sample of `genesis config file`.