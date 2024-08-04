CREATE TABLE IF NOT EXISTS orders (
    id INT NOT NULL AUTO_INCREMENT,
    customer_id INT NOT NULL,
    no_contract VARCHAR(100) NOT NULL,
    otr DECIMAL(16, 2) NOT NULL,
    admin_fee DECIMAL(16, 2) NOT NULL,
    installment_value DECIMAL(16, 2) NOT NULL,
    interest_value DECIMAL(16, 2) NOT NULL,
    asset_name VARCHAR(100) NOT NULL,
    tenor INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY(id)
)