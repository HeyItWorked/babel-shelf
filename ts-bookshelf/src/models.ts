// models.ts — Defines the Book type and status constants.
// Shared across all layers. Equivalent of models.go in Go.
//
// TypeScript advantage: union types catch invalid status at compile time.
// Go uses string constants — typos only caught at runtime.

export type BookStatus = "want to read" | "reading" | "finished";

export const STATUS_WANT_TO_READ: BookStatus = "want to read";
export const STATUS_READING: BookStatus = "reading";
export const STATUS_FINISHED: BookStatus = "finished";
export const VALID_STATUSES: BookStatus[] = [STATUS_WANT_TO_READ, STATUS_READING, STATUS_FINISHED];

export interface Book {
  id: number;
  title: string;
  author: string;
  status: BookStatus;
}
