# Example 3

> Typical project style: DDD, Clean architecture, Hexagonale.

More complicated case: domain model `User` differs from database (one domain
model stored in two tables).

File `./adapters/mysql/user_repo.go` is crafted by hands to unite two generated
repo in one.
