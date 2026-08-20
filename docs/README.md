# Invoice Extract Service

## What It Does
Accepts invoice/receipt images and PDFs, extracts key fields using OCR + AI, and exports to CSV/Excel.

## API Endpoints
- `POST /extract` — Upload invoice image/PDF
- `GET /invoices` — List all extracted invoices
- `GET /invoices/{id}` — Get specific extraction
- `GET /export` — Download all as CSV
- `GET /health` — Health check

## Roadmap
- [ ] GPT-4o vision for extraction
- [ ] Email monitoring for incoming invoices
- [ ] Multi-format export (XLSX, JSON)

## Related
- Motto: "Burn VPS, Not Tokens."
