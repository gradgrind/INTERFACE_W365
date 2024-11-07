# Data Structures.

I use the term Activity to refer to a unit which can be placed in the timetable. It has a length (number of slots), and various needed resources – Teacher(s), Student Group(s) and Room(s). It has a list of Lesson references, which is the link to the elements in the database it represents. Normally this list will only have a single entry, but allowing multiple entries allows compulsory parallel Lessons to be handled within a single Activity. There is also a list of permissible time slots, which is described later. The Activities are generated afresh after any changes of relevant database data.

To detect clashes, each resource has a vector of Activity references, one entry for each time slot of the week. An empty slot refers to no activity, meaning that the resource is available. There is a special value for a slot which is "not available" to activities. This value cannot be removed.

By also arranging these clash vectors within a vector, access is kept quite simple. Each resource has an index to this outer vector. These indexes are also used in the activity fields – instead of whatever references are used in the primary representation of the activities. In other words, the relevant data structures are rebuilt for the placement and clash monitoring. Essentially, the idea is to use arrays instead of maps to increase access speed.

Taking slots blocked by the resource not being available and those blocked as a result of fixed Lessons, a reduced list of potentially available slots can be built for each Activity. This reduces the search space.
