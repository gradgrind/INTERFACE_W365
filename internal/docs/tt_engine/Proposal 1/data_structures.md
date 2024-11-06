# Data Structures.

I use the term Activity to refer to a unit which can be placed in the timetable. It has a length (number of slots), and various needed resources – Teacher(s), Student Group(s) and Room(s). It also provides a Subject (here for display purposes only).

To detect clashes, each resource has a vector of Activity references. These are accessed as two-dimensional arrays (days x hours). By also arranging the clash arrays (one for teachers, one for groups, one for rooms) as vectors, access is kept quite simple. An empty slot refers to no activity, meaning that the resource is available.

There is a special value for a slot which is "not available" to activities. This value cannot be removed.

All accessing is done by means of 0-based indexes, which means that each resource must also get a corresponding index. These indexes are also used in the activity fields – instead of whatever references are there in the primary representation of the activities. In other words, the relevant data structures are rebuilt for the placement and clash monitoring. Essentially, the idea is to use arrays instead of maps to increase access speed.
