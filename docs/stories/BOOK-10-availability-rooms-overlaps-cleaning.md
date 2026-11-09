# BOOK-10 · Availability: rooms, overlaps and cleaning time
Status: approved

**As the** cafe bar
**I want** availability to be based on our rooms, overlapping bookings and cleaning time
**So that** we don't get double bookings and exceed our capacity

## Acceptance criteria
- There are 4 rooms, all of the same type.
- The user does not pick a room; a free room is assigned automatically.
- The 15-minute cleaning buffer is added to every booking and stays reserved (e.g. 10:00–11:00 blocks the room until 11:15).
- Two bookings overlap when their times, including the cleaning buffer, overlap (e.g. 10:00–11:00 and 10:30–11:30 overlap).
- A time slot is available as long as at least one room is free for the whole slot, including its cleaning buffer.
- Held slots (BOOK-8) count as taken.

## Out of scope (demo MVP)
- Multiple room types.
- Changing the number of rooms.

## Open questions
- _(none)_
