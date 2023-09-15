# Part I. The problem

There are many different microservices in our company.
Many of them want to interact with the user's balance in one way or another.

**The Architectural Committee** decided to centralize the work with the user's balance in a separate service.

# Part II. Task definition

You need to implement a microservice which handle users wallet balance.

Through this microservice we should be able to:

1. Deposit fund
2. Charge funds for some service in two-step way: **hold** funds first, then **charge** or **unhold** them
3. Show balance of a specific user
4. Show history of transactions for a specific user
5. Give summary report of all users transactions for certain period

# Part III. Requirements

User's balance - very sensitive data. Here we work with real money actually, **so data integrity and consistency is a must.**

For example: situation when user's balance can go below zero is absolutely unacceptable

As an input data we can rely on this fields coming from outside:

1. UserID
2. OrderID
3. ServiceID
4. Amount