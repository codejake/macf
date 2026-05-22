#!/bin/bash

echo "[TEST] Should print usage statement"
go run main.go
echo "---"

echo "[TEST] Mixed-case run:"
go run main.go c6:89:f2:d2:DC:3E
echo "---"

echo "[TEST] Invalid MAC -- too short run:"
go run main.go c6:89:f2:d2:DC:3
echo "---"

echo "[TEST] Invalid MAC -- too long run:"
go run main.go c6:89:f2:d2:DC:3e:4d:ab
echo "---"

echo "[TEST] None format uppercase:"
go run main.go --format none --uppercase c6:89:f2:d2:DC:3E
echo "---"

echo "[TEST] Invalid format:"
go run main.go --format bogus c6:89:f2:d2:DC:3E
echo "---"
