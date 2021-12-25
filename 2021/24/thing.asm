u: +13 +12 +10 -11 +14 +13 +12  -5 +10 +0 -11 -13 -13 -11
q:  +8 +16  +4  +1 +13  +5  +0 +10  +7 +2 +13 +14 +15  +9

[w: 0] [x: 0] [y: 0] [z: 0]
[w: a] [x: 0 + u[i]] [y: 0] [z: 0]
inp w
mul x 0
add x z
mod x 26
div z 1
add x 13
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 8
mul y x
add z y
[w: a] [x: 1] [y: a+8] [z: a+8]

inp w
mul x 0
add x z
mod x 26
div z 1
add x 12
eql x w 
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 16
mul y x
add z y
[w: b] [x: 1] [y: b+16] [z: 26*(a+8)+b+16]
[w: b] [x: 1] [y: b+16] [z: 26*(a+8)+b+16]

inp w
mul x 0
add x z
mod x 26
div z 1
add x 10
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 4
mul y x
add z y
[w: c] [x: 1] [y: c+4] [z: 26*(26*a+208+b+16)+c+4]
[w: c] [x: 1] [y: c+4] [z: 26*(26*(a+8)+b+16)+c+4]

inp w
mul x 0
add x z
mod x 26
div z 26
add x -11
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 1
mul y x
add z y
(pick d == 2 & c == 9) (min: c=8;d=1)
[w: d] [x: 0] [y: 0] [z: (26*a+208+b+16)]
[w: 2] [x: 0] [y: 0] [z: 26*(a+8)+b+16)]
[w: d] [x: c-7] [y: c+4] [z: 26*(26*(a+8)+b+16)+c+4]

inp w
mul x 0
add x z
mod x 26
div z 1
add x 14
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 13
mul y x
add z y
[w: e] [x: 1] [y: e+676*a+26*b+5837] [z: 676*a+5824+26*b]
[w: e] [x: 1] [y: e+13] [z: 26*(26*a+208+b+16)+e+13]
[w: e] [x: 1] [y: e+13] [z: 26*(26*(a+8)+b+16)+e+13]

inp w
mul x 0
add x z
mod x 26
div z 1
add x 13
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 5
mul y x
add z y
[w: f] [x: 1] [y: f+5] [z: 17576*a+151424+676*b+f+5]
[w: f] [x: 1] [y: f+5] [z: 26*(26*a+208+b+16)+e+13+f+5]
[w: f] [x: 1] [y: f+5] [z: 26*(26*(26*(a+8)+b+16)+e+13)+f+5]

inp w
mul x 0
add x z
mod x 26
div z 1
add x 12
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 0
mul y x
add z y
// don't pick f = 9
[w: g] [x: 1] [y: g] [z: 26*(17576*a+151424+676*b+14)+g]
[w: g] [x: e+13+f+5] [y: f+5] [z: 26*(26*a+208+b+16)+e+13+f+5]
[w: g] [x: 1] [y: g] [z: 26*(26*(26*(26*(a+8)+b+16)+e+13)+f+5)+g]

inp w
mul x 0
add x z
mod x 26
div z 26
add x -5
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 10
mul y x
add z y
pick (g=9 h = 4) min(g=6;h=1)
[w: h] [x: 0] [y: 0] [z: 17576*a+151424+676*b+14]
[w: h] [x: 0] [y: 0] [z: 26*(26*(26*(a+8)+b+16)+e+13)+f+5]
[w: h] [x: g-5] [y: g] [z: 26*(26*(26*(26*(a+8)+b+16)+e+13)+f+5)+g]

inp w
mul x 0
add x z
mod x 26
div z 1
add x 10
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 7
mul y x
add z y
[w: i] [x: 1] [y: i+7] [z: 26*(17576*a+151424+676*b+14)+i+7]
[w: i] [x: 1] [y: i+7] [z: 26*(26*(26*(26*(a+8)+b+16)+e+13)+f+5)+i+7]

inp w
mul x 0
add x z
mod x 26
div z 26
add x 0
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 2
mul y x
add z y
pick i = 2, j = 9 (min: i=1,j=8)
[w: j] [x: 0] [y: 0] [z: 17576*a+151424+676*b+14]
[w: 9] [x: 0] [y: 0] [z: 26*(26*(26*(a+8)+b+16)+e+13)+f+5]
[w: j] [x: i+7] [y: i+7] [z: 26*(26*(26*(26*(a+8)+b+16)+e+13)+f+5)+i+7]

inp w
mul x 0
add x z
mod x 26
div z 26
add x -11
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 13
mul y x
add z y
pick k = 3; f=9 (min: f=7;k=1)
[w: k] [x: 0] [y: 0] [z: 676*a+5824+26*b]
[w: 3] [x: 0] [y: 0] [z: 26*(26*(a+8)+b+16)+e+13]
[w: k] [x: f-6] [y: 0] [z: 26*(26*(26*(a+8)+b+16)+e+13)+f+5]

inp w
mul x 0
add x z
mod x 26
div z 26
add x -13
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 15
mul y x
add z y
e & l = 9 (min: e & l = 1)
[w: l] [x: 1] [y: l+15] [z: 26*(26*a+224+b)+l+15]
[w: 9] [x: 0] [y: 0] [z: 26*(a+8)+b+16]
[w: l] [x: e] [y: 0] [z: 26*(26*(a+8)+b+16)+e+13]

inp w
mul x 0
add x z
mod x 26
div z 26
add x -13
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 14
mul y x
add z y
// don't pick m = 9, l =7
actual pick b=6, m=9 (min: b=1, m=4)
[w: m] [x: 0] [y: 0] [z: 26*a+224+b]
[w: 9] [x: 0] [y: 0] [z: (a+8)]
[w: m] [x: b+3] [y: 0] [z: 26*(a+8)+b+16]

inp w
mul x 0
add x z
mod x 26
div z 26
add x -11
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y 9
mul y x
add z y
//don't pick (b = 4, n = 9)
actual pick a=9, n=6 (min: a=4, n=1)
[w: n] [x: 0] [y: 0] [z: a+16]
[w: 6] [x: 6] [y: 0] [z: 0]
[w: n] [x: a-3] [y: 0] [z: (a+8)]
