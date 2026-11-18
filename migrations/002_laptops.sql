-- Módulo de préstamo de portátiles.
-- Cada equipo se controla de manera individual mediante su número de serie.

CREATE TABLE IF NOT EXISTS laptops (
  id INT AUTO_INCREMENT PRIMARY KEY,
  brand VARCHAR(100) NOT NULL,
  model VARCHAR(120) NOT NULL,
  serial_number VARCHAR(120) NOT NULL UNIQUE,
  status ENUM('available', 'loaned', 'maintenance') NOT NULL DEFAULT 'available',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS laptop_loans (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  laptop_id INT NOT NULL,
  loaned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  due_at TIMESTAMP NOT NULL,
  returned_at TIMESTAMP NULL,
  status ENUM('active', 'returned') NOT NULL DEFAULT 'active',
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (laptop_id) REFERENCES laptops(id),
  INDEX idx_laptop_loans_user (user_id),
  INDEX idx_laptop_loans_laptop_status (laptop_id, status)
);
