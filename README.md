#  CrackHash

**CrackHash** — распределённая система для подбора слов по заданному хэшу методом brute-force. Система реализована на Golang с использованием Docker Compose и состоит из двух основных компонентов: менеджера и воркера.

---

##  Описание проекта

Проект моделирует работу распределённой системы для подбора паролей по хэшу (в данной лабораторной — MD5). Менеджер принимает запрос, разбивает задачу на подзадачи и рассылает их воркерам, которые производят перебор слов в своём диапазоне. Результаты собираются и возвращаются пользователю.

---

##  Архитектура

- **Менеджер**  
  Принимает запросы от клиентов, делит задачу на части и отправляет воркерам. Хранит состояние задач в оперативной памяти. Обеспечивает REST API.

- **Воркеры**  
  Принимают подзадачи от менеджера, перебирают слова на основе алфавита и длины, вычисляют их MD5-хэши и возвращают найденные совпадения.

 Взаимодействие между сервисами — через HTTP внутри Docker Compose.

 Форматы:
- Менеджер ⇄ Клиент: JSON
- Менеджер ⇄ Воркер: XML (генерация моделей из XSD-схем)

---

##  Запуск через Docker Compose

```bash
docker-compose up --build
```

## Формат отправки запросов клиентом

Отправка запроса на взлом хэша
```bash
сurl -X POST http://localhost:8080/api/hash/crack \                                               
  -H "Content-Type: application/json" \
  -d '{"hash": "81dc9bdb52d04dc20036dbd8313ed055", "maxLength": 4}'
```
Получение результата по requestId
```bash
curl -X GET "http://localhost:8080/api/hash/status?requestId=61a6d633-8392-4367-aba0-50abe39cdcf8"
```
