# Go Web Auth Project

Простой веб-проект на Go для регистрации и входа пользователей с использованием PostgreSQL.  

## Описание

Проект реализует базовую систему аутентификации с веб-интерфейсом:

- Регистрация новых пользователей.
- Вход существующих пользователей.
- Хранение данных в PostgreSQL.
- Минимальный веб-интерфейс через HTML-страницы (`login.html` и `register.html`).

Администратор добавляется автоматически при первом запуске (`username: admin`, `password: admin1`).

---

## Стек технологий

- Go (Golang)
- PostgreSQL
- [pgx](https://github.com/jackc/pgx) — драйвер для работы с PostgreSQL
- [godotenv](https://github.com/joho/godotenv) — загрузка переменных окружения из `.env`
- HTML (для фронтенда)

---

Переработанный и более правильный код на Go можно посмотреть здесь - https://github.com/hobgoblin167/FULLSTACK_GO_x_ReactNative_Application
