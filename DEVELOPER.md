
# Assumptions

## Public api should always work with region position vs termbox position
This is the purpose of this library, to remove the need to calculate absolute positioning.
However, internal api will need to deal with termbox position from time to time, and that is okay.

## All array are stored as YX instead of XY.
There are many use cases for storing cells as `cells[Y][X]`.

For example, this allow you define menu as rows, which can then be split into cells, much more easily.
It then allows us to do write "borders" like

```
rows[0] = "XXXX"
rows[1] = "X  X"
rows[1] = "X XX"
rows[1] = "XXX "
```

and remove the need to "flip" the x y position of the cells.

However, this is only for internal representation and for params that are 2D array.

For Param where x and y are required, i.e. `SetCell`, the ordering should still be x, y.
