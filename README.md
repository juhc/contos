# contos
ContOS - это минимальная операционная система созданная для запуска Docker контейнеров. Создана с помощью Buildroot/

## Функции
- Удалены неиспользуемые функции и модули ядра
- Минимальный набор утилит (ssh, bash, sudo и другие) добавлен через BusyBox
- Используется OverylayFS для реализации атомарности

## Сборка
Для сборки iso образа потребуется:
- docker
- make

    ```bash
    git clone https://gitflic.ru/project/stud0000241558-utmn-ru/contos.git
    cd contos
    make
    ```
- Пользователь - contos
- Пароль - contos

