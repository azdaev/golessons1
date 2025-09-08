| Status | HTTP Method | Endpoint | Method Name | Description |
|--------|-------------|----------|-------------|-------------|
| ✅ | POST | `/users/` | CreateUser | Создание нового пользователя |
| ✅ | GET | `/users/:id/` | GetUser | Получение информации о пользователе по ID |
| ⏳ | GET | `/users/:id/posts/` | GetPostsByUserID | Получение списка всех постов пользователя
| ⏳ | POST | `/posts/` | CreatePost | Создание нового поста |
| ⏳ | GET | `/posts/` | GetPosts | Получение списка всех постов |
| ⏳ | GET | `/posts/:id/` | GetPost | Получение информации о посте по ID |
| ⏳ | POST | `/posts/:id/comments?reply_to=231e119` | CreateComment | Создание комментария к посту |
| ⏳ | GET | `/posts/:id/comments/` | GetCommentsByPostID | Получение всех комментариев к посту |
