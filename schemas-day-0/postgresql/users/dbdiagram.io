// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

Table users {
  id integer [pk, increment]
  name varchar(300) [not null]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
}

Table roles {
  id integer [pk, increment]
  name varchar(300) [not null]
  description text
  may_have_tasks bool [note: 'Tasks can be assigned to this role if the value is true']
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
}

Table user_has_roles {
  user_id integer [not null, ref: > users.id]
  role_id integer [not null, ref: > roles.id]
  created_at timestamp [default: `now()`]

  Indexes {
    (user_id, role_id) [pk]
  }
}

Table schema_version {
  version varchar(20)
  dirty boolean [default: false]
}
