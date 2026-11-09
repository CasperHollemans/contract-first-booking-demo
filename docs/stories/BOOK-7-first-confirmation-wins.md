# BOOK-7 · First confirmation wins
Status: approved

**As the** cafe bar
**I want** that when two users try to book the last free room for the same slot at (roughly) the same time, the one who confirmed first gets it
**So that** we don't get double bookings and exceed our capacity

## Acceptance criteria
- A conflict only happens when the last free room for that time slot (BOOK-10) is being booked.
- Overlapping slots count as the same slot (see BOOK-10).
- When two users try to book the last free room for the same or an overlapping slot at the same or roughly the same time, the user who confirmed first gets the slot.
- The other user's booking is not made.
- The 5-minute hold from BOOK-8 is an extra prevention layer; this rule still applies on top of it.

## Open questions
- _(none)_
