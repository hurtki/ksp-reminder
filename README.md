# ksp-reminder
Rest service that manages your ksp products to be on the shop branch

# Endpoints 

- **`/`** **GET** - Retrieves all your reminders in this structure 
```json
[
  {
    "Article": 298053,
    "BranchesIDs": [
      51,
      63
    ],
    "Updates": [
      {
        "Time": "2025-09-26T22:20:30.355179+03:00",
        "IsFound": false,
        "BranchFoundOn": {
          "id": 0,
          "name": ""
        },
        "Error": "API error: 429 429 Too Many Requests"
      }
    ]
  }
]
```
# Structure overview
![alt text](docs/get_reminder_api_response_structure.png)

- **`/`** **POST** - Endpoint to add a new reminder, your body need to look in this way

```json
{
    "Article": 298043,
    "Branches": [45, 53, 65]
}
```

## So the article you get from the link to the item
 ![alt text](docs/article-from-link.png)
## And Branches you are getting from 
**[File with branches ids](docs/branches.json)**
