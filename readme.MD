# ТЗ yadro
Программа написана с использованием стандартной библиотеки. Используется для парсинга txt файла и выводит результат в консоль
## Запуск
Через консоль:
```sh
$ go run main.go --path=PATH_TO_THE_FILE
```
Через собранный проект для windows:
```sh
$ go build
$ golang-yadro-test.exe --path=PATH_TO_THE_FILE
```
## Требования:
Если не указывать параметр path, то файл будет искаться в дерриктории проекта по имени test_file.txt Поэтому параметр path обязателен для собранного проекта под windows или перед сборкой необходимо добавить test_file.txt в каталог.
