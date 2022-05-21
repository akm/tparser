# Test for same idents in some units

## Output

```
-- Unit4 Initialization----
Foo@Unit4
Bar@Unit4
Project1@Unit4
-- Unit1 Initialization----
Foo@Unit1
Bar@Unit1
Unit2@Unit1
-- Unit2 Initialization----
Foo@Unit2
Bar@Unit2
Unit1@Unit2
-- Unit3 Initialization----
Foo@Unit2
Bar@Unit2
Project1@Unit2
-- Project1 ----
Foo@Unit4
Bar@Project1
Unit2@Unit1
Unit1@Unit2
Bar@Project1

-- Unit3 Finalization----
Foo@Unit2
Bar@Unit2
Project1@Unit2

-- Unit2 Finalization----
Foo@Unit2
Bar@Unit2
Unit1@Unit2

-- Unit1 Finalization----
Foo@Unit1
Bar@Unit1
Unit2@Unit1

-- Unit4 Finalization----
Foo@Unit4
Bar@Unit4
Project1@Unit4
```
