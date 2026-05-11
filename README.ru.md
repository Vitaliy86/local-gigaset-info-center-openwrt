# Альтернативный сервис погоды для `info.gigaset.net`

Этот проект позволяет развернуть собственный мини-сервис для предоставления прогнозов погоды на ваш IP-телефон Gigaset (например, Gigaset C430A Go), поскольку Siemens прекратил работу `info.gigaset.net`, который ранее предоставлял источники данных.

В настоящее время эта замена крайне ограничена и предоставляет данные о прогнозе погоды только для одного настраиваемого местоположения. Настройка выполняется на стороне сервера, а не через телефон (как раньше).

Этот проект разработан для работы на старом PHP5, потому что именно на нём работает мой старый проверенный QNAP TS-109 Pro II.

Конечно, было бы легко переписать это на `node`, `python` или любом другом языке по вашему выбору.

Делитесь своими впечатлениями и присылайте PR с исправлениями и улучшениями!

## Краткое содержание

Этот проект:
- перенаправит `info.gigaset.net` на ваш собственный HTTP-сервер в локальной сети
- будет предоставлять данные о погоде через два URL, которые используются телефоном для получения данных «Info Center»:
  - `http://info.gigaset.net/info/menu.jsp?lang=2&tz=120&mac=7C2F80XXXXXX&cc=49&handsetid=XXXXXXXXXX&provid=11`\
    Этот URL вызывается при доступе к «Info Center» из меню «Extras».
  - `http://info.gigaset.net/info/request.do?lang=2&tz=120&tick=true&mac=7C2F80XXXXXX&cc=49&provid=11&handsetid=XXXXXXXXXX`\
    Этот URL вызывается при получении данных для бегущей строки на экране ожидания и когда «Info Center» используется как заставка. Параметр `tick=true` иногда присутствует, иногда нет. Я не углублялся.

## Первый шаг: Регистрация бесплатного аккаунта на OpenWeatherMap

