# Post's Correspondence Problem

Emil, a Polish mathematician, sent a simple puzzle by post to his British friend, Alan. Alan sent a reply saying he didn’t have an infinite amount of time he could spend on such non-essential things. Emil modified his puzzle (making it a bit more restricted) and sent it back to Alan. Alan then solved the puzzle.

Here is the original puzzle Emil sent: given a sequence of pairs of strings (a1,b1),(a2,b2),…,(ak,bk), find a non-empty sequence s1,s2,…,sm such that the following is true:

as1as2…asm=bs1bs2…bsm
where as1as2… indicates string concatenation. The modified puzzle that Emil sent added the following restriction: for all i≠j, si≠sj.

# How it works

For each case, the program displays the case number followed by the sequence found (if it is possible to form one) or “IMPOSSIBLE” (if it is not possible to solve the problem). 

If it is possible but there are multiple sequences, the program will prefer the shortest one (in terms of the number of characters output). If there are multiple shortest sequences, will be choose the one that is lexicographically first.

# Usage

```sh
go run .\main.go
```

The program will read the input, but it can be modify to read a file.

# Example input/output

## Sample input

    5
    are yo
    you u
    how nhoware
    alan arala
    dear de
    8
    i ie
    ing ding
    resp orres
    ond pon
    oyc y
    hello hi
    enj njo
    or c
    3
    efgh efgh
    d cd
    abc ab
    3
    a ab
    b bb
    c cc

## Sample output

    Case 1: dearalanhowareyou
    Case 2: ienjoycorresponding
    Case 3: abcd
    Case 4: IMPOSSIBLE

# TODO

The program was designed based on this [master degree research](https://webdocs.cs.ualberta.ca/~games/PCP/thesis/pcp.pdf) and there are some rules that was not implemented:

1. One importante point in this problem is the depth of the solving (how deep it can go to salve the problem). Currently the program has a fixed depth.
2. There are at least 3 rules to check if the instance is solvable, but was implemented just 2.
3. The performance can be a problem in some cases, so is importante to improve the performance of the iterations.