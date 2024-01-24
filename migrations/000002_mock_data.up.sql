INSERT INTO users(id, username, password)
VALUES 
    (1, 'alice', '123456'),
    (2, 'bob', '123456');

INSERT INTO accounts(id, name, bank, balance, user_id)
VALUES 
    (1, 'alice 1', 'VCB', 12000, 1),
    (2, 'bob 1', 'VCB', 12000, 2),
    (3, 'alice 2', 'ACB', 12000, 1),
    (4, 'bob 2', 'ACB', 12000, 2),
    (5, 'alice 3', 'VIB', 12000, 1);