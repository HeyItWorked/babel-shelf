// index.ts — Entry point. Connects to Postgres, starts the server.

import { Pool } from "pg"
import { pool, setPool } from "./db"
import app from "./routes"

// create pool needs dbUrl -> only index can
const databaseUrl = process.env.DATABASE_URL
if (!databaseUrl){
    console.error("DATABASE_URL is required")
    process.exit(1)
}

// db.ts has the parking spot, index.ts buys the car and parks it.
// db.ts just drives (pool.query). only index.ts and tests know the connection string.
setPool(new Pool({ connectionString: databaseUrl }))
await pool.query("SELECT 1")
console.log("connected to database")

export default { port: 8081, fetch: app.fetch }
