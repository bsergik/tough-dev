// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

Table accounts {
  id integer [pk, increment]
  user_id integer [not null]
  current_amount integer [not null, note: 'Current money amount.']
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
}

Table reasons {
  id integer [pk]
  title varchar(200) [not null]
  created_at timestamp [default: `now()`]
}

Table account_per_operation_history {
  id integer [pk, increment]
  account_id integer [not null, ref: > accounts.id]
  amount integer [not null, note: 'Amount of money added or subtracted.']
  reason_id integer [not null, ref: > reasons.id]
  created_at timestamp [default: `now()`]
}

Table account_per_day_history {
  id integer [pk, increment]
  account_id integer [not null, ref: > accounts.id]
  amount integer [not null]
  created_at timestamp [default: `now()`]
}

Table schema_version {
  version varchar(20)
  dirty boolean [default: false]
}