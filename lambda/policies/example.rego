package top.next.example

default allow = false

allow = true {
    email != null
    endswith(lower(email), "@example.com")
}

user := input.membership.user.login
email := object.get(input.membership.user, "mail", null)