# review structure

Are the pieces connected correctly and separably?

## What to Look For

### Boundaries
- Does the change reach into modules, packages, or services it shouldn't? Does it import from internal paths of another module?
- Does the change create dependencies between components that were previously independent?
- Does the change weaken an existing boundary by exposing internals, adding pass-through methods, or creating shared mutable state?
- Does the change introduce a new concept that should have its own boundary but doesn't?
- Does data or control flow skip layers or go in the wrong direction?
- **Single Responsibility** — does each module/class have one reason to change? If a change requires modifying a component for two unrelated reasons, the responsibility is too broad.
- **Bounded contexts** — if the codebase has domain boundaries, does the change respect them? Are different models for the same concept properly isolated?
- **Bulkhead isolation** — can failure in this component cascade to others? Are failure domains properly contained?

### Cohesion
- Do the things that change together live together? If this change requires touching files scattered across unrelated packages, it may indicate poor cohesion.
- Is there a **god object** — a class or module that knows too much or does too much?
- **Shotgun surgery** — does a single logical change require modifications in many different modules? That's a cohesion problem.
- Does the change put new behavior in the right place, or is it adding logic to a component that doesn't own that responsibility (**feature envy**)?

### Dependencies
- Do dependencies flow from higher-level to lower-level? Is anything depending upward or sideways? (**Dependency Inversion**)
- Does the change create or extend a **dependency cycle**? The dependency graph must be acyclic.
- **Stable Dependencies** — does the change depend on something less stable than itself? Depend in the direction of stability.
- Could this change work with fewer imports or a narrower interface? (**Interface Segregation** at the module level)
- Does the change introduce new third-party dependencies? Are they justified, maintained, and appropriately scoped?
- Does the change leak a dependency's types or behavior through a public interface, coupling consumers to something they shouldn't know about?
- **Inversion of Control** — is the dependency direction correct? Are high-level components controlling flow, or are low-level details driving architecture?

### Removal
- How many files would need to change to remove this feature? If the answer is more than the feature itself, the boundaries are likely wrong.
- Does the change introduce shared state that other features will come to depend on?
- Does the change hook into core paths (middleware, event loops, initialization) in ways that would be hard to untangle?
- If this turns out to be the wrong approach, what's the cost of backing it out?
- Could this feature be turned off or removed without affecting unrelated functionality?
- **Orthogonality** — are the components independent? Changing this feature shouldn't force changes in unrelated features.

## Do Not Flag

- Internal reorganization within a module that doesn't change its external interface
- Boundary crossings that follow established patterns in the codebase — one more import in an existing pattern isn't a new problem
- Test code reaching into internals — tests have different boundary rules
- Dependencies that are standard for the ecosystem
- Internal utility imports within a module
- Core features that are definitionally permanent — not everything needs to be removable
- Small utilities where removal cost is trivial regardless of structure
- Cohesion issues that predate this change — don't blame new code for old structure unless it's making it worse
