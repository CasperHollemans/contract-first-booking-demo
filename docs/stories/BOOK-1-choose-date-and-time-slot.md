# BOOK-1 · Choose a date and time slot
Status: approved

**As a** user
**I want** to choose a date and time slot for a room
**So that** I can make use of the room for a prespecified time period

## Acceptance criteria
- The user can select a date from today up to and including 30 days ahead (e.g. on 9 Nov, the last bookable date is 9 Dec).
- The time slot is at least 30 minutes, in 15-minute increments.
- Start and end times fall on the 15-minute grid (:00, :15, :30, :45), e.g. 10:15–11:00.
- The time slot must fall within opening hours, 08:00–19:00, every day of the week (weekends included).
- A booking may end exactly at 19:00; its 15-minute cleaning time may run past closing.
- On today's date, a time slot can only be selected and confirmed before it starts (e.g. at 10:07, 10:15 can be selected but 10:00 cannot; from 10:15 onwards, 10:15 can no longer be selected or confirmed).

## Examples
_Approved._ Each example is one spec in `acceptance/specs/book1_choose_date_and_time_slot_test.go`, grouped per acceptance criterion. Today is Monday 9 November, 00:00, unless stated otherwise.

**Date range** (`TestBOOK1_DateRange`)
- Exactly 30 days ahead (9 December, 10:00–11:00) can be selected.
- 31 days ahead (10 December, 10:00–11:00) cannot be selected.
- A past date (8 November, 10:00–11:00) cannot be selected.
- A weekend day (Saturday 14 November, 10:00–11:00) can be selected.

**Opening hours** (`TestBOOK1_OpeningHours`, on 10 November)
- Starting at opening (08:00–08:30) can be selected.
- Ending exactly at closing (18:30–19:00) can be selected.
- The whole day (08:00–19:00) can be selected.
- Starting before opening (07:45–08:30) cannot be selected.
- Ending after closing (18:45–19:15) cannot be selected.

**Length and grid** (`TestBOOK1_LengthAndGrid`, on 10 November)
- 30 minutes (10:00–10:30) can be selected.
- 45 minutes (10:00–10:45) can be selected.
- 29 minutes (10:00–10:29) cannot be selected.
- 15 minutes (10:00–10:15) cannot be selected.
- A start time off the grid (10:07–10:45) cannot be selected.
- An end time off the grid (10:15–10:50) cannot be selected.

**Today** (`TestBOOK1_Today`, on 9 November)
- At 10:07, 10:15–11:00 hasn't started and can be selected.
- At 10:07, 10:00–11:00 has already started and cannot be selected.
- At exactly 10:15, 10:15–11:00 starts right now and cannot be selected.

_Moved to BOOK-9:_ "Selected 10:15–11:00 at 10:10 and confirmed at exactly 10:15: the booking cannot be confirmed." Confirming is built in BOOK-9, so that's where this example is proven.

## Open questions
- _(none)_
