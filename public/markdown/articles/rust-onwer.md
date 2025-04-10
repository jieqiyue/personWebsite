# 所有权
> 什么样子的数据才会有所有权

Rust中数据可以分为基本类型和非基本类型。
基本类型可以有如下这些：
- 所有整数类型，比如 u32
- 布尔类型，bool，它的值是 true 和 false
- 所有浮点数类型，比如 f64
- 字符类型，char
- 元组，当且仅当其包含的类型也都是 Copy 的时候。比如，(i32, i32) 是 Copy 的，但 (i32, String) 就不是
- 不可变引用 &T ，例如转移所有权中的最后一个例子，但是注意: 可变引用 &mut T 是不可以 Copy的

可以通过以下两个方面来理解这些东西的所有权。

1. 因为这些类型的大小是固定的，所以在编译期就已经知道大小了。所以没有必要放到堆中。直接在栈上分配就行了。那对于这些类型的赋值
```rust
let x = 5;
let y = x;

println!("x = {}, y = {}", x, y);
```
就不会发生所有权的转化。也就是说所有权是针对于堆上分配的来说的。

2. 这些类型都实现了Copy这种特性。所有实现了Copy特征的类型，在赋值的时候就会进行拷贝了。任何基本类型的组合可以 Copy ，不需要分配内存或某种形式资源的类型是可以 Copy 的。

> 所有权转移是怎么样的

每一个变量都是有作用域的。并且在离开这个作用域的时候，会调用Drop方法来释放这个变量所指代的内存。那么在这个作用域中，如果内存的所属权从一个变量转移到了另外一个变量。那就是发生了所有权的转移。所以，这也得出了以下的先决条件。
1. Rust 中每一个值都被一个变量所拥有，该变量被称为值的所有者
2. 一个值同时只能被一个变量所拥有，或者说一个值只能拥有一个所有者
3. 当所有者(变量)离开作用域范围时，这个值将被丢弃(drop)

