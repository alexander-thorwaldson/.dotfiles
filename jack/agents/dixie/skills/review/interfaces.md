# review interfaces

Are the public contracts clean?

## What to Look For

### Contract Quality
- **Contract clarity** — are input/output shapes obvious from the signature? Would a consumer know how to use this without reading the implementation?
- **Principle of Least Astonishment** — does the interface behave the way a caller would expect from its name and signature? No surprising side effects, no hidden state changes.
- **Design by Contract** — are preconditions, postconditions, and invariants clear? Can a caller know what's expected of them and what they'll get back?
- **Uniform Access** — is the notation consistent regardless of whether values are computed or stored?

### Stability and Evolution
- **Breaking changes** — does this modify an existing public interface in a way that breaks current consumers?
- **Hyrum's Law** — even undocumented behavior is depended on. If the change alters observable behavior that isn't part of the contract, it may still break consumers.
- **Open/Closed** — can this interface be extended without modifying existing code? Are new behaviors added via extension rather than modification?
- **Semantic versioning awareness** — does the change warrant a major, minor, or patch version bump? Is the magnitude of the contract change appropriate?

### Surface Area
- **Minimal exposure** — is the change exposing more than it needs to? Could any of the new public surface be kept internal?
- **Interface Segregation** — are consumers forced to depend on methods they don't use? Should this be split into focused interfaces?
- **Encapsulation** — is internal state properly hidden? Can consumers reach into implementation details they shouldn't?
- **Law of Demeter** — does the change encourage callers to chain through objects (`a.getB().getC().doThing()`)? Each component should only talk to its immediate collaborators.
- **Leaky Abstractions** — does the interface leak implementation details that consumers shouldn't need to know about? All abstractions leak eventually, but new ones shouldn't leak immediately.

### Semantics
- **Command-Query Separation** — does each method either perform an action or return data, but not both?
- **Idempotency** — for operations that may be retried (APIs, message handlers), does calling them multiple times produce the same result as calling once?
- **Postel's Law** — is the interface liberal in what it accepts and conservative in what it produces? Strict output contracts, tolerant input handling.
- **Error contracts** — are failure modes explicit? Can consumers distinguish between different kinds of errors?
- **Symmetry** — do related interfaces follow consistent patterns?

## Do Not Flag

- Internal interfaces between private components — these can change freely
- Interface style preferences unless they create actual ambiguity
- Error handling implementation details — how errors are caught is the engineer's concern, whether they're surfaced correctly in the contract is yours
- Law of Demeter violations in builder/fluent APIs — method chaining is the intended pattern there
- Hyrum's Law concerns for internal-only interfaces — it matters for public contracts, not private ones
