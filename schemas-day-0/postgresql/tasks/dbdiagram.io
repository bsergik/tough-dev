// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

Table tasks {
  id integer [pk, increment]
  title varchar(300) [not null]
  description text
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
}

Table user_has_tasks {
  id integer [pk, increment]
  task_id integer [not null, ref: > tasks.id]
  user_id integer [not null]
  created_at timestamp [default: `now()`]
}

Table statuses {
  id integer [pk, increment]
  name varchar(50) [not null]
  description text
  created_at timestamp [default: `now()`]
}

Table task_has_statuses {
  task_id integer [not null, ref: > tasks.id]
  status_id integer [not null, ref: > statuses.id]
  user_id integer [not null, note: 'Who changed status']
  created_at timestamp [default: `now()`]
}

Table schema_version {
  version varchar(20)
  dirty boolean [default: false]
}