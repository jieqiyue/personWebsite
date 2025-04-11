从一段代码引起的思考：
```rust
struct Val<T> {
    val: T,
}

impl<T> Val<T> {
    fn value(&self) -> &T {
        &self.val
    }
}


fn main() {
    let x = Val{ val: 3.0 };
    let y = Val{ val: "hello".to_string()};
    println!("{}, {}", x.value(), y.value());
}
```
在这段代码中为什么value函数要返回val的引用呢？
可不可以不返回引用。
不返回引用换句话说，就是要把val的所有权给转移出去。
那么要实现这一点要怎么做呢？
可以想到两种办法，
1. 给这个T类型实现一下Copy trail，实现了之后，这个类型就可以直接Copy了，而不会发生所有权的转移。
2. 转移val的所有权。
那么接下来一一看下，要实现这些东西，代码要怎么写
首先看第一种：
```rust
use std::env::var;

#[derive(Debug, Clone,Copy, PartialEq, Eq)]
struct A;
struct Val<T:Copy> {
    val: T,
}

impl<T:Copy> Val<T> {
    fn value(&self) -> T {
        self.val
    }
}


fn main() {
    let x = Val{ val: 3.0 };
    //let y = Val{ val: "hello".to_string()};
    let z = Val{val:A};
   // println!("{}, {}", x.value(), y.value());
    println!("{}", x.val);
    println!("{:?}", z.val);
}
```
然后来看第二种：
```rust
struct Val<T> {
    val: T,
}

impl<T> Val<T> {
    fn value(self) -> T {
        self.val
    }
}


fn main() {
    let x = Val{ val: 3.0 };
    let y = Val{ val: "hello".to_string()};
    println!("{}, {}", x.value(), y.value());
}
```
在这个代码中，val的所有权被转移了。
在 Rust 中，如果你想编写一个方法，将结构体中某个字段的所有权转移出去，你不可以使用 &mut self 作为参数，因为 &mut self 是一个可变借用，而不是所有权的转移。为了转移所有权，你需要使用 self 作为方法的接收者。