BIN=$(dirname "$([[ $0 == /* ]] && echo "$0" || echo "$PWD/${0#./}")")

# путь к директории хранилища ключей
export _DIR_STORE_="${BIN}/keystore1"

# алиас закрытого ключа/сертификата
export _ALIAS_CERT_="clearing2021"

# алиас сертификата клирингового сервиса
export _ALIAS_2_="way4_2021"

# путь к файлу хранилища ключей
#export _STORE_="${_DIR_STORE_}/keystore.kdb"
export _STORE_="${_DIR_STORE_}/mq-ams.kdb"

# путь к файлу конфига
export _CONF_="${_DIR_STORE_}/keystore.conf"

# пароль к файлу хранилища ключей
export _PASS_="passw0rd"

# Common Name сертификата
export _DNAME_="CN=mskmqis02v.inet.vtb.ru"

# путь к файлу запроса на выпуск сертификата afsc
export _REQ_CERT_="${_DIR_STORE_}/req-cert-way4.vtb.pem"

# путь к файлу запроса на выпуск сертификата afsc
export _CERT_="${_DIR_STORE_}/cert.vtb.pem"

if [ ! -d "$_DIR_STORE_" ]; then
  mkdir "$_DIR_STORE_"
  [ $? -ne 0 ] && echo "Ошибка при создании директории '$_DIR_STORE_'" && exit 1
fi

export PATH="${PATH}:/opt/mqm/bin"
