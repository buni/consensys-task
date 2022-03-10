# ConsenSys task - Hristo Karchokov
## Endpoints
The API has two endpoints:

POST endpoint `localhost:8080/api/v1/links/` where you need to send a URL list as plain text, where each url is valid and is new line delimited.
The response is `{"data":{"job_id":"some_uuid"}}` where job_id is used to query for the results of the "job" as its run asynchronous

GET endpoint `localhost:8080/api/v1/links/status/{jobID}` which returns either the job results, or an empty response with status 202 as to indicate that he proccessing of the job hasn't finished yet 
```json
{
   "data":{
      "results":[
         {
            "id":"6b0b045b-8eb4-4926-bb8d-e1936f9c368b",
            "job_id":"dc0eb029-ef6d-4906-b442-08f1a1b32470",
            "page_url":"http://google.com/",
            "internal_links_count":4,
            "external_links_count":15,
            "success":true,
            "error":null,
            "created_at":"2022-03-10T11:36:15.0006864Z",
            "updated_at":"2022-03-10T11:36:15.0006865Z"
         }
      ]
   }
}
```

### Example requests
    curl -X POST -d $'http://google.com/\nhttp://youtube.com/\n' http://localhost:8080/api/v1/links/
    {"data":{"job_id":"dc0eb029-ef6d-4906-b442-08f1a1b32470"}}
    
```json
curl http://localhost:8080/api/v1/links/status/dc0eb029-ef6d-4906-b442-08f1a1b32470
{
   "data":{
      "results":[
         {
            "id":"6b0b045b-8eb4-4926-bb8d-e1936f9c368b",
            "job_id":"dc0eb029-ef6d-4906-b442-08f1a1b32470",
            "page_url":"http://google.com/",
            "internal_links_count":4,
            "external_links_count":15,
            "success":true,
            "error":null,
            "created_at":"2022-03-10T11:36:15.0006864Z",
            "updated_at":"2022-03-10T11:36:15.0006865Z"
         },
         {
            "id":"23b18f7a-dfd7-4d68-adae-265ec85ba1be",
            "job_id":"dc0eb029-ef6d-4906-b442-08f1a1b32470",
            "page_url":"http://youtube.com/",
            "internal_links_count":6,
            "external_links_count":8,
            "success":true,
            "error":null,
            "created_at":"2022-03-10T11:36:15.0006891Z",
            "updated_at":"2022-03-10T11:36:15.0006892Z"
         }
      ]
   }
}
```

### Setup
Run make up/down to start/stop the service (make sure you don't have anything running on port 8080)
