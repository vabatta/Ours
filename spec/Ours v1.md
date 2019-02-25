## OURS VERSION 1.0
---
### Description
This file describe the file syntax version `1.0` for the Ours application.

Here a general overview of the file:
```txt
/* <syntax version> */
<id>@<color>@<name>

<id>:<type>:<classroom>:<start time>:<end time>
```

### Basic structure
Every document must start with syntax version rule (see below).  
Every line represents a single class hours or a single class name.
The separator for each rule in a single line is done through the character `@` or `:`.  
Each rule is then recognized by the syntax structure of itself.

### Rules
The rules are divided in partitions. Partitions are ordered and each one is marked with an initial 
`<number>: ` before its syntax to identify the position in the rule. 
Those marked with a `+` after the number means that could be repetead more times in a single rule.   
Good practice is to put *Name rule*s before *Class rule*s. Remember that a *Class rule* overwrites 
a *Class rule* written before it if the hours are overlapping.

### Syntax rule
`/* ours@<syntax version> */ : <syntax version> : ^\/\*\s*ours@(\d\.\d)\s*\*\/$`  
This represents the current Ours syntax used inside the file.

### Name rule
`<id>@<color>@<name>`  
The name represents the course information and how to display it in the timetable.

#### Id
`1: <id> : [A-Z0-9]+`  
The id identify the course and is used to identify the class hours.  

#### Color
`2: <color> : [A-Z]+`  
The color is used for coloring the hours in the timetable.

#### Name
`3: <name> : [a-zA-Z0-9\s]+`  
The full name of the course.

### Class rule
`<id>:<type>:<day>:<start>:<end>:<classroom>`  
Represents a single class hour in the timetable.

#### Id
`1: <id> : [A-Z0-9]+`  
The id to identify the course from the *Name rule*s.

#### Type
`2: <type> : LAB|THY`  
The type could be either *LAB* or *THY* meaning laboratories and theories classes respectively.

#### Day
`3: <day> : MON|TUE|WED|THU|FRI`  
The day is when the class is teached.

#### Start
`4: <start> : [0-9]{4}`  
The start time of the class written as full integer number.

#### End
`5: <end> : [0-9]{4}`  
The end time of the class written as full integer number.

#### Classroom
`6: <classroom> : [a-zA-Z0-9\s.-]+`  
The classroom place represeting both the building and number.

### Example
```txt
/* ours@1.0 */
SSR@YELLOW@Systemas y Servicios en Red

SSR:LAB:1G 2S-17:MONDAY:1900:2030
SSR:THY:1E 0.3:TUESDAY:1700:1830
```