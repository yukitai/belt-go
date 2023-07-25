# the belt programming language

### modern, simple and fast

```rust
struct MyStruct {
    a: int,
    b: float,
}

enum Option {
    Some('a),
    None,
}

fn unwrap_option(opt: Option<'a>) -> 'a {
    match opt {
        Some(v) => v,
        None => panic("unwrap a None variant"),
    }
}

fn main() {
    let my_struct = MyStruct { a: 0, b: 0. }
    let MyStruct { a, b } = my_struct
    let some = Option::Some(a)
    let none: Option<int> = Option::None
    print_int(unwrap_option(some))
    print_int(unwrap_option(none)) // will panic
}
```