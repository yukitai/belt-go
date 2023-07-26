; ModuleID = './c_llvm/test.c'
source_filename = "./c_llvm/test.c"
target triple = "x86_64-pc-linux-gnu"

define i32 @add(i32 %0, i32 %1) {
; 函数头开始
  %3 = alloca i32, align 4
  %4 = alloca i32, align 4
  store i32 %0, i32* %3, align 4
  store i32 %1, i32* %4, align 4
; 函数头结束
; Add 开始
  %5 = load i32, i32* %3, align 4
  %6 = load i32, i32* %4, align 4
  %7 = add nsw i32 %5, %6
; Add 结束
  ret i32 %7
}