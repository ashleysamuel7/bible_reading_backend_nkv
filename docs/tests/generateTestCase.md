You are a test engineer. Generate comprehensive test cases for every API endpoint of the backend. The backend has resources: Auth, Users, Todos, and Health. For each endpoint produce:

- Endpoint metadata: Method, Path, Auth required (Yes/No)
- Purpose (one line)
- Test cases (list) with for each:
  - id: a short id (e.g., AUTH.LOGIN.001)
  - title
  - type: unit/integration/contract/security
  - preconditions (e.g., user exists, DB has seed data)
  - steps (numbered)
  - request (HTTP method, full URL template, headers, query params, JSON body example)
  - expected status code(s)
  - expected response JSON schema (concise JSON Schema)
  - expected response example (concrete JSON)
  - cleanup steps (if any)
  - severity (critical/high/medium/low)

Return ONLY JSON with top-level key `tests` (array of test-case objects). No commentary.
