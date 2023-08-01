@$const_0 = global [7 x i8] c"Hello, "
@$const_1 = global [7 x i8] c"World!\0A"

declare i32 @puts(i8* %0)

declare i8* @__bl_str_connect(i8* %0, i8* %1)

define i8* @get_string() {
entry:
        %0 = alloca i8*
        %1 = call i8* @__bl_str_connect(i8* getelementptr ([7 x i8], [7 x i8]* @$const_0, i64 0, i64 0), i8* getelementptr ([7 x i8], [7 x i8]* @$const_1, i64 0, i64 0))
        store i8* %1, i8** %0
        %2 = alloca i64
        store i64 114, i64* %2
        %3 = alloca i64
        store i64 514, i64* %3
        %4 = alloca i64
        %5 = load i64, i64* %2
        %6 = mul i64 %5, 1000
        %7 = load i64, i64* %3
        %8 = add i64 %6, %7
        store i64 %8, i64* %4
        %9 = load i8*, i8** %0
        ret i8* %9
}

define i1 @test_tyinfer() {
entry:
        %0 = alloca i64
        %1 = alloca fp128
        %2 = alloca i1
        %3 = alloca i64
        %4 = alloca fp128
        %5 = alloca i64
        %6 = load i64, i64* %3
        %7 = load i64, i64* %0
        %8 = add i64 %6, %7
        store i64 %8, i64* %5
        %9 = alloca fp128
        %10 = load fp128, fp128* %4
        %11 = load fp128, fp128* %1
        %12 = fadd fp128 %10, %11
        store fp128 %12, fp128* %9
        %13 = alloca i8*
        %14 = alloca i8*
        %15 = alloca i8*
        %16 = load i8*, i8** %13
        %17 = load i8*, i8** %14
        %18 = call i8* @__bl_str_connect(i8* %16, i8* %17)
        store i8* %18, i8** %15
        %19 = load i1, i1* %2
        ret i1 %19
}