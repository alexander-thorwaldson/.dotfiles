# review economy

Is this doing more than it needs to?

## What to Look For

### Complexity
- Layers of indirection that don't buy anything. If you can't explain why a layer exists in one sentence, it probably shouldn't. (**Abstraction Inversion** — is a simple operation implemented on top of a complex abstraction?)
- Control flow that requires tracking 4+ levels of context. Could early returns flatten it?
- Generic solutions to specific problems. Plugin systems for one plugin. Factory patterns for one type. Configuration for values that never change.
- Wrappers, adapters, and middleware that just pass through. Every hop costs comprehension.
- Could the same result be achieved with less machinery? Three lines of straightforward code is better than a clever one-liner or a premature abstraction. (**Worse Is Better** — simplicity of implementation matters)
- **Gall's Law** — is the change trying to build a complex system from scratch rather than evolving from something simple that works?

### YAGNI
- Parameterizing things that only have one value, making things configurable that nobody will configure, accepting interfaces when only one implementation exists.
- Code paths for scenarios that don't exist yet. "We might need this" is not a requirement.
- Hooks, plugins, event systems, middleware chains built for future extensibility that has no concrete plan.
- Building for scale/load/concurrency that the system doesn't face and has no timeline to face. (**Premature Optimization** — don't optimize before profiling)
- Extra options, extra modes, extra flexibility that nobody asked for. (**Zawinski's Law** — resist the natural tendency toward feature bloat)

### Proportionality
- **Rule of Three** — is the change extracting an abstraction prematurely? Duplication is acceptable twice. The third time, refactor. Not the first time.
- Is the architectural investment proportional to the problem? A 200-line framework for a 20-line problem is a net negative.
- **Lehman's Laws** — is the change actively reducing complexity, or is it adding to the natural growth of complexity that comes with evolution? Every change that doesn't reduce complexity increases it.
- **Pareto awareness** — is the effort focused on the part of the system that matters most, or is it gold-plating something that rarely runs?

## Do Not Flag

- Complexity that is inherent to the problem domain — some things are genuinely complex
- Abstractions that have multiple consumers — if three callers use it, the abstraction is justified
- Standard patterns used in their standard way
- Abstractions that serve current, concrete needs — even if they also happen to be extensible
- Standard patterns that are cheap to implement — an interface with one implementation is fine if it's at a natural boundary
- Known upcoming work — if the roadmap has a concrete plan for a second consumer, designing for it now is reasonable
- Performance work backed by profiling data — that's not premature optimization, that's engineering
