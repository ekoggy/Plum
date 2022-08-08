# Установка
  Прежде, чем приступать к установке данной утилиты, убедитесь, что на вашем ПК присутствуют:
  - Установите Golang на вашу систему.
``` 
  sudo apt install golang
``` 
  - Произведите установку PostgreSQL.
```
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get update
sudo apt-get -y install postgresql
```
  - Proxy-сервер от TOR, который поднят на 127.0.0.1:9050 ([Ubuntu documentation](https://help.ubuntu.ru/wiki/tor) до пункта "проверка")
```
sudo apt-get install tor tor-geoipdb privoxy 
sudo gedit /etc/privoxy/config
Вставляем следующее и сохраняем:
confdir /etc/privoxy
logdir /var/log/privoxy
# actionsfile standard  # Internal purpose, recommended
actionsfile default.action   # Main actions file
actionsfile user.action      # User customizations
filterfile default.filter
 
logfile logfile
#jarfile jarfile
#debug   0    # show each GET/POST/CONNECT request
debug   4096 # Startup banner and warnings
debug   8192 # Errors - *we highly recommended enabling this*
 
user-manual /usr/share/doc/privoxy/user-manual
listen-address  127.0.0.1:8118
toggle  1
enable-remote-toggle 0
enable-edit-actions 0
enable-remote-http-toggle 0
buffer-limit 4096
```
  - Любой VPN сервис (при тестировании использовался [ProtonVPN](https://protonvpn.com/ru/))
            
При наличии всех компонентов скачайте/импортируйте/клонируйте репозиторий в необходимую вам среду.
    Важным для внимания является файл config.ini, в котором содержатся данные о пользователе, собирающем данные. Отредактируйте файл в соответствии с вашим настроенным пользователем.
    В момент, когда вы имеете правильный config.ini и исходный код программы, запустите на своем ПК proxy-сервер от TOR, VPN и пропишите в косоли `make`, после чего получите .ехе файл, готовый к работе (в процессе создания исполняемого файла будет скачиваться и устанавливаться TelegramAPI, поэтому рекомендую запастись терпением и чаем) 
# Структура
  Структура проекта составляется несколькими пакетами, в которых лежат одноименные файлы с исходным кодом. Кроме того, в структуру входят html-файл страницы и файл стилей, которые лежат в своих собственных папках, в связи с тем, чтобы при возможной модификации и добавлением новых страниц, все однотипные объекты лежали в одной директории.
  Каждый из файлов исходного кода содержит в себе свой уникальный функционал:
  - **postgre.go** отвечает за функции связанные с работой с БД. В данный момент там находится только внесение записи в таблицу, но в будущем файл может быть расширен
  - **telegram.go** содержит в себе простейший клиент для социальной сети телеграмм. В качестве функционала в данном клиенте представлены функции авторизации, считывания нескольких сообщений из канала и мониторинг обновлений
  -  **darknet.go** описывает функции считывания и мониторинга записей на сайте по продаже баз данных. Здесь представлены элементы фреймворка **Colly**, при помощи которого было реализовано считывание данных со страницы
  -  **view.go** включает в себя управление сервером, который поднимается на 127.0.0.1:3000. Если быть точнее, то в этом файле описана функция для реакции на обновление страницы и функция запуска сервера    
  -  **main.go** Соединяет в себе все вышеперечисленные инструменты
# Алгоритмы работы
  Как и было упомянуто, при разработке программы использовались некоторые интрументы, расширяющие функционал языка GO. Это были и драйвера для БД, и фреймворк для скраппинга страниц, и официальная библиотека для телеграмма, портированная для GO.
  1. Принцип работы функций, связанных с PostgreSQL, описывается весьма просто: 
     - Запрос - при помощи драйвера и функций Exec и Prepare отправляется запрос к БД, в котором содержится либо установление соединения с таблицей, либо проверка записи на существование, либо добавление записи в таблицу
     - Выполнение - БД выполняет то, что от неё потребовали
     - Ответ - программа получает ответ и обрабатывает его
2. При работе с Telegram выполнялась следующая последовательность действий:
   - Создание нового клиента - используется функция NewClient, которая и возвращает указатель на новый объект, принимающий данные
   - Авторизация - использование Authorize, которая выполняет передачу параметров пользователя и авторизацию
   - Считывание старых сообщений - GetChatHistory позволила считать N сообщений, начиная с последнего
   - Установка фильтра - написание функции фильтрации приходящих обновлений посредством проверки ID чата, из которого обновление пришло
   - Мониторинг обновлений - создание ресивера при помощи AddEventReceiver, который мониторит события, указанные нами
3. Считывание данных из Darknet происходит по следующему алгоритму:
   - Создание коллектора - при помощи одноименной функции создается коллектор (оператор, который обращается к сайту и хранит его ответы на запросы)
   - Установка proxy-сервера - SetProxyFunc позволяет настроиться на прокси сервер TORа и подключаться к .onion ресурсам
   - Формирование запросов - функции onHTML, onError, ... позволяют сформировать требования от сайта в том или ином случае 
   - Подключение к сайту - фунция visit отправляет запрос к сайту
   - Принятие ответов на запросы - принятие кода страницы сайта, который отфильтрован с учетом запросов
