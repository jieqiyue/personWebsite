> 在rust模式匹配中，单纯的下划线和下划线开头的变量有什么区别呢？
```rust
let s = Some(String::from("Hello!"));

if let Some(_s) = s {
    println!("found a string");
}

println!("{:?}", s);
```
这段代码会报错，因为s的所有权转移到了_s上面，但是如果不用_s而是使用_的话，就不会报错了。

> 模式匹配的每一个分支的返回值类型必须相同

每一个分支的最后一句话必须是一个语句，代表要返回的东西。
所以match本身也是一个表达式，可以用来赋值。

```rust
enum Direction {
    East,
    West,
    North,
    South,
}

fn main() {
    let dire = Direction::South;
    match dire {
        Direction::East => println!("East"),
        Direction::North | Direction::South => {
            println!("South or North");
        },
        _ => println!("West"),
    };
}
```
在和这个例子中，所有分支的返回值类型都是相同的。
但是下面例子呢？
```rust
enum Direction {
    East,
    West,
    North,
    South,
}

fn main() {
    let dire = Direction::South;
    match dire {
        Direction::East => println!("East"),
        Direction::North | Direction::South => {
            println!("South or North");
            12
        },
        Direction::West => println!("West"),
    };
}
```
运行结果是：
```rust
Compiling world_hello v0.1.0 (/home/jie/rust/world_hello)
error[E0308]: `match` arms have incompatible types
  --> src/main.rs:14:13
   |
10 | /     match dire {
11 | |         Direction::East => println!("East"),
   | |                            ---------------- this is found to be of type `()`
12 | |         Direction::North | Direction::South => {
13 | |             println!("South or North");
14 | |             12
   | |             ^^ expected `()`, found integer
15 | |         },
16 | |         Direction::West => println!("West"),
17 | |     };
   | |_____- `match` arms have incompatible types

For more information about this error, try `rustc --explain E0308`.
error: could not compile `world_hello` (bin "world_hello") due to 1 previous error
```
所以分支的返回值类型必须相同。

> if let 匹配

```rust
fn main() {
    let v = Some(3u8);
    match v {
        Some(3) => println!("three"),
        _ => (),
    }

    if let Some(3) = v{
        println!("hello");
    }

    if v == Some(3){
        println!("world");
    }
}
```
这个if  let是只用于匹配match中单一的某一种类型的时候用的。
那为什么不直接用一个if来判断不就行了吗？
观察如下情况：
```rust
fn main() {
   
    let five = Some(5);
    
    let six = plus_one(five);
    println!("{:?}", six);
}

fn plus_one(x: Option<i32>) -> Option<i32> {
    if x == Some(i) {  // 这行代码报错
        
    }

    if let Some(i) = x{
        return Some(i + 1);
    }

    return None;
}
```
上面这种情况汇总，if是不能解构出i的值的，只能是一个确定的值。
本质原因是一种叫模式绑定的东西。
模式匹配的另外一个重要功能是从模式中取出绑定的值，例如：
Option的定义是
```rust
enum Option<T> {
    Some(T),
    None,
}
```
这个Option绑定了一个T。当我们需要获取到这个T的时候，就可以用到模式绑定的语法特性。
但是if是不能获取到的。
```rust
fn plus_one(x: Option<i32>) -> Option<i32> {
    match x {
        None => None,
        Some(i) => Some(i + 1),
    }
}
```
在这个代码中，Option中的T被绑定到了i上面。