---
title: OURS 2.0 SPECIFICATION
author: Valerio Battaglia
date: 26.09.2018
version: 2.0
lang: en
---

### Description
This file describe the file syntax version `2.0` for the Ours application.

Here a general overview of the file:
```txt
/* <syntax version> */

<id>@<color>@<name>

<id>:<icon>:<location>:<start time>:<end time>
```

### Basic structure
Every document must start with syntax version rule (see below).  
Every line represents a single activity or a single slot's activity.
The separator for each rule in a single line is done through the characters `@` and `:`.

### Rules
The rules are divided and explained using partitions. A partition is a subpart of the rule: they are ordered and each one is marked with an initial `<number>: ` before its syntax definition to identify its position in the rule.  
Good practice is to put *Activity rule*s before *Slot rule*s. Remember that a *Slot rule* overwrites 
a previously written *Slot rule* if the hours are overlapping.

<!-- \pagebreak -->

### Syntax rule
`/* ours@<syntax version> */ : <syntax version> : \/\*\s*ours@(\d\.\d)\s*\*\/`  
This represents the current Ours syntax used inside the file. Syntax version for this document: 2.0

### Activity rule
`<id>@<color>@<name>`  
The name represents the activity information and how to display it in the timetable.

#### Id
`1: <id> : [A-Z0-9]+`  
The id identify the activity and is used to identify its slots. Must be unique in the whole file.  

#### Color
`2: <color> : [A-Z]+ | #[a-fA-F0-9]{6}#[a-fA-F0-9]{6}`  
The color of the slots in the timetable, either represented as built-in color or as a pair of hex colors for background and foreground (e.g. #ffffff#000000).

#### Name
`3: <name> : [a-zA-Z0-9\s]+`  
The full name of the activity.

### Slot rule
`<id>:<icon>:<location>:<day>:<start>:<end>`  
Represents a single slot's activity in the timetable.

#### Id
`1: <id> : [A-Z0-9]+`  
The id to identify the course from the *Activity rule*s.

#### Icon
`2: <icon> : [-a-zA-Z0-9\s]*`  
The icon to draw in the upper right corner of the slot.

#### Location
`3: <location> : [a-zA-Z0-9\s.-]+`  
The location represeting where the slot is done.

#### Day
`4: <day> : MON|TUE|WED|THU|FRI|01|02|03|04|05`  
The day of the slot. The enumeration of the week starts from 01 for Monday to 05 for Friday.

#### Start
`5: <start> : (?:0[089]|1[0-9]|2[0-2])[0-5][0-9]`  
The start time of the slot between 08:00 and 22:00 in 24h format, written as full integer number without semicolon `:`.

#### End
`6: <end> : (?:0[089]|1[0-9]|2[0-3])[0-5][0-9]`  
The end time of the slot between 08:00 and 23:00 in 24h format, written as full integer number without semicolon `:`.

### Comment rule
`#<text>`  
A comment in the `ours` file which will be ignored by the parser.

#### Text
`1: <text> : .+`
The comment content of the file.

### Example
```txt
/* ours@2.0 */

# Activities
SSR@YELLOW@Systemas y Servicios en Red
TSR@#e74c3c#efefef@Tecnología de Sistemas de Información en la Red

# Slots
SSR:FLASK:1G 2S-17:MON:1900:2030
SSR:BOOK:1E 0.3:TUE:1700:1830

TSR:FLASK:1B DSIC 5:WED:1030:1200
TSR:BOOK:1E 2.0:WED:1300:1430
TSR:BOOK:1E 2.0:FRI:0930:1100
```
