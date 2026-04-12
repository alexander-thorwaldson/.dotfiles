# review patterns

Does this follow what's established?

## What to Look For

### Consistency
- Does the change solve a problem that's already been solved elsewhere in the codebase? If so, does it use the established approach or invent a new one? (**DRY** — at the pattern level, not just code duplication)
- Does the change follow the codebase's conventions for file organization, error handling, logging, configuration, and similar concerns?
- Does this introduce a second way of doing something? Two ways means every future developer has to choose, and half will choose wrong.
- If the change intentionally breaks from convention, is there a clear reason? New patterns are fine when the old one is genuinely wrong — but it should be explicit, not accidental.
- **Composition over Inheritance** — if the change uses inheritance, would composition be clearer and less coupled? Deep inheritance hierarchies are a pattern smell.
- **Fail Fast** — does the error handling strategy match the codebase convention? Are errors caught early and surfaced visibly, or swallowed silently?

### Naming
- Does each name accurately describe what the thing does? A function called `validate` that also transforms data is misnamed.
- Are names at the right level of abstraction? Too specific ties the name to an implementation. Too vague obscures the purpose.
- Does the change use the same terminology as the rest of the codebase for the same concepts?
- If a module, class, or function can't be named clearly, its responsibility probably isn't clear. The naming problem is a design problem.
- Are there names that suggest one thing but do another?

### Single Source of Truth
- Does the change duplicate data, configuration, or state that already has a canonical location? Every piece of knowledge should have one authoritative representation.
- Are there constants, types, or schemas defined in multiple places that should be shared?
- If the change introduces a new source for something that already exists, which is authoritative?

### Conformance
- Use `/kb article find` to locate KB articles relevant to the changed code. Does the change follow the patterns and conventions documented there?
- Check the `reference/` topic for decision log entries that apply. Is the change consistent with prior architectural decisions?
- Are there KB articles the author should have read before writing this change? Surface them.
- Does the change introduce patterns or concepts that aren't in the KB but should be? Flag these for future KB updates.
- Does the change undo or work against something the KB explicitly documents as intentional?

### Substitutability
- **Liskov Substitution** — if the change introduces or modifies a subtype, can it be used anywhere the parent type is expected without breaking correctness?
- **Protected Variations** — are points of likely change or instability wrapped behind stable interfaces?

## Do Not Flag

- Code style and formatting — that's linters and the engineer's concern
- Patterns that are already inconsistent across the codebase — don't blame this change for a pre-existing mess unless it's making it worse
- New patterns in new modules — if there's no established way of doing something in this area, the change is setting the convention
- Names that follow existing codebase convention even if you'd choose differently
- Local variable names in short-lived scopes
- Domain-standard terminology even if not self-explanatory to outsiders
- Changes in areas the KB doesn't cover — no documentation doesn't mean non-conformance
- Intentional departures from documented patterns when the change includes rationale
- KB articles that are themselves outdated — if the code is right and the docs are wrong, that's a KB maintenance issue
- Inheritance that is idiomatic for the framework in use
