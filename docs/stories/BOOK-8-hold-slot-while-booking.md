# BOOK-8 · Hold a slot while booking
Status: approved

**As the** cafe
**I want** a time slot to be reserved for 5 minutes once a user selects it, so they can finish their booking
**So that** another user doesn't start booking the same time slot and get disappointed later

## Acceptance criteria
- When a user selects a time slot, it is held for them for 5 minutes.
- A held slot counts as taken for other users (e.g. with all 4 rooms held or booked, the slot shows as unavailable).
- The user sees a 5-minute countdown while their slot is held.
- If the user changes their selection, the previous hold is released immediately.
- If the user doesn't finish their booking within 5 minutes, the hold is released and the slot becomes available to other users again.
- When the 5 minutes run out, the user is shown that the slot is now available for others to book.
- The hold is an extra layer of prevention on top of BOOK-7; BOOK-7 ("the one who confirmed first wins") still applies.

## Out of scope (demo MVP)
- Re-selecting the same slot after the hold has run out.

## Open questions
- _(none)_
