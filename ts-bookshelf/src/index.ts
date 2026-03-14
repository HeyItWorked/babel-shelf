// index.ts — Entry point. Connects to Postgres, starts the server.
// Equivalent of main() in Go's main.go.
//
// 1. Read DATABASE_URL from environment
// 2. Initialize the pg Pool (imported from db.ts)
// 3. Test connection with pool.query("SELECT 1")
// 4. Start server on port 8081 using Bun.serve

// import { Pool } from "pg"
// import { pool } from "./db"      — to assign the connection
// import app from "./routes"        — the Hono app with all routes wired

// read DATABASE_URL
// create new Pool({ connectionString: databaseUrl })
// assign to db.ts's exported pool
// test connection: await pool.query("SELECT 1")
// start: export default { port: 8081, fetch: app.fetch }