Перейдите на страницу [регистрации](https://home.openweathermap.org/users/sign_up) OpenWeatherMap и создайте аккаунт. Затем перейдите на страницу [API-ключей](https://home.openweathermap.org/api_keys) и скопируйте ваш ключ для следующего шага.

## Требование: Включение необходимых расширений PHP

Перед настройкой сервиса убедитесь, что на вашем сервере включены следующие расширения PHP:

- **cURL** (`curl`) — необходимо для получения данных о погоде из API OpenWeatherMap
- **GD** (`gd`) — необходимо для конвертации PNG-иконок погоды в формат Gigaset fnt (битовая карта)

На Synology DiskStation вы можете включить эти расширения в **Web Station** в настройках PHP для вашего виртуального хоста.

## Второй шаг: Настройка сервиса

Скопируйте содержимое этого репозитория в директорию, обслуживаемую как `http://<ваш-сервер>/info`. В моём случае это `/Qweb/info/` на NAS.

Чтобы ваш сервер `apache` корректно обрабатывал предоставленные скрипты, добавьте что-то подобное в файл конфигурации `apache`, например `/usr/local/apache/conf/apache.conf`:

```apache
<Directory "/share/Qweb/info">
    DirectoryIndex menu.jsp
    Order deny,allow
    Deny from all
    Allow from localhost
    Allow from 192.168.10.0/24
	AddType application/x-httpd-php .jsp .do
	SetEnv OPENWEATHERMAP_API_KEY <ключ из шага 1>
	SetEnv CITY "Berlin"
	SetEnv LATITUDE 52.52437
	SetEnv LONGITUDE 13.41053
</Directory>
```

Очевидно, замените `192.168.10.0` на вашу локальную IP-сеть и обновите строки `SetEnv` в соответствии с вашими данными.

Проверьте и активируйте конфигурацию:
```term
# /usr/local/apache/bin/apachectl configtest
# /usr/local/apache/bin/apachectl restart
```

Теперь откройте `http://<ваш-сервер>/info`. Вы должны увидеть что-то подобное:

<div style="overflow:scroll; height:12rem">
  <p style='text-align:center'>Do, 23.05.2024<br/>19,3/23,6°C/0 mm<br/>Bedeckt</p><p style='text-align:center'>Fr, 24.05.2024<br/>17,5/24,8°C/0 mm<br/>Bedeckt</p><p style='text-align:center'>Sa, 25.05.2024<br/>14,2/24,2°C/5 mm<br/>Leichter Regen/Bed.</p><p style='text-align:center'>So, 26.05.2024<br/>13,7/24,8°C/1 mm<br/>Leichter Regen/Mäßig bew.</p><p style='text-align:center'>Mo, 27.05.2024<br/>14,3/25,7°C/0 mm<br/>Bed./Überw. bew.</p><p style='text-align:center'>Di, 28.05.2024<br/>13,2/16,6°C/4 mm<br/>Leichter Regen</p>
</div>

Также проверьте `http://<ваш-сервер>/info/menu.jsp` и `http://<ваш-сервер>/info/request.do`, которые должны показывать то же самое.

## Альтернатива: Настройка Synology DiskStation

Если вы используете Synology DiskStation, вы можете использовать `.htaccess` вместо редактирования глобальной конфигурации Apache.

1. Скопируйте файлы репозитория в директорию с веб-доступом, например `/volume1/web/info/`.
2. В **Web Station** убедитесь, что PHP включён для виртуального хоста, обслуживающего эту директорию. Также убедитесь, что включены необходимые расширения PHP:
   - **cURL** — для получения данных OpenWeatherMap
   - **GD** — для обработки иконок
3. Скопируйте `.htaccess.example` в `.htaccess` и заполните ваши значения:

```apache
DirectoryIndex menu.php

Require local
Require ip 10.0.0.0/8

SetEnv OPENWEATHERMAP_API_KEY <ваш-API-ключ>
SetEnv CITY "Berlin"
SetEnv LATITUDE 52.52437
SetEnv LONGITUDE 13.41053

Options +FollowSymLinks
RewriteEngine On
RewriteRule ^menu\.jsp$ menu.php [L]
RewriteRule ^request\.do$ menu.php [L]
```

   Замените `10.0.0.0/8` на диапазон вашей локальной сети и обновите строки `SetEnv` с вашим городом и координатами.

   Директивы `Require local` и `Require ip` ограничивают доступ только вашей локальной сетью.

### Дополнительно: Отключение иконок погоды

По умолчанию сервис отображает иконки погоды рядом с прогнозом. Если иконки не загружаются или вы предпочитаете текстовый вид, добавьте следующую переменную окружения в ваш `.htaccess`:

```apache
SetEnv SHOW_ICONS false
```

Когда иконки отключены, показывается одно немецкое слово о погоде (например, «Regen», «Sonnig», «Nebel»). Это даёт быструю визуальную сводку ожидаемой погоды без необходимости поддержки изображений.

4. `.htaccess` указан в `.gitignore`, поэтому ваш API-ключ не будет случайно отправлен в систему версионирования.

## Третий шаг: Перенаправление телефона Gigaset на ваш собственный сервер

Это сильно зависит от того, какой у вас роутер. Легче всего для роутеров, которые предоставляют собственный (кэширующий) DNS-сервис, например OpenWRT. Также, если вы можете вручную настроить сопоставление имён хостов с IP-адресами, у вас может всё получиться. Ключевая часть — заставить DNS-сервер, настроенный в базовой станции Gigaset (вероятно, через DHCP), разрешать:

    info.gigaset.net -> <IP вашего сервера>

В OpenWRT вы можете настроить это в `https://<ваш-роутер>/cgi-bin/luci/admin/network/dhcp` на вкладке _Hostnames_.

## Сборка для OpenWRT

### Сборка в OpenWRT SDK

Сначала клонируйте OpenWRT SDK:

```bash
# Скачайте OpenWRT SDK для вашей целевой архитектуры
# Пример для x86_64:
wget https://downloads.openwrt.org/releases/23.05.3/targets/x86/64/openwrt-sdk-23.05.3-x86-64_gcc-12.3.0_musl_x86_64.tar.xz
tar xf openwrt-sdk-23.05.3-x86-64_gcc-12.3.0_musl_x86_64.tar.xz
cd openwrt-sdk-23.05.3-x86-64_gcc-12.3.0_musl_x86_64

# Клонируйте этот репозиторий в директорию packages
git clone https://github.com/Vitaliy86/local-gigaset-info-center-openwrt.git package/gigaset-info-center
```

Затем соберите пакет:

```bash
# Make menuconfig и добавьте gigaset-info-center в Network -> Other packages

# Соберите пакет
make package/gigaset-info-center/compile V=s

# Пакет .ipk будет в:
# bin/packages/<architecture>/base/
```

### Установка на устройство OpenWRT

```bash
# Загрузите пакет .ipk на ваше устройство OpenWRT
scp bin/packages/*/base/gigaset-info-center*.ipk root@ВАШ_IP_OPENWRT:/tmp/

# Установите пакет
opkg install /tmp/gigaset-info-center*.ipk

# Включите и запустите сервис
rc-update add gigaset-info-center default
/etc/init.d/gigaset-info-center start
```

### Настройка сервиса

После установки отредактируйте файл конфигурации окружения:

```bash
vi /etc/gigaset-env.example
# или скопируйте в реальную конфигурацию:
cp /etc/gigaset-env.example /etc/gigaset-env
```

Отредактируйте следующие значения:

| Переменная | Описание | Пример |
|---|---|---|
| `OPENWEATHERMAP_API_KEY` | Ваш API-ключ OpenWeatherMap | `a1b2c3d4e5f6...` |
| `CITY` | Название города для отображения | `"Berlin"` |
| `LATITUDE` | Широта | `52.52437` |
| `LONGITUDE` | Долгота | `13.41053` |
| `SHOW_ICONS` | Показывать иконки погоды (true/false) | `"true"` |

### Настройка lighttpd

Добавьте конфигурацию gigaset-info-center в вашу конфигурацию lighttpd:

```bash
# Включите конфигурацию пакета в /etc/lighttpd/lighttpd.conf
echo 'include "/etc/gigaset-info-center.conf"' >> /etc/lighttpd/lighttpd.conf

# Перезапустите lighttpd
/etc/init.d/lighttpd restart
```

### Настройка DNS-перенаправления для телефона Gigaset

В OpenWRT LuCI настройте DNS-перенаправление:

1. Перейдите на вкладку **Network > DHCP > DNS**
2. Добавьте алиас имени хоста: `info.gigaset.net` -> IP вашего роутера

Или через CLI:

```bash
# Добавьте статическую DNS-запись в /etc/config/dhcp
uci add dhcp static_host
uci set dhcp.static_host.ip="IP_ВАШЕГО РОУТЕРА"
uci set dhcp.static_host.name="info.gigaset.net"
uci commit dhcp
/etc/init.d/network restart
```

### Проверка установки

Проверьте, работает ли сервис:

```bash
# Проверьте статус сервиса
/etc/init.d/gigaset-info-center status

# Протестируйте локально (требуется curl)
curl http://127.0.0.1:8081/
```

Ожидаемый вывод должен показывать данные о погоде в формате XHTML-GP для вашего телефона Gigaset.

## Система сборки

### Структура файлов

```
├── Makefile                    # Определение сборки пакета OpenWrt
├── gigaset-info-center.init    # Init-скрипт OpenWrt
├── etc/
│   ├── lighttpd/               # Конфигурация lighttpd
│   │   └── gigaset-info-center.conf
│   └── gigaset-env.example     # Шаблон конфигурации окружения
├── icons/                      # Иконки погоды
├── .github/workflows/          # GitHub Actions CI/CD
│   └── build-apk.yml           # Определение workflow сборки
```

### Makefile пакета OpenWrt

[`Makefile`](Makefile) следует соглашениям системы сборки пакетов OpenWrt:

```bash
# Сборка пакета (запускается из корня OpenWrt SDK)
make package/gigaset-info-center/compile V=s

# Очистка артефактов сборки
make package/gigaset-info-center/clean

# Настройка в menuconfig
make menuconfig
# Перейдите: Network -> Other packages -> gigaset-info-center
```

### Переменные пакета

| Переменная | Описание |
|---|---|
| `PKG_NAME` | Имя пакета: `gigaset-info-center` |
| `PKG_VERSION` | Версия пакета: `1.7` |
| `PKG_RELEASE` | Номер релиза пакета |
| `DEPENDS` | Зависимости времени выполнения |

### GitHub Actions Workflow

[`.github/workflows/build-apk.yml`](.github/workflows/build-apk.yml:1) workflow:
- Собирает пакет `.apk` при push в main/master или создании тега
- Загружает артефакт для ручного скачивания
- Создаёт GitHub Release при тегировании (`v*`)
- Проверяет целостность структуры пакета

## Ссылки

Я использовал информацию из

- https://www.ip-phone-forum.de/threads/gigaset-infodienst-selbst-gemacht.174719/ (спасибо [VoIPMaster](https://www.ip-phone-forum.de/members/voipmaster.95683/))
- https://copyandpastecode.blogspot.com/2008/08/siemens-s685ip-s68h.html (спасибо [Jon Bright](https://www.blogger.com/profile/13465823659620242219))
- и особенно: http://www.ensued.net/request.do (написано на ruby)

Посмотрите последнее для дополнительных идей о том, как включить RSS-обновления и информацию о общественном транспорте.
