The Challenge
=============


Your task is to build an app which reads customer data from a CSV file (customers.csv) and sends out reminders of unsettled invoices based on the specified schedule. 
The schedule's values are offsets relative to the first message being sent. 
To send a message make a POST request to the commservice's `/messages` endpoint. 
It expects the message to be JSON encoded in the request body and contain the customer's email and message's text, e.g. look like

```{.json}
{
    "email": "user@mail.com",
    "text": "hello user"
}
```

If the communication service has been able to decode the message successfully, it responds with a `201` HTTP status code. 
It is possible that a customer settles an invoice. 
In this case, the communication service's JSON encoded response body additionally contains `{"paid": true}` and the customer shouldn't receive any additional messages.
After your service has sent out all messages terminate the commservice to get a report.


# Develop and run

All artifacts are placed inside `/bin/<platform>/` directory. Supported platforms are `darwin` and `linux`.


## How to build

    make build

## Run tests

    make test-unit

## Run intergration tests

    make test-integration

## Run all tests

    make test

## Clear environment

    make clean
