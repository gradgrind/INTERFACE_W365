# Timetable Generation

There are two distinct primary areas of concern here. Firstly it should be possible to place "activities" manually, getting information as to where an activity can be placed without conflicts, and what conflicts there are when trying to place it in currently blocked time slots. A second desirable feature is to have the activities placed automatically taking the various constraints into account. These two approaches can share some basic framework, but they also differ in a number of points.

## The Timetable and Activities

In these discussions it will be assumed that we are dealing with a timetable covering a week (which may have any number of days) and that each day has a certain number of time slots into which activities may be placed. To keep it simple, all days will have the same number of slots – they don't all need to be filled. An activity can cover any number of slots in one day, the number of slots being referred to as the length of the activity.

An activity has a subject, which is a label for what is happening in the time slot(s) covered by the activity. The activities will generally represent lessons, in which case the subject will be the name of a taught course, say Mathematics, English or Sport. An activity can, however also represent other occurrences, such as conferences. There should always be an abbreviation (say, "Ma", "En" or "Sp"), which can be used in the displays and printouts to allow a dense representation of the timetable.

An activity occupies zero or more teachers, zero or more student groups and zero or more rooms. These can be regarded (when considering the technical aspects of timetable construction) as "resources", each of which can only be used once in any given time slot.

Students are divided into classes, and groups within the classes. A class may be divided in different ways in different subjects. For a further discussion of this point, see "Atomic Groups". A teacher may of course not be divided. It is conceivable that certain rooms could be divisible, but this possibility will not be covered here.

I think starting with a description of manual placement is probably sensible, as it can offer a fairly gentle introduction to the theme without getting too bogged down in questions of algorithmic efficiency.

## Manual Activity Placement

### Hard Constraints

There are certain constraints which should make it impossible to place an activity in a particular slot. One of these is mentioned above, there may be no conflict of resources. This is a "hard" constraint, which must be respected. It may, however, be permissible to remove other activities from a time slot, so that the conflict no longer exists.

There may also be certain constraints that can't be resolved. For example, a class may have lessons only in the morning, or a teacher is available only at particular times. Or some activities may be fixed in particular time slots, so that their "resources" would then not be available for other activities.

By keeping track of these restrictions on the resources, the program can help the user in placing the activities. If an activity is selected, the slots where it can be placed without conflict could be highlighted green (say). Slots which would be available if one or more other activities were removed might be highlighted orange (for example). By some simple action it should also be possible to determine which activities need to be removed in order to free up a particular slot. Some further action can then cause the placement to be performed, removing conflicting activities as necessary.

In some cases there may be a room conflict which could be resolved by choosing a different room for some activity or other. See "Handling Rooms" for more on this rather complicated area.

#### Additional hard constraints

There are other constraints which may be implemented as hard requirements, for example the following ones:

 - **Different Days**: The activities of a set with this constraint must all lie on different days. This would normally be used to ensure that two lessons in the same subject do not occur on the same day. Under some circumstances a softer version might be acceptable.
 
 - **Parallel Activities**: This constraint allows multiple activities to be coupled: forced to start in the same slot. It can be useful if certain lessons must always take place at the same time as each other. Again, a softer version may also be useful.
