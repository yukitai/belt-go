@str = global [13 x i8] c"Hello, World!"

declare i32 @puts(i8* %0)

define i8* @main() {
entry:
        %0 = alloca i8*
        store i8* getelementptr ([13 x i8], [13 x i8]* @str, i64 0, i64 0), i8** %0
        %1 = load i8*, i8** %0
        %2 = call i32 (i8*) @puts(i8* %1)
        ret i8* %1
}