# Работа над ошибками
## Проблемы с производительностью
У кода очень большие проблемы с производительностью. Получение следующего символа приводит к аллокации памяти и копированию данных.
## Варианты решений
### Использовать цикл
Цикл по диапазону можно применять к строке. В этом случае GO неявно выполняет декодирование UTF-8. Такой подход потребляет меньше памяти, делает меньше аллокаций, и в целом отрабатывает быстрее.
### Использовать `DecodeRune`
В отличие от `DecodeRuneInString`, который работет с неизменяемыми строками, `DecodeRune` работает со срезами. А срез это легковесная структура данных, которая ссылается на элементы массива, а не создает новые, как это происходит при выделении подстроки.
## Результаты измерений
Добавил бенчмарк, измеряющий производительность метода `Unpack`:
```
func BenchmarkUnpack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unpack("a4bc2d5e")
	}
}
```
### Исходная неоптимизированная реализация
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkUnpack$ github.com/astak/otus-golang-homework/hw02_unpack_string

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkUnpack-12    	 2448382	       488.4 ns/op	     168 B/op	      13 allocs/op
PASS
ok  	github.com/astak/otus-golang-homework/hw02_unpack_string	1.700s
```

### Реализация через цикл
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkUnpack$ github.com/astak/otus-golang-homework/hw02_unpack_string

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
=== RUN   BenchmarkUnpack
BenchmarkUnpack
BenchmarkUnpack-12       2411959               485.8 ns/op           168 B/op         13 allocs/op
PASS
ok      github.com/astak/otus-golang-homework/hw02_unpack_string        1.686s
```
### Исходная оптимизированная реализация
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkUnpack$ github.com/astak/otus-golang-homework/hw02_unpack_string

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
=== RUN   BenchmarkUnpack
BenchmarkUnpack
BenchmarkUnpack-12       2460854               486.8 ns/op           168 B/op         13 allocs/op
PASS
ok      github.com/astak/otus-golang-homework/hw02_unpack_string        1.702s
```
### Добавил бенчмарк на оба варианта tokenizer
#### Итератор
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkScan$ github.com/astak/otus-golang-homework/hw02_unpack_string/iterator

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string/iterator
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkScan-12    	 4747164	       248.9 ns/op	     128 B/op	       8 allocs/op
PASS
ok  	github.com/astak/otus-golang-homework/hw02_unpack_string/iterator	1.446s
```
#### Цикл
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkScan$ github.com/astak/otus-golang-homework/hw02_unpack_string/loop

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string/loop
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkScan-12    	 5691138	       210.3 ns/op	     128 B/op	       8 allocs/op
PASS
ok  	github.com/astak/otus-golang-homework/hw02_unpack_string/loop	1.420s
```
### Избавился от аллокаций токена паредавая его по значению а не по ссылке
#### Бенчмарк tokenizer на итераторе
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkScan$ github.com/astak/otus-golang-homework/hw02_unpack_string/iterator

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string/iterator
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkScan-12    	18537831	        66.64 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/astak/otus-golang-homework/hw02_unpack_string/iterator	1.308s
```
#### Бенчмарк tokenizer на цикле
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkScan$ github.com/astak/otus-golang-homework/hw02_unpack_string/loop

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string/loop
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkScan-12    	46649948	        26.62 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/astak/otus-golang-homework/hw02_unpack_string/loop	1.276s
```
#### Бенчмарк метода Unpack на базе итератора
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkUnpack$ github.com/astak/otus-golang-homework/hw02_unpack_string

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkUnpack-12    	 3565096	       328.6 ns/op	      40 B/op	       5 allocs/op
PASS
ok  	github.com/astak/otus-golang-homework/hw02_unpack_string	1.519s
```
#### Бенчмарк метода Unpack на базе цикла
```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkUnpack$ github.com/astak/otus-golang-homework/hw02_unpack_string

goos: linux
goarch: amd64
pkg: github.com/astak/otus-golang-homework/hw02_unpack_string
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkUnpack-12    	 4314504	       272.9 ns/op	      40 B/op	       5 allocs/op
PASS
ok  	github.com/astak/otus-golang-homework/hw02_unpack_string	1.468s
```