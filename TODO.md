# TODO

Lifetime analyzer for gno

- [ ] Scope Anaylze
  - [ ] Identify the block structure in the AST
  - [ ] Determine the scope of each variable's declaration
  - [ ] Handle nested scopes
  - [ ] Adding scope information to `VaraibleTracker`

- [ ] Control Flow Analysis
  - [ ] Designing a control flow graph (CFG) structure
  - [ ] Handling if-else statements
  - [ ] Handling loops
  - [ ] Handling switch statements
  - [ ] Handling defer statements
  - [ ] Handling goto statements
  - [ ] Handling panic/recover statements
  - [ ] Implement CFG generation functions

- [ ] Analyzing Data Flow
  - [ ] Trace variable statements
  - [ ] Trace variable reads
  - [ ] Handling parameter passing in function calls
  - [ ] Return value processing
  - [ ] Pointer analysis (handling the address operator `&` and the dereference operator `*`)

- [ ] Cross-Functional Analysis
  - [ ] Create a function call graph (FCG)
  - [ ] Implement inter-procedural analysis
  - [ ] Handling recursive functions

- [ ] Type Checking and Inference
  - [ ] Implementing basic type checking
  - [ ] Handling `struct` and `interface`
  - [ ] Implementing type inference logic
  - [ ] Handling generic types (optional)

- [ ] Handling Packages
  - [ ] Model the behavior of standard library functions
  - [ ] Establishing a strategy for handling third-party packages
  - [ ] Analyzing package imports
