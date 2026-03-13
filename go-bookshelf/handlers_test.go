// handlers_test.go — Integration tests. Sends real HTTP requests through
// the full stack (router → handler → db → Postgres) and checks the responses.
// TDD: write tests first, watch them fail, then build to make them pass.
package main
