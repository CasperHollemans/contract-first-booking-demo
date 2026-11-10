# Testing strategy
Status: draft

## Principles
1. **One locked contract.** `acceptance/specs/` is the contract. It changes only with explicit approval, never during Build.
2. **Every behaviour gets one spec in the DSL.** Each spec is written in business language and traces back to a story example.
3. **Each spec lists the drivers that prove it.** Choose the cheapest driver that can prove the behaviour. A spec may list more than one driver when each adds proof the others can't give.
4. **Every story has at least one UI spec.** This is enforced: the run fails otherwise.

## Layers
- `acceptance/specs/`: the contract. One `TestBOOKn_<Criterion>` function per acceptance criterion, with one `Spec` per behaviour.
- `acceptance/dsl/`: business-language steps (`GivenToday`, `WhenUserChooses`, …) and `Spec`.
- `acceptance/drivers/`: how the steps reach the system.

## Writing a spec
```go
func TestBOOK1_DateRange(t *testing.T) {
	Spec(t, "When the user chooses exactly 30 days ahead, expect the slot can be selected", func(b *Booking) {
		b.GivenToday("9 November")
		b.WhenUserChooses("9 December", "10:00–11:00")
		b.ThenSlotCanBeSelected()
	}, Domain, API, UI)
}
```

- `Spec` runs `t.Run(behaviour)` and, inside it, one subtest per driver: `TestBOOK1_DateRange/When_…/UI`.
- Each driver run gets a fresh `Booking`, with a fresh system and a fresh fake clock.
- The story is the prefix of the test name: `TestBOOK1_DateRange` belongs to `BOOK1`.
- Examples are in 2026 and in the cafe's time zone (UTC). Slots take an en dash or a hyphen.

## Drivers
| Driver | Runs against | Speed | Use it for |
|---|---|---|---|
| Domain | the service layer (`internal/service`), with an in-memory repository and the fake clock; no HTTP | µs | rules and logic: boundaries, combinations, time-based behaviour |
| API | the real server over HTTP (`httptest`), with the fake clock | ms | errors, wiring, anything concurrent (e.g. first confirmation wins) |
| UI | the real server through headless Chrome, with the fake clock | s | the happy path only |

Postgres comes in when a story stores data (BOOK-9, BOOK-10). Then the API and UI drivers run against real Postgres, and the Domain driver gets an in-memory repository.

## Choosing a driver
Choose per behaviour, not per story. List the behaviours in a story, then pick the cheapest driver that can prove each one. A story gets as many UI specs as it has behaviours that only a browser can prove, and at least one.

- Could a rule or a combination of rules break it? Use the **Domain** driver.
- Does it depend on HTTP, validation, storage, transactions or races? Add the **API** driver.
- Is it the journey a person actually takes? Add the **UI** driver.

A rule rejection (e.g. a past date) is a rule, so Domain proves it. Add API once per story to prove that a rejection reaches the client correctly (`200 {"selectable": false}`). Malformed input (`400`) is not a business behaviour: controller tests cover it.

## Guards
- **Never ran.** A `Spec` with no drivers fails with `this spec never ran: add Domain, API or UI`.
- **Story needs UI.** `TestMain` records which drivers ran per story. After the run, it fails any story that ran specs but none through the UI. The guard is skipped when the run is filtered with `-run` or `-skip`, because a partial run can't tell a missing UI spec from a filtered one.

## BOOK-1
| Driver | Specs |
|---|---|
| Domain | all 18 |
| API | exactly 30 days ahead (selectable), past date (not selectable) |
| UI | exactly 30 days ahead (the happy path) |

The confirm-at-start-time example moved to BOOK-9, where confirming exists.

## Open questions
- How does Postgres run in tests: docker compose or testcontainers-go? (Not needed until a story stores data, BOOK-9 / BOOK-10.)
- Should a repository contract test suite run against both the in-memory repository and Postgres, so the Domain driver can't drift?
