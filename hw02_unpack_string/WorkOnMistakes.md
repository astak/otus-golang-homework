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