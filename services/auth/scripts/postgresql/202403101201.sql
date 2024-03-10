INSERT INTO roles (id, title) VALUES 
  (1, 'admin'),
  (2, 'user')
  ON CONFLICT DO NOTHING;