这三条是所有权的规则。
那在实际发生所有权转移的时候，内存中发生了什么事情呢？
![在这里插入图片描述](https://i-blog.csdnimg.cn/direct/a155c22e31e64395884777a3c1cd9ddc.png)
```rust
let s1 = String::from("hello");
let s2 = s1;

println!("{}, world!", s1);  // 报错，因为s1此时已经被转移了所有权，所以它不再能访问堆上的那个之前的内存地址了
```
上面的s1，s2其实是存在栈中的。当s1赋值给s2的时候，就发生了所有权的转移。
可以从图里看出来发生了什么。

> 函数调用和函数返回

将值传递给函数，一样会发生 移动 或者 复制，就跟 let 语句一样，下面的代码展示了所有权、作用域的规则：
```rust
fn main() {
    let s = String::from("hello");  // s 进入作用域

    takes_ownership(s);             // s 的值移动到函数里 ...
                                    // ... 所以到这里不再有效

    let x = 5;                      // x 进入作用域

    makes_copy(x);                  // x 应该移动函数里，
                                    // 但 i32 是 Copy 的，所以在后面可继续使用 x

} // 这里, x 先移出了作用域，然后是 s。但因为 s 的值已被移走，
  // 所以不会有特殊操作

fn takes_ownership(some_string: String) { // some_string 进入作用域
    println!("{}", some_string);
} // 这里，some_string 移出作用域并调用 `drop` 方法。占用的内存被释放

fn makes_copy(some_integer: i32) { // some_integer 进入作用域
    println!("{}", some_integer);
} // 这里，some_integer 移出作用域。不会有特殊操作
```
你可以尝试在 takes_ownership 之后，再使用 s，看看如何报错？例如添加一行 println!("在move进函数后继续使用s: {}",s);。

同样的，函数返回值也有所有权，例如:
```rust
fn main() {
    let s1 = gives_ownership();         // gives_ownership 将返回值
                                        // 移给 s1

    let s2 = String::from("hello");     // s2 进入作用域

    let s3 = takes_and_gives_back(s2);  // s2 被移动到
                                        // takes_and_gives_back 中,
                                        // 它也将返回值移给 s3
} // 这里, s3 移出作用域并被丢弃。s2 也移出作用域，但已被移走，
  // 所以什么也不会发生。s1 移出作用域并被丢弃

fn gives_ownership() -> String {             // gives_ownership 将返回值移动给
                                             // 调用它的函数

    let some_string = String::from("hello"); // some_string 进入作用域.

    some_string                              // 返回 some_string 并移出给调用的函数
}

// takes_and_gives_back 将传入字符串并返回该值
fn takes_and_gives_back(a_string: String) -> String { // a_string 进入作用域

    a_string  // 返回 a_string 并移出给调用的函数
}
```
所有权很强大，避免了内存的不安全性，但是也带来了一个新麻烦： 总是把一个值传来传去来使用它。 传入一个函数，很可能还要从该函数传出去，结果就是语言表达变得非常啰嗦，幸运的是，Rust 提供了新功能解决这个问题。

> 本质

其实这个过程有没有像是编译器通过一个栈来控制了一个变量的生命周期。
比如说一个在函数体内创建的变量，它的引用是存在栈上的。那么这个栈在这次函数调用结束之后就会被回收。那么在这个函数体中创建的堆上的变量，到底要不要继续使用呢？
如果要继续使用那就得程序员自己将它返回到外面（函数被调用处）。
这个栈其实是操作系统去实现的，原本操作系统并没有清理堆区的任务的，它仅仅是创建栈帧，然后压栈，然后销毁栈。但是我们的Rust编译器多赋予了它一个功能，就是清理栈帧的同时，将不用堆内存也一块清理了。
所以对于c++这种语言来说，它没有gc，也就是说得程序员去清理堆区。
对于go语言，它有gc，直接go语言的垃圾收集器就把这个堆内存清理了。
所以其实就是一个堆内存的管理问题。如果堆内存处理不当，很容易内存泄漏。
# 借用
> 可变引用

如果要一直考虑所有权的问题的话，写代码会很麻烦。有没有简单的方式呢？
```rust
fn main() {
    let x = 5;
    let y = &x;

    assert_eq!(5, x);
    assert_eq!(5, *y);
}
```
上面的代码中，y是x的引用。
```rust
fn main() {
    let s1 = String::from("hello");

    let len = calculate_length(&s1);

    println!("The length of '{}' is {}.", s1, len);
}

fn calculate_length(s: &String) -> usize {
    s.len()
}
```
那么，在这个代码里面就不会发生所有权的转移。因为String的所有权并没有发生转移，因为传给函数的参数是一个引用。
以上代码在内存中是这样的：
![](https://i-blog.csdnimg.cn/direct/38cd29e4bcc14ac0b47922f91afc5ac5.png)
所以s在使用完毕之后，并不会去销毁s1。
上面这种方式叫做不可变引用。
>可变引用

那如果想在函数中修改这个传入引用指向的堆内存。我们就需要用到可变引用。
```rust
fn main() {
    let mut s = String::from("hello");

    change(&mut s);
}

fn change(some_string: &mut String) {
    some_string.push_str(", world");
}
```
你以为这么简单就结束了吗？
可变引用有如下几个很重要的点，非常重要：
1.  同一作用域，特定数据只能有一个可变引用
那么什么叫作用域？
在早期的rust编译器引用和变量的作用域是相同的。也就是从定义的地方开始，到最近的一个`}`结束，那么这样其实很麻烦。因为如果你刚开始想使用一个不可变引用，然后又想定义一个可变引用，那怎么办（对的，可变引用和不可变引用也不能同时存在）？
后期Rust编译器进行了改进，一个引用的作用域是从定义到最后一次使用结束。
就像下面这个例子，r1的作用域是从定义开始到println结束。
那如果这个引用定义了之后都没有使用呢？那就等于立马结束了。也就是说，如果把下面代码的println删除掉。那么这个代码可以通过编译。
```rust
let mut s = String::from("hello");

let r1 = &mut s;
let r2 = &mut s;

println!("{}, {}", r1, r2);
```

这种限制的好处就是使 Rust 在编译期就避免数据竞争，数据竞争可由以下行为造成：

- 两个或更多的指针同时访问同一数据
- 至少有一个指针被用来写入数据
- 没有同步数据访问的机制
试想，如果都没有两个可变引用指向同一块内存的话，那么数据竞争从何谈起呢？
所以就从根源上断绝了数据竞争。
2. 同一个作用域，不能既有可变引用，又既有不可变引用
其实这个也很好理解，正在借用不可变引用的用户，肯定不希望他借用的东西，被另外一个人莫名其妙改变了。多个不可变借用被允许是因为没有人会去试图修改数据，每个人都只读这一份数据而不做修改，因此不用担心数据被污染。

# 所有权在方法上的应用
```rust
pub struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    fn width(&self) -> bool {
        self.width > 0
    }
}

fn main() {
    let rect1 = Rectangle {
        width: 30,
        height: 50,
    };

    if rect1.width() {
        println!("The rectangle has a nonzero width; it is {}", rect1.width);
    }
}
```
width是Rectangle结构体上的方法。
这里面的width的参数就是&self，表示不可变借用。在调用这个方法之后，所有权并不会转移。
需要注意的是，self 依然有所有权的概念：
- self 表示 Rectangle 的所有权转移到该方法中，这种形式用的较少
- &self 表示该方法对 Rectangle 的不可变借用
- &mut self 表示可变借用
一般来说都不会用到self，因为self会导致所有权转移，除非想在调用结束之后，原本的结构体直接就不用了。