# Requirements and context

## Data models
### User
- has accounts
### Account
- has at least the following fields:
- `name`: name of account
- `bank`: name of bank (3 possible values: `VCB`, `ACB`, `VIB`)
- has transactions
### Transaction
- has at least the following fields:
- `amount`: amount of money
- `transaction_type`: type of transaction (2 possible values: `withdraw`, `deposit`)
## Detail of endpoints
### 1. Get transactions of an user
- URL path: `/api/users/<user_id>/transactions`
- HTTP method: `GET`
- Request:
- Parameters:
  
| Name         | Required | Data type | Description  |
| ------------ | -------- | --------- | ------------ |
| `user_id`    | Yes      | Integer   | User's ID    |
| `account_id` | No       | Integer   | Account's ID |

- Note: When `account_id` is not specified, return all transactions of the user.
- Please have validations for required fields
- Response:
- Content type: `application/json`
- HTTP status: `200 OK`
- Body: Array of user's transactions, each of which has the following fields:
  
| Name               | Data type | Description                 |
| :----------------- | :-------- | :-------------------------- |
| `id`               | Integer   | Transaction's ID            |
| `account_id`       | Integer   | Account's id                |
| `amount`           | Decimal   | Amount of money             |
| `bank`             | String    | Bank's name                 |
| `transaction_type` | String    | Type of transaction         |
| `created_at`       | String    | Created date of transaction |

- Example: GET `/api/users/1/transactions?account_id=1`
- Response:
```json
[{
"id": 1,
"account_id": 1,
"amount": 100000.00,
"bank": "VCB",
"transaction_type": "deposit",
"created_at": "2020-02-10 20:00:00 +0700"
}, { ... }]
```

### 2. Create a transaction for an user
- URL path: `/api/users/<user_id>/transactions`
- HTTP method: `POST`
- Request:
- Parameters:
  
| Name      | Required | Data type | Description |
| --------- | -------- | --------- | ----------- |
| `user_id` | Yes      | Integer   | User's ID   |

- Body:
  
| Name               | Required | Data type | Description         |
| ------------------ | -------- | --------- | ------------------- |
| `account_id`       | Yes      | Integer   | Account's ID        |
| `amount`           | Yes      | Decimal   | Amount of money     |
| `transaction_type` | Yes      | String    | Type of transaction |

- Please have validations for required fields
- Response:
- Content type: `application/json`
- HTTP status: `201 Created`
- Body: Details of the created transaction with the following fields:
  
| Name               | Data type | Description                 |
| ------------------ | --------- | --------------------------- |
| `id`               | Integer   | Transaction's ID            |
| `account_id`       | Integer   | Account's id                |
| `amount`           | Decimal   | Amount of transaction       |
| `bank`             | String    | Bank's name                 |
| `transaction_type` | String    | Type of transaction         |
| `created_at`       | String    | Created date of transaction |

- Example: POST `/api/users/1/transactions`
- Request body:
```json
{
"account_id": 2,
"amount": 100000.00,
"transaction_type": "deposit"
}
```
- Response
```json
{
"id": 10,
"account_id": 2,
"amount": 100000.00,
"bank": "VCB",
"transaction_type": "deposit",
"created_at": "2020-02-10 20:10:00 +0700"
}
```
### 3. Implement PUT and DELETE as your preference

## Requirements

There are three APIs below. (These are actual API, so you can access them.)
Please create a program with some of these APIs. It should take a user ID as input and return the name, account list and balances of that user. In addition, the program should be designed and implemented with maintainability in mind.

You can use Golang to do this programming.

```sh
- Web API
  - https://mfx-recruit-dev.herokuapp.com/users/1
    - レスポンス例 / Sample response
    {
        "id": 1,
        "name": "Alice",
        "account_ids": [
          1,
          3,
          5
        ]
    }
- https://mfx-recruit-dev.herokuapp.com/users/1/accounts
    - レスポンス例 / Sample response
    [
        {
            "id": 1,
            "user_id": 1,
            "name": "A銀行",
            "balance": 20000
        },
        {
            "id": 3,
            "user_id": 1,
            "name": "C信用金庫",
            "balance": 120000
        },
        {
            "id": 5,
            "user_id": 1,
            "name": "E銀行",
            "balance": 5000
        }
    ]
  - https://mfx-recruit-dev.herokuapp.com/accounts/2
    - レスポンス例 / Sample response
    {
        "id": 2,
        "user_id": 2,
        "name": "Bカード",
        "balance": 200
    }
